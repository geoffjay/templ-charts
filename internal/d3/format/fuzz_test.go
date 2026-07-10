package d3format

import (
	"math"
	"testing"
)

// FuzzFormatString exercises the format-specifier parser with arbitrary spec
// strings and values. FormatString drives a hand-written grammar parser
// (parseSpec) plus padding/grouping/SI string manipulation, so a malformed
// spec, an extreme width/precision, or a non-finite value must never panic.
func FuzzFormatString(f *testing.F) {
	specs := []string{
		"", ",.2f", "+08.3e", "$,.0%", ".9s", "~g", "09d", "x", "(.2f",
		"d", ".0f", ",", "%", "^20.4g", "-.3~r",
	}
	values := []float64{
		0, -0, 1, -1, 1234.5, 1e308, -1e308, 1e-308,
		math.NaN(), math.Inf(1), math.Inf(-1),
	}
	for _, s := range specs {
		for _, v := range values {
			f.Add(s, v)
		}
	}

	f.Fuzz(func(t *testing.T, spec string, v float64) {
		// The contract under test is simply "does not panic". FormatString is
		// reachable from public chart Props via caller-supplied format strings.
		_ = FormatString(spec, v)
	})
}
