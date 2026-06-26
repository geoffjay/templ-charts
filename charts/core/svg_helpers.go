package core

import (
	"fmt"
	"math"
	"strings"
)

// fmtFloat formats a float for SVG output (matching d3-path rounding: 3 dp).
func fmtFloat(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	return fmt.Sprintf("%.3g", v)
}

// attrFloat emits `name="value"` for a non-zero float, or empty string.
func attrFloat(name string, v float64) string {
	if v == 0 {
		return ""
	}
	return fmt.Sprintf(` %s="%v"`, name, v)
}

// attrStr emits `name="value"` for a non-empty string, or empty string.
func attrStr(name, v string) string {
	if v == "" {
		return ""
	}
	return fmt.Sprintf(` %s="%s"`, name, escapeAttr(v))
}

// escapeAttr escapes a string for use in an XML attribute value.
func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// escapeText escapes a string for use in XML text content.
func escapeText(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
