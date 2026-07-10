package d3color

import "testing"

// FuzzRGBColor exercises the CSS-color string parser, which does substring
// indexing on user input (hex slicing, rgb()/hsl() function parsing, percent
// handling) — a classic source of index-out-of-range panics on short or
// malformed input. RGBColor is reachable from the public API via
// charts/colors with caller-supplied color strings. It may legitimately return
// nil for unparseable input; the contract under test is "does not panic".
func FuzzRGBColor(f *testing.F) {
	seeds := []string{
		"", "#", "#f", "#ff", "#fff", "#ffff", "#ffffff", "#ffffffff",
		"rgb(255,0,0)", "rgba(0,0,0,0.5)", "rgb(300%,-1,x)", "rgb(",
		"hsl(120,50%,50%)", "hsla(0,0%,0%,1)", "steelblue", "notacolor",
		"  #FFF  ", "RGB(1 2 3)",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		_ = RGBColor(s)
	})
}
