package d3format

import (
	"math"
	"testing"
)

func TestEmptySpec(t *testing.T) {
	f := Format("")
	if f == nil {
		t.Fatal("Format(\"\") returned nil")
	}
	// Default "g"-like, 6 sig figs.
	if got := f(1234.56789); got != "1234.57" {
		t.Errorf("empty spec for 1234.56789 = %q want 1234.57", got)
	}
}

func TestFixedPoint(t *testing.T) {
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{".2f", 3.14159, "3.14"},
		{".0f", 3.7, "4"},
		{".0f", 3.2, "3"},
		{".3f", 1.5, "1.500"},
		{"+.2f", -3.14, "-3.14"},
		{"+.2f", 3.14, "+3.14"},
		{"$.2f", 12.5, "$12.50"},
		{"(.2f", -3.14, "(3.14)"},
	}
	for _, c := range cases {
		got := FormatString(c.spec, c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestDecimalInteger(t *testing.T) {
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{"d", 42, "42"},
		{"d", 3.7, "4"},   // rounds
		{"d", -3.7, "-4"}, // negative rounds away from zero
		{"05d", 42, "00042"},
		{",d", 1234567, "1,234,567"},
		{"+d", 5, "+5"},
	}
	for _, c := range cases {
		got := FormatString(c.spec, c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestPercent(t *testing.T) {
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{".0%", 0.456, "46%"},
		{".2%", 0.456, "45.60%"},
		{"$.0%", 0.5, "$50%"},
	}
	for _, c := range cases {
		got := FormatString(c.spec, c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestSignVariants(t *testing.T) {
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{"+d", 5, "+5"},
		{" d", 5, " 5"},
		{"-d", 5, "5"},
		{"(.2f", -3.14, "(3.14)"},
	}
	for _, c := range cases {
		got := FormatString(c.spec, c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestWidthAlign(t *testing.T) {
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{"10.2f", 3.14, "      3.14"},
		{">10.2f", 3.14, "      3.14"},
		{"<10.2f", 3.14, "3.14      "},
		{"^10.2f", 3.14, "   3.14   "},
		{"=10.2f", -3.14, "-     3.14"},
		{".2f", 3.14, "3.14"},
	}
	for _, c := range cases {
		f := Format(c.spec)
		if f == nil {
			t.Errorf("%q: parser returned nil (want %q)", c.spec, c.want)
			continue
		}
		got := f(c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestGrouping(t *testing.T) {
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{",.0f", 1234567, "1,234,567"},
		{",.2f", 1234567.89, "1,234,567.89"},
		{"$,d", 1234567, "$1,234,567"},
	}
	for _, c := range cases {
		got := FormatString(c.spec, c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestTrim(t *testing.T) {
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{".5~f", 3.14, "3.14"},
		{".3~f", 3.5, "3.5"},
		{".0%", 0.5, "50%"},
		{".5~%", 0.123, "12.3%"},
	}
	for _, c := range cases {
		got := FormatString(c.spec, c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestSI(t *testing.T) {
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{".3s", 1234, "1.23k"},
		{".3s", 1234567, "1.23M"},
		{".3s", 0.00123, "1.23m"},
		{".2s", 1500, "1.5k"},
		{".2s", 1500000, "1.5M"},
	}
	for _, c := range cases {
		got := FormatString(c.spec, c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestSpecialValues(t *testing.T) {
	if got := FormatString(".2f", math.NaN()); got != "NaN" {
		t.Errorf("NaN = %q want NaN", got)
	}
	if got := FormatString(".2f", math.Inf(1)); got != "∞" {
		t.Errorf("+Inf = %q want ∞", got)
	}
	if got := FormatString(".2f", math.Inf(-1)); got != "-∞" {
		t.Errorf("-Inf = %q want -∞", got)
	}
}

func TestParseFailure(t *testing.T) {
	// invalid type char
	if f := Format("z"); f != nil {
		t.Errorf("Format(\"z\") should return nil")
	}
	// precision with no digits
	if f := Format("."); f != nil {
		t.Errorf("Format(\".\") should return nil")
	}
}

func TestGeneralType(t *testing.T) {
	// 'g' uses significant digits
	cases := []struct {
		spec string
		v    float64
		want string
	}{
		{".3g", 12345, "1.23e+04"},
		{".3g", 0.00123, "0.00123"},
		{".6g", 123.456789, "123.457"},
	}
	for _, c := range cases {
		got := FormatString(c.spec, c.v)
		if got != c.want {
			t.Errorf("%q %v = %q want %q", c.spec, c.v, got, c.want)
		}
	}
}

func TestHex(t *testing.T) {
	if got := FormatString("x", 255); got != "ff" {
		t.Errorf("x 255 = %q want ff", got)
	}
	if got := FormatString("X", 255); got != "FF" {
		t.Errorf("X 255 = %q want FF", got)
	}
	if got := FormatString("08x", 255); got != "000000ff" {
		t.Errorf("08x 255 = %q want 000000ff", got)
	}
}
