package core_test

import (
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
)

func TestSMILAnimate(t *testing.T) {
	got := core.SMILAnimate("r", "0", "6", "")
	for _, want := range []string{
		`attributeName="r"`, `from="0"`, `to="6"`, `begin="0s"`,
		`dur="0.6s"`, `fill="freeze"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("SMILAnimate missing %q in %q", want, got)
		}
	}
}

func TestSMILAnimate_NoOpWhenEqual(t *testing.T) {
	if got := core.SMILAnimate("opacity", "1", "1", ""); got != "" {
		t.Errorf("expected empty when from==to, got %q", got)
	}
}

func TestSMILFadeIn(t *testing.T) {
	got := core.SMILFadeIn("0.1s")
	if !strings.Contains(got, `attributeName="opacity"`) ||
		!strings.Contains(got, `from="0"`) || !strings.Contains(got, `to="1"`) ||
		!strings.Contains(got, `begin="0.1s"`) {
		t.Errorf("unexpected fade-in markup: %q", got)
	}
}

func TestStaggerBegin(t *testing.T) {
	cases := []struct {
		i       int
		stagger float64
		want    string
	}{
		{0, 0.05, "0s"}, // i<=0 → shared begin
		{3, 0, "0s"},    // stagger<=0 → shared begin
		{1, 0.05, "0.05s"},
		{4, 0.05, "0.2s"},
		{2, 0.1, "0.2s"},
	}
	for _, c := range cases {
		if got := core.StaggerBegin(c.i, c.stagger); got != c.want {
			t.Errorf("StaggerBegin(%d, %v) = %q, want %q", c.i, c.stagger, got, c.want)
		}
	}
}
