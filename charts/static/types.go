// Package static mirrors @nivo/static: a ChartType enum, a ChartsMapping
// registry mapping chart types to their components/defaults/runtime-props, a
// RenderChart dispatcher that applies static defaults + chart defaults +
// user props + whitelisted overrides and renders to an SVG string, and a
// Samples map of demo data for bar/line/pie.
package static

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// ChartType enumerates the supported chart types.
type ChartType string

const (
	ChartTypeBar  ChartType = "bar"
	ChartTypeLine ChartType = "line"
	ChartTypePie  ChartType = "pie"
)

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

// ChartsMapping maps each ChartType to its Mapping. Mirrors @nivo/static
// chartsMapping.
var ChartsMapping = map[ChartType]Mapping{
	ChartTypeBar: {
		Component:    barComponent{},
		RuntimeProps: []string{"width", "height", "colors", "groupMode"},
		Defaults: map[string]any{
			"margin": map[string]any{"top": 40, "right": 50, "bottom": 40, "left": 50},
		},
	},
	ChartTypeLine: {
		Component:    lineComponent{},
		RuntimeProps: []string{"width", "height", "colors"},
		Defaults: map[string]any{
			"margin": map[string]any{"top": 40, "right": 50, "bottom": 40, "left": 50},
		},
	},
	ChartTypePie: {
		Component:    pieComponent{},
		RuntimeProps: []string{"width", "height", "colors", "groupMode"},
		Defaults: map[string]any{
			"margin": map[string]any{"top": 40, "right": 50, "bottom": 40, "left": 50},
		},
	},
}

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
// `chartType` selects the chart; `props` is the chart-specific props struct
// (bar.BarProps, line.LineProps, or pie.PieProps); `override` is a map of
// runtime-overridable props (only keys in Mapping.RuntimeProps are applied).
func RenderChart(chartType ChartType, props any, override map[string]any) (string, error) {
	mapping, ok := ChartsMapping[chartType]
	if !ok {
		return "", fmt.Errorf("static: unknown chart type %q", chartType)
	}

	// Apply whitelisted overrides.
	over := map[string]any{}
	if override != nil {
		over = pick(override, mapping.RuntimeProps)
	}

	return mapping.Component.Render(props, over)
}

// mergeProps overlays `overlay` onto `base`. If base is a map[string]any, the
// overlay is merged key-by-key (recursive for nested maps). If base is a
// struct (bar.BarProps etc.), the overlay is applied via reflection where
// possible; otherwise base is returned as-is (the chart's applyDefaults
// handles the struct-level defaults).
//
// For the struct case, we convert the struct to a map, merge, and pass the map
// to the component's Render (which knows how to interpret it). In practice,
// the static package passes chart-specific props structs directly and relies
// on the component's applyDefaults for the rest; staticProps/defaults are
// applied as map overlays only when the base is already a map.
func mergeProps(base any, overlay map[string]any) any {
	if overlay == nil {
		return base
	}
	// If base is a map, merge recursively.
	if bm, ok := base.(map[string]any); ok {
		out := make(map[string]any, len(bm)+len(overlay))
		for k, v := range bm {
			out[k] = v
		}
		for k, v := range overlay {
			if existing, ok := out[k]; ok {
				if em, ok1 := existing.(map[string]any); ok1 {
					if vm, ok2 := v.(map[string]any); ok2 {
						out[k] = mergeProps(em, vm)
						continue
					}
				}
			}
			out[k] = v
		}
		return out
	}
	// If base is a struct, we can't easily merge a map onto it without
	// reflection. The chart components handle defaults internally via
	// applyDefaults, so we return the base as-is. The staticProps
	// (animate=false, isInteractive=false) are handled by passing them as
	// override fields that the component respects.
	return base
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

// --- chart component adapters ---

// barComponent adapts bar.Bar to the ChartComponent interface.
type barComponent struct{}

func (barComponent) Render(props any, override map[string]any) (string, error) {
	bp, ok := props.(bar.BarProps)
	if !ok {
		return "", fmt.Errorf("static: bar render expects bar.BarProps, got %T", props)
	}
	// Force static overrides.
	bp.Animate = false
	bp.Interactive = false
	bp.Theme = &theming.DefaultTheme
	// Apply chart defaults (margin).
	if bp.Margin == (core.Margin{}) {
		bp.Margin = core.Margin{Top: 40, Right: 50, Bottom: 40, Left: 50}
	}
	// Apply whitelisted overrides.
	if v, ok := override["width"]; ok {
		bp.Width = toFloat(v)
	}
	if v, ok := override["height"]; ok {
		bp.Height = toFloat(v)
	}
	if v, ok := override["groupMode"]; ok {
		if gm, ok := v.(bar.GroupMode); ok {
			bp.GroupMode = gm
		}
	}
	if v, ok := override["colors"]; ok {
		if cfg, err := colors.ParseOrdinalColorScaleConfig(v); err == nil {
			bp.Colors = cfg
		}
	}
	return renderComponent(bar.Bar(bp)), nil
}

// lineComponent adapts line.Line to the ChartComponent interface.
type lineComponent struct{}

func (lineComponent) Render(props any, override map[string]any) (string, error) {
	lp, ok := props.(line.LineProps)
	if !ok {
		return "", fmt.Errorf("static: line render expects line.LineProps, got %T", props)
	}
	lp.Animate = false
	lp.Interactive = false
	lp.Theme = &theming.DefaultTheme
	if lp.Margin == (core.Margin{}) {
		lp.Margin = core.Margin{Top: 40, Right: 50, Bottom: 40, Left: 50}
	}
	if v, ok := override["width"]; ok {
		lp.Width = toFloat(v)
	}
	if v, ok := override["height"]; ok {
		lp.Height = toFloat(v)
	}
	if v, ok := override["colors"]; ok {
		if cfg, err := colors.ParseOrdinalColorScaleConfig(v); err == nil {
			lp.Colors = cfg
		}
	}
	return renderComponent(line.Line(lp)), nil
}

// pieComponent adapts pie.Pie to the ChartComponent interface.
type pieComponent struct{}

func (pieComponent) Render(props any, override map[string]any) (string, error) {
	pp, ok := props.(pie.PieProps)
	if !ok {
		return "", fmt.Errorf("static: pie render expects pie.PieProps, got %T", props)
	}
	pp.Animate = false
	pp.Interactive = false
	pp.Theme = &theming.DefaultTheme
	if pp.Margin == (core.Margin{}) {
		pp.Margin = core.Margin{Top: 40, Right: 50, Bottom: 40, Left: 50}
	}
	if v, ok := override["width"]; ok {
		pp.Width = toFloat(v)
	}
	if v, ok := override["height"]; ok {
		pp.Height = toFloat(v)
	}
	if v, ok := override["colors"]; ok {
		if cfg, err := colors.ParseOrdinalColorScaleConfig(v); err == nil {
			pp.Colors = cfg
		}
	}
	return renderComponent(pie.Pie(pp)), nil
}

// toFloat coerces an any to float64.
func toFloat(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	}
	return 0
}

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}
