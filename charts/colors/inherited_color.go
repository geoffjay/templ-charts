package colors

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/geoffjay/templ-charts/charts/theming"
)

// ColorModifier is a single [kind, amount] color modifier pair. Kind is one
// of "brighter", "darker", "opacity". Mirrors @nivo/colors ColorModifier.
type ColorModifier [2]any

// InheritedColorConfig is the union of nivo's inherited-color strategies:
//   - static color string
//   - custom func(datum) string
//   - {theme: "path.into.theme"} — color from theme
//   - {from: "path.into.datum", modifiers: [...]} — color from datum with
//     optional brighter/darker/opacity modifiers
//
// Modeled as a struct with a discriminator field. Use the Is* helpers to
// detect the variant.
type InheritedColorConfig struct {
	// Static is the static color string variant. Set when Type == TypeStatic.
	Static string
	// Func is the custom-function variant. Set when Type == TypeFunc.
	Func func(datum any) string
	// ThemePath is the theme-path variant. Set when Type == TypeTheme.
	ThemePath string
	// FromPath is the datum-path for the from-context variant.
	FromPath string
	// Modifiers applies to the from-context variant.
	Modifiers []ColorModifier
	// Type discriminates the variant.
	Type InheritedColorType
}

// InheritedColorType discriminates InheritedColorConfig variants.
type InheritedColorType int

const (
	InheritedColorTypeStatic InheritedColorType = iota
	InheritedColorTypeFunc
	InheritedColorTypeTheme
	InheritedColorTypeFromContext
)

// NewStaticColor returns a static-color config.
func NewStaticColor(color string) InheritedColorConfig {
	return InheritedColorConfig{Type: InheritedColorTypeStatic, Static: color}
}

// NewFuncColor returns a custom-function config.
func NewFuncColor(fn func(any) string) InheritedColorConfig {
	return InheritedColorConfig{Type: InheritedColorTypeFunc, Func: fn}
}

// NewThemeColor returns a {theme: path} config.
func NewThemeColor(path string) InheritedColorConfig {
	return InheritedColorConfig{Type: InheritedColorTypeTheme, ThemePath: path}
}

// NewFromContextColor returns a {from: path, modifiers: ...} config.
func NewFromContextColor(path string, modifiers []ColorModifier) InheritedColorConfig {
	return InheritedColorConfig{Type: InheritedColorTypeFromContext, FromPath: path, Modifiers: modifiers}
}

// ParseInheritedColorConfig accepts the loose forms nivo allows (string,
// func, or map[string]any with "theme"/"from") and returns a typed
// InheritedColorConfig. Used by chart packages that accept `any` props.
func ParseInheritedColorConfig(v any) (InheritedColorConfig, error) {
	switch x := v.(type) {
	case nil:
		return InheritedColorConfig{Type: InheritedColorTypeStatic, Static: ""}, nil
	case string:
		return NewStaticColor(x), nil
	case func(any) string:
		return NewFuncColor(x), nil
	case func(any) any:
		return NewFuncColor(func(d any) string { return fmt.Sprintf("%v", x(d)) }), nil
	case map[string]any:
		if path, ok := x["theme"].(string); ok {
			return NewThemeColor(path), nil
		}
		if from, ok := x["from"].(string); ok {
			var mods []ColorModifier
			if rawMods, ok := x["modifiers"].([]any); ok {
				for _, m := range rawMods {
					if arr, ok := m.([]any); ok && len(arr) >= 2 {
						mods = append(mods, ColorModifier{arr[0], arr[1]})
					}
				}
			}
			return NewFromContextColor(from, mods), nil
		}
		return InheritedColorConfig{}, fmt.Errorf("invalid inherited color config: missing 'theme' or 'from'")
	}
	return InheritedColorConfig{}, fmt.Errorf("invalid inherited color config type: %T", v)
}

// GetInheritedColorGenerator returns a func(datum) string for the given
// config + theme. Mirrors @nivo/colors getInheritedColorGenerator.
func GetInheritedColorGenerator(config InheritedColorConfig, theme *theming.Theme) func(any) string {
	switch config.Type {
	case InheritedColorTypeFunc:
		return config.Func
	case InheritedColorTypeTheme:
		return func(any) string {
			if theme == nil {
				return ""
			}
			return getThemePath(*theme, config.ThemePath)
		}
	case InheritedColorTypeFromContext:
		getColor := func(d any) string {
			return getPathValue(d, config.FromPath)
		}
		if len(config.Modifiers) == 0 {
			return getColor
		}
		return func(d any) string {
			base := getColor(d)
			mods := make([][2]any, len(config.Modifiers))
			for i, m := range config.Modifiers {
				mods[i] = [2]any{m[0], m[1]}
			}
			return ApplyColorModifiers(base, mods)
		}
	default: // static
		return func(any) string { return config.Static }
	}
}

// getThemePath resolves a dot path into a Theme, returning the value as a
// string. Unknown paths return "".
func getThemePath(t theming.Theme, path string) string {
	v, ok := getPath(reflect.ValueOf(t), strings.Split(path, "."))
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// getPathValue resolves a dot path over a datum (struct fields, map keys).
func getPathValue(d any, path string) string {
	v, ok := getPath(reflect.ValueOf(d), strings.Split(path, "."))
	if !ok {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// getPath walks a dot path over struct fields and map keys.
func getPath(rv reflect.Value, parts []string) (any, bool) {
	for _, p := range parts {
		for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Struct:
			f := rv.FieldByName(p)
			if !f.IsValid() {
				rt := rv.Type()
				for i := 0; i < rt.NumField(); i++ {
					if strings.EqualFold(rt.Field(i).Name, p) {
						f = rv.Field(i)
						break
					}
				}
				if !f.IsValid() {
					return nil, false
				}
			}
			rv = f
		case reflect.Map:
			mk := reflect.ValueOf(p)
			if !mk.Type().AssignableTo(rv.Type().Key()) {
				return nil, false
			}
			mv := rv.MapIndex(mk)
			if !mv.IsValid() {
				return nil, false
			}
			rv = mv
		default:
			return nil, false
		}
	}
	return rv.Interface(), true
}
