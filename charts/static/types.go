// Package static mirrors @nivo/static: a ChartType enum, a ChartsMapping
// registry mapping chart types to their components/defaults/runtime-props, a
// RenderChart dispatcher that applies static defaults + chart defaults +
// user props + whitelisted overrides and renders to an SVG string, and a
// Samples map of ready-to-render demo data for every chart family.
//
// The per-family wiring (ChartType constants, ChartsMapping entries, and the
// Samples registry) lives in registry.go; this file holds the generic
// reflection-based adapter that every family shares. Rather than a hand-written
// component per chart, one reflectComponent forces the static flags
// (animate/isInteractive off, default theme), applies a default margin, and
// applies the whitelisted overrides by field name — set-if-assignable, so a
// family whose Colors field is a bespoke type (heatmap, calendar) simply skips
// the generic `colors` override. The only per-family code is a one-line render
// closure that type-asserts the props and calls the templ constructor.
package static

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// ChartType names a supported chart family. The concrete constants live in
// registry.go.
type ChartType string

// ChartComponent is the render interface every chart type implements. Render
// applies static defaults + chart defaults + user props + whitelisted overrides
// and returns the SVG string.
type ChartComponent interface {
	Render(props any, override map[string]any) (string, error)
}

// Mapping is one entry in ChartsMapping: the component, the set of prop names
// that may be overridden at render time (runtimeProps), and the chart-specific
// static defaults (e.g. margin).
type Mapping struct {
	Component    ChartComponent
	RuntimeProps []string
	Defaults     map[string]any
}

// defaultMargin is the static-render margin applied to any family whose Margin
// prop is left zero (matches the legacy bar/line/pie behaviour).
var defaultMargin = core.Margin{Top: 40, Right: 50, Bottom: 40, Left: 50}

// staticProps mirrors nivo's staticProps: the base overrides applied to every
// static render (no animation, no interactivity, no wrapper, empty theme).
var staticProps = map[string]any{
	"animate":       false,
	"isInteractive": false,
	"renderWrapper": false,
	"theme":         map[string]any{},
}

// RenderChart mirrors @nivo/static renderChart: applies staticProps +
// chart.defaults + props + whitelisted override, renders to SVG string.
// `chartType` selects the chart; `props` is the chart-specific props struct;
// `override` is a map of runtime-overridable props (only keys in
// Mapping.RuntimeProps are applied).
func RenderChart(chartType ChartType, props any, override map[string]any) (string, error) {
	mapping, ok := ChartsMapping[chartType]
	if !ok {
		return "", fmt.Errorf("static: unknown chart type %q", chartType)
	}
	over := map[string]any{}
	if override != nil {
		over = pick(override, mapping.RuntimeProps)
	}
	return mapping.Component.Render(props, over)
}

// reflectComponent is the shared, reflection-based ChartComponent. `render`
// type-asserts the (defaults-applied) props and constructs the templ component.
type reflectComponent struct {
	render func(props any) (templ.Component, error)
}

// Render applies the static overrides + default margin + whitelisted overrides
// to a copy of props (via reflection), then renders.
func (rc reflectComponent) Render(props any, override map[string]any) (string, error) {
	nv, ok := mutableCopy(props)
	if ok {
		// Force static flags.
		setBoolField(nv, "Animate", false)
		setBoolField(nv, "Interactive", false)
		setThemeField(nv, "Theme", &theming.DefaultTheme)
		// Default margin when unset.
		if mf := nv.FieldByName("Margin"); mf.IsValid() && mf.CanSet() {
			if m, ok := mf.Interface().(core.Margin); ok && m == (core.Margin{}) {
				mf.Set(reflect.ValueOf(defaultMargin))
			}
		}
		// Whitelisted overrides.
		for k, v := range override {
			applyOverride(nv, k, v)
		}
		props = nv.Interface()
	}
	comp, err := rc.render(props)
	if err != nil {
		return "", err
	}
	return renderComponent(comp), nil
}

// mutableCopy returns an addressable copy of a struct value so its fields can
// be set via reflection. Returns ok=false for non-structs.
func mutableCopy(props any) (reflect.Value, bool) {
	v := reflect.ValueOf(props)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	nv := reflect.New(v.Type()).Elem()
	nv.Set(v)
	return nv, true
}

// applyOverride maps a runtime-prop key to a struct field and sets it when the
// value is assignable. The `colors` key is parsed into an ordinal color config
// first (skipped for families whose Colors field is a bespoke type).
func applyOverride(v reflect.Value, key string, val any) {
	switch key {
	case "width":
		setFloatField(v, "Width", val)
	case "height":
		setFloatField(v, "Height", val)
	case "colors":
		if cfg, err := colors.ParseOrdinalColorScaleConfig(val); err == nil {
			setAssignableField(v, "Colors", reflect.ValueOf(cfg))
		}
	default:
		setAssignableField(v, capitalize(key), reflect.ValueOf(val))
	}
}

// setAssignableField sets field `name` to rv when the field exists, is
// settable, and rv is assignable to it (otherwise it is silently skipped).
func setAssignableField(v reflect.Value, name string, rv reflect.Value) {
	f := v.FieldByName(name)
	if !f.IsValid() || !f.CanSet() || !rv.IsValid() {
		return
	}
	if rv.Type().AssignableTo(f.Type()) {
		f.Set(rv)
	}
}

func setBoolField(v reflect.Value, name string, b bool) {
	f := v.FieldByName(name)
	if f.IsValid() && f.CanSet() && f.Kind() == reflect.Bool {
		f.SetBool(b)
	}
}

func setThemeField(v reflect.Value, name string, theme *theming.Theme) {
	f := v.FieldByName(name)
	if f.IsValid() && f.CanSet() {
		rv := reflect.ValueOf(theme)
		if rv.Type().AssignableTo(f.Type()) {
			f.Set(rv)
		}
	}
}

func setFloatField(v reflect.Value, name string, val any) {
	f := v.FieldByName(name)
	if !f.IsValid() || !f.CanSet() {
		return
	}
	if fv, ok := asFloat(val); ok && f.Kind() == reflect.Float64 {
		f.SetFloat(fv)
	}
}

func asFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	}
	return 0, false
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// pick returns a new map containing only the keys from `src` that are listed
// in `allowed`. Mirrors lodash.pick.
func pick(src map[string]any, allowed []string) map[string]any {
	out := map[string]any{}
	for _, k := range allowed {
		if v, ok := src[k]; ok {
			out[k] = v
		}
	}
	return out
}

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}
