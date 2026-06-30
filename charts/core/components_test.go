package core_test

import (
	"context"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
)

func renderWrapper(t *testing.T, props core.SvgWrapperProps) string {
	t.Helper()
	var b strings.Builder
	if err := core.SvgWrapper(props, "<g></g>").Render(context.Background(), &b); err != nil {
		t.Fatalf("SvgWrapper.Render: %v", err)
	}
	return b.String()
}

func TestSvgWrapper_TitleDesc(t *testing.T) {
	out := renderWrapper(t, core.SvgWrapperProps{
		Width: 100, Height: 80, Role: "img",
		Title: "Sales by region",
		Desc:  "A bar chart of quarterly sales across regions.",
	})
	if !strings.Contains(out, "<title>Sales by region</title>") {
		t.Errorf("expected <title>; got %q", out)
	}
	if !strings.Contains(out, "<desc>A bar chart of quarterly sales across regions.</desc>") {
		t.Errorf("expected <desc>; got %q", out)
	}
	// <title> must precede the inner content (first child for AT name resolution).
	if ti, gi := strings.Index(out, "<title>"), strings.Index(out, "<g>"); ti == -1 || ti > gi {
		t.Errorf("<title> should come before the inner <g>")
	}
}

func TestSvgWrapper_NoTitleDescByDefault(t *testing.T) {
	out := renderWrapper(t, core.SvgWrapperProps{Width: 100, Height: 80, Role: "img"})
	if strings.Contains(out, "<title>") || strings.Contains(out, "<desc>") {
		t.Errorf("title/desc must not be emitted when unset; got %q", out)
	}
	if !strings.Contains(out, `role="img"`) {
		t.Errorf("expected role=img")
	}
}

func TestSvgWrapper_Focusable(t *testing.T) {
	out := renderWrapper(t, core.SvgWrapperProps{Width: 10, Height: 10, IsFocusable: true})
	if !strings.Contains(out, `tabindex="0"`) {
		t.Errorf("IsFocusable should emit tabindex=0")
	}
}
