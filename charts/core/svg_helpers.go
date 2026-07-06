package core

import (
	"fmt"
	"math"
	"strings"
)

// fmtFloat formats a float for SVG output (up to 3 decimal places, trailing
// zeros trimmed).
func fmtFloat(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	s := fmt.Sprintf("%.3f", v)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "-0" {
		// Normalize -0 (tiny negatives rounded to zero) so output is
		// stable across platforms (arm64 vs amd64).
		s = "0"
	}
	return s
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
