package bar

import (
	"github.com/geoffjay/templ-charts/charts/theming"
)

// resolveTheme returns props.Theme or the default theme when nil.
func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

// themeBackground returns the theme background color (or "" for transparent).
func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}

// themeLabelFontSize extracts the labels.text.fontSize as a float64. nivo
// stores font sizes as any (number or string); we coerce to float.
func themeLabelFontSize(t *theming.Theme) float64 {
	if t == nil {
		return 0
	}
	return toFloatAny(t.Labels.Text.FontSize)
}

// toFloatAny coerces an any-typed number (float64/int/string) to float64.
func toFloatAny(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case string:
		// ignore "11px"-style strings; return 0
		return 0
	}
	return 0
}
