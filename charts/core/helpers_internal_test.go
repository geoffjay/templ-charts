package core

import (
	"math"
	"testing"
	"time"
)

func TestFmtFloat(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0"},
		{2, "2"},
		{1.5, "1.5"},
		{1.23456, "1.235"},
		{-3.10, "-3.1"},
		{math.NaN(), "0"},
		{math.Inf(1), "0"},
		{math.Inf(-1), "0"},
		// Tiny negatives rounding to zero must not emit "-0".
		{-0.0001, "0"},
	}
	for _, c := range cases {
		if got := fmtFloat(c.in); got != c.want {
			t.Errorf("fmtFloat(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestAttrHelpers(t *testing.T) {
	if got := attrFloat("x", 0); got != "" {
		t.Errorf("attrFloat zero = %q, want empty", got)
	}
	if got := attrFloat("x", 2.5); got != ` x="2.5"` {
		t.Errorf("attrFloat = %q", got)
	}
	if got := attrStr("fill", ""); got != "" {
		t.Errorf("attrStr empty = %q, want empty", got)
	}
	if got := attrStr("fill", `a"b<c>&d`); got != ` fill="a&quot;b&lt;c&gt;&amp;d"` {
		t.Errorf("attrStr = %q", got)
	}
}

func TestEscapeText(t *testing.T) {
	if got := escapeText(`<a> & "b"`); got != `&lt;a&gt; &amp; "b"` {
		t.Errorf("escapeText = %q", got)
	}
}

func TestSmallHelpers(t *testing.T) {
	if opacityOr1(0) != 1 || opacityOr1(-0.5) != 1 || opacityOr1(0.3) != 0.3 {
		t.Errorf("opacityOr1 wrong")
	}
	if patternNum(0, 4) != 4 || patternNum(2, 4) != 2 {
		t.Errorf("patternNum wrong")
	}
	if patternStr("", "#fff") != "#fff" || patternStr("#000", "#fff") != "#000" {
		t.Errorf("patternStr wrong")
	}
	if textAnchorOr("", "middle") != "middle" || textAnchorOr("start", "middle") != "start" {
		t.Errorf("textAnchorOr wrong")
	}
	if got := degreesToRadians(180); math.Abs(got-math.Pi) > 1e-12 {
		t.Errorf("degreesToRadians(180) = %v", got)
	}
}

func TestItoa(t *testing.T) {
	cases := []struct {
		in   int
		want string
	}{{0, "0"}, {7, "7"}, {12345, "12345"}, {-42, "-42"}}
	for _, c := range cases {
		if got := itoa(c.in); got != c.want {
			t.Errorf("itoa(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestToFloat64(t *testing.T) {
	type myInt int8
	type myUint uint16
	type myFloat float32
	tm := time.UnixMilli(1700000000000).UTC()
	cases := []struct {
		name string
		in   any
		want float64
	}{
		{"nil", nil, 0},
		{"float64", 1.5, 1.5},
		{"float32", float32(2.5), 2.5},
		{"int", 3, 3},
		{"int64", int64(4), 4},
		{"int32", int32(5), 5},
		{"bool true", true, 1},
		{"bool false", false, 0},
		{"numeric string", "6.5", 6.5},
		{"bad string", "x", 0},
		{"time", tm, 1700000000000},
		{"named int", myInt(7), 7},
		{"named uint", myUint(8), 8},
		{"named float", myFloat(9.5), 9.5},
		{"unsupported", struct{}{}, 0},
	}
	for _, c := range cases {
		if got := toFloat64(c.in); got != c.want {
			t.Errorf("%s: toFloat64(%v) = %v, want %v", c.name, c.in, got, c.want)
		}
	}
}

func TestToString(t *testing.T) {
	tm := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	if toString(nil) != "" || toString("s") != "s" || toString(42) != "42" {
		t.Errorf("toString basic cases wrong")
	}
	if got := toString(tm); got != "2024-05-01T00:00:00Z" {
		t.Errorf("toString(time) = %q", got)
	}
}

func TestToBool(t *testing.T) {
	cases := []struct {
		in   any
		want bool
	}{
		{nil, false},
		{true, true},
		{false, false},
		{0.0, false},
		{1.5, true},
		{"true", true},
		{"false", false},
		{"", false},
		// Unparseable non-empty strings are truthy.
		{"yes", true},
		{42, true},
	}
	for _, c := range cases {
		if got := toBool(c.in); got != c.want {
			t.Errorf("toBool(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestToTime(t *testing.T) {
	tm := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	if got, ok := toTime(tm); !ok || !got.Equal(tm) {
		t.Errorf("toTime(time.Time) = %v, %v", got, ok)
	}
	if got, ok := toTime(&tm); !ok || !got.Equal(tm) {
		t.Errorf("toTime(*time.Time) = %v, %v", got, ok)
	}
	if _, ok := toTime((*time.Time)(nil)); ok {
		t.Errorf("toTime(nil ptr) should fail")
	}
	if got, ok := toTime("2024-05-01T00:00:00Z"); !ok || !got.Equal(tm) {
		t.Errorf("toTime(rfc3339) = %v, %v", got, ok)
	}
	if _, ok := toTime("nope"); ok {
		t.Errorf("toTime(bad string) should fail")
	}
	if _, ok := toTime(nil); ok {
		t.Errorf("toTime(nil) should fail")
	}
	if _, ok := toTime(42); ok {
		t.Errorf("toTime(int) should fail")
	}
}

func TestComputeMarkerLabel(t *testing.T) {
	const (
		w  = 100.0
		h  = 80.0
		ox = 10.0
		oy = 5.0
	)
	cases := []struct {
		axis, position, orientation string
		x, y, rotation              float64
		anchor                      string
	}{
		{"x", "top-left", "horizontal", -ox, oy, 0, "end"},
		{"x", "top", "horizontal", 0, -oy, 0, "middle"},
		{"x", "top", "vertical", 0, -oy, -90, "start"},
		{"x", "top-right", "horizontal", ox, oy, 0, "start"},
		{"x", "top-right", "vertical", ox, oy, -90, "end"},
		{"x", "right", "horizontal", ox, h / 2, 0, "start"},
		{"x", "right", "vertical", ox, h / 2, -90, "middle"},
		{"x", "bottom-right", "horizontal", ox, h - oy, 0, "start"},
		{"x", "bottom", "horizontal", 0, h + oy, 0, "middle"},
		{"x", "bottom", "vertical", 0, h + oy, -90, "end"},
		{"x", "bottom-left", "horizontal", -ox, h - oy, 0, "end"},
		{"x", "bottom-left", "vertical", -ox, h - oy, -90, "start"},
		{"x", "left", "horizontal", -ox, h / 2, 0, "end"},
		{"x", "left", "vertical", -ox, h / 2, -90, "middle"},
		{"y", "top-left", "horizontal", ox, -oy, 0, "start"},
		{"y", "top", "horizontal", w / 2, -oy, 0, "middle"},
		{"y", "top", "vertical", w / 2, -oy, -90, "start"},
		{"y", "top-right", "horizontal", w - ox, -oy, 0, "end"},
		{"y", "top-right", "vertical", w - ox, -oy, -90, "start"},
		{"y", "right", "horizontal", w + ox, 0, 0, "start"},
		{"y", "right", "vertical", w + ox, 0, -90, "middle"},
		{"y", "bottom-right", "horizontal", w - ox, oy, 0, "end"},
		{"y", "bottom", "horizontal", w / 2, oy, 0, "middle"},
		{"y", "bottom", "vertical", w / 2, oy, -90, "end"},
		{"y", "bottom-left", "horizontal", ox, oy, 0, "start"},
		{"y", "bottom-left", "vertical", ox, oy, -90, "end"},
		{"y", "left", "horizontal", -ox, 0, 0, "end"},
		{"y", "left", "vertical", -ox, 0, -90, "middle"},
	}
	for _, c := range cases {
		x, y, rot, anchor := computeMarkerLabel(c.axis, w, h, c.position, ox, oy, c.orientation)
		if x != c.x || y != c.y || rot != c.rotation || anchor != c.anchor {
			t.Errorf("computeMarkerLabel(%s %s %s) = (%v, %v, %v, %q), want (%v, %v, %v, %q)",
				c.axis, c.position, c.orientation, x, y, rot, anchor, c.x, c.y, c.rotation, c.anchor)
		}
	}
}

func TestComputeMarkerLine(t *testing.T) {
	m := CartesianMarker{Axis: "y", Value: 1}
	got := computeMarkerLine(m, 100, 80, nil, func(any) float64 { return 25 })
	if got != (markerLine{X: 0, Y: 25, X2: 100, Y2: 0}) {
		t.Errorf("y marker line = %+v", got)
	}
	m = CartesianMarker{Axis: "x", Value: 1}
	got = computeMarkerLine(m, 100, 80, func(any) float64 { return 60 }, nil)
	if got != (markerLine{X: 60, Y: 0, X2: 0, Y2: 80}) {
		t.Errorf("x marker line = %+v", got)
	}
}

func TestComputeMarkerLayoutDefaults(t *testing.T) {
	// Empty position/orientation/offsets default to top-right, horizontal, 14.
	l := computeMarkerLayout(CartesianMarker{Axis: "y"}, 100, 80)
	if l.X != 86 || l.Y != -14 || l.Rotation != 0 || l.TextAnchor != "end" {
		t.Errorf("default layout = %+v", l)
	}
}
