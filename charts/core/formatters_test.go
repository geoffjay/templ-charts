package core_test

import (
	"testing"
	"time"

	"github.com/geoffjay/templ-charts/charts/core"
)

func TestGetValueFormatter_Default(t *testing.T) {
	f := core.GetValueFormatter[float64](nil)
	if got := f(1.5); got != "1.5" {
		t.Errorf("nil format = %q, want %q", got, "1.5")
	}
	e := core.GetValueFormatter[float64]("")
	if got := e(2.0); got != "2" {
		t.Errorf("empty format = %q, want %q", got, "2")
	}
	w := core.GetValueFormatter[float64]("   ")
	if got := w(3.0); got != "3" {
		t.Errorf("blank format = %q, want %q", got, "3")
	}
}

func TestGetValueFormatter_Func(t *testing.T) {
	f := core.GetValueFormatter[float64](func(v float64) string { return "v!" })
	if got := f(1); got != "v!" {
		t.Errorf("func format = %q, want %q", got, "v!")
	}
}

func TestGetValueFormatter_D3Spec(t *testing.T) {
	f := core.GetValueFormatter[float64](".2f")
	if got := f(3.14159); got != "3.14" {
		t.Errorf(".2f = %q, want %q", got, "3.14")
	}
	// String values are coerced to float64 before formatting.
	s := core.GetValueFormatter[any](".1f")
	if got := s("2.55"); got != "2.5" && got != "2.6" {
		t.Errorf(".1f on string = %q, want 2.5 or 2.6", got)
	}
}

func TestGetValueFormatter_TimeSpec(t *testing.T) {
	tm := time.Date(2024, 5, 1, 12, 30, 0, 0, time.UTC)
	f := core.GetValueFormatter[time.Time]("time:%Y-%m-%d")
	if got := f(tm); got != "2024-05-01" {
		t.Errorf("time spec = %q, want %q", got, "2024-05-01")
	}
	g := core.GetValueFormatter[any]("time:%Y")
	if got := g(&tm); got != "2024" {
		t.Errorf("time spec on *time.Time = %q, want %q", got, "2024")
	}
	if got := g("2024-05-01T12:30:00Z"); got != "2024" {
		t.Errorf("time spec on RFC3339 string = %q, want %q", got, "2024")
	}
	// Non-time values fall back to %v formatting.
	if got := g(42); got != "42" {
		t.Errorf("time spec on int = %q, want %q", got, "42")
	}
	if got := g("not a date"); got != "not a date" {
		t.Errorf("time spec on bad string = %q, want passthrough", got)
	}
	if got := g((*time.Time)(nil)); got != "<nil>" {
		t.Errorf("time spec on nil *time.Time = %q, want %q", got, "<nil>")
	}
}

func TestGetValueFormatter_UnsupportedType(t *testing.T) {
	f := core.GetValueFormatter[int](42)
	if got := f(7); got != "7" {
		t.Errorf("unsupported format = %q, want %q", got, "7")
	}
}
