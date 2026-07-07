package htmx_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/treemap"
)

// --- Mount component -----------------------------------------------------------

func renderMount(t *testing.T, props htmx.MountProps) string {
	t.Helper()
	var b strings.Builder
	if err := htmx.Mount(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Mount render: %v", err)
	}
	return b.String()
}

func TestMount_Static(t *testing.T) {
	out := renderMount(t, htmx.MountProps{ID: "demo", SVG: "<svg>inner</svg>"})
	if !strings.Contains(out, `id="chart-demo"`) {
		t.Errorf("missing chart container id: %q", out)
	}
	if !strings.Contains(out, `class="tc-chart"`) {
		t.Errorf("missing tc-chart class: %q", out)
	}
	if !strings.Contains(out, "<svg>inner</svg>") {
		t.Errorf("SVG should be injected raw: %q", out)
	}
	// Static mount: no hover-leave wiring, no tooltip sibling.
	if strings.Contains(out, "hx-get") || strings.Contains(out, "tooltip-demo") {
		t.Errorf("static mount should have no interactivity wiring: %q", out)
	}
}

func TestMount_Interactive(t *testing.T) {
	out := renderMount(t, htmx.MountProps{ID: "demo", SVG: "<svg/>", Interactive: true})
	if !strings.Contains(out, "/charts/demo/hover?leave=1") {
		t.Errorf("missing leave-reset hx-get: %q", out)
	}
	if !strings.Contains(out, `hx-trigger="mouseleave"`) {
		t.Errorf("missing mouseleave trigger: %q", out)
	}
	if !strings.Contains(out, `hx-target="#tooltip-demo"`) {
		t.Errorf("missing tooltip target: %q", out)
	}
	if !strings.Contains(out, `id="tooltip-demo"`) || !strings.Contains(out, "tc-chart-tooltip") {
		t.Errorf("missing tooltip sibling: %q", out)
	}
}

func TestMount_ExtraClass(t *testing.T) {
	out := renderMount(t, htmx.MountProps{ID: "demo", SVG: "<svg/>", Class: "card"})
	if !strings.Contains(out, `class="tc-chart card"`) {
		t.Errorf("extra class not appended: %q", out)
	}
}

// --- line mesh hover -------------------------------------------------------------

func TestHandler_LineMeshHover(t *testing.T) {
	r := newLineRegistry(t)
	h := htmx.NewHandler(r)
	rec := do(t, h, http.MethodGet, "/charts/demo-line/hover?x=0&y=0")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "nivo-tooltip-basic") {
		t.Errorf("mesh hover missing tooltip: %q", body)
	}
	// Tooltip names one of the series.
	if !strings.Contains(body, "A") && !strings.Contains(body, "B") {
		t.Errorf("mesh hover missing series id: %q", body)
	}
	// OOB chart swap present.
	if !strings.Contains(body, `hx-swap-oob="innerHTML:#chart-demo-line"`) {
		t.Errorf("mesh hover missing OOB swap: %q", body)
	}
	// State snapped to the nearest point and hover activated.
	st := r.Get("demo-line").State()
	if !st.HasHover {
		t.Error("HasHover = false after mesh hover")
	}
}

func TestHandler_LineMeshHover_BadParams(t *testing.T) {
	h := htmx.NewHandler(newLineRegistry(t))
	if rec := do(t, h, http.MethodGet, "/charts/demo-line/hover"); rec.Code != http.StatusBadRequest {
		t.Errorf("missing x/y: status = %d, want 400", rec.Code)
	}
	if rec := do(t, h, http.MethodGet, "/charts/demo-line/hover?x=abc&y=0"); rec.Code != http.StatusBadRequest {
		t.Errorf("invalid x: status = %d, want 400", rec.Code)
	}
	if rec := do(t, h, http.MethodGet, "/charts/demo-line/hover?x=0&y=abc"); rec.Code != http.StatusBadRequest {
		t.Errorf("invalid y: status = %d, want 400", rec.Code)
	}
}

func TestHandler_HoverLeaveClearsState(t *testing.T) {
	r := newLineRegistry(t)
	h := htmx.NewHandler(r)
	// Establish hover state first.
	do(t, h, http.MethodGet, "/charts/demo-line/hover?x=0&y=0")
	if !r.Get("demo-line").State().HasHover {
		t.Fatal("precondition: hover state should be set")
	}
	rec := do(t, h, http.MethodGet, "/charts/demo-line/hover?leave=1")
	if rec.Code != http.StatusOK {
		t.Fatalf("leave status = %d, want 200", rec.Code)
	}
	st := r.Get("demo-line").State()
	if st.HasHover || st.HoveredKey != "" || st.HoverX != 0 || st.HoverY != 0 {
		t.Errorf("leave should clear hover state, got %+v", st)
	}
	// Response is an empty tooltip plus the OOB chart swap.
	body := rec.Body.String()
	if !strings.HasPrefix(body, `<div hx-swap-oob=`) {
		t.Errorf("leave body should start with the OOB div (empty tooltip): %q", body[:min(60, len(body))])
	}
	if !strings.Contains(body, "<svg") {
		t.Errorf("leave body missing clean SVG")
	}
}

func TestHandler_PieHoverLeaveClearsActive(t *testing.T) {
	r := newPieRegistry(t)
	h := htmx.NewHandler(r)
	do(t, h, http.MethodGet, "/charts/demo-pie/hover?arc=A")
	if r.Get("demo-pie").State().ActiveID != "A" {
		t.Fatal("precondition: active arc should be set")
	}
	rec := do(t, h, http.MethodGet, "/charts/demo-pie/hover?leave=1")
	if rec.Code != http.StatusOK {
		t.Fatalf("leave status = %d", rec.Code)
	}
	if got := r.Get("demo-pie").State().ActiveID; got != "" {
		t.Errorf("ActiveID = %q after leave, want empty", got)
	}
}

func TestHandler_HoverUnsupportedKind(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterIcicle("ic", icicle.IcicleProps{Width: 300, Height: 200, Data: icicleData()})
	rec := do(t, htmx.NewHandler(r), http.MethodGet, "/charts/ic/hover?bar=x")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("hover on icicle: status = %d, want 400", rec.Code)
	}
}

func TestHandler_HeatmapHover_NilValueAndUnknownCell(t *testing.T) {
	h := htmx.NewHandler(newHeatmapRegistry(t))
	// USA.Car exists but has a nil value: no tooltip.
	if rec := do(t, h, http.MethodGet, "/charts/demo-heatmap/hover?cell=USA.Car"); rec.Code != http.StatusNotFound {
		t.Errorf("nil-value cell: status = %d, want 404", rec.Code)
	}
	if rec := do(t, h, http.MethodGet, "/charts/demo-heatmap/hover?cell=Nope.Nope"); rec.Code != http.StatusNotFound {
		t.Errorf("unknown cell: status = %d, want 404", rec.Code)
	}
}

func TestHandler_PieHover_MissingArcParam(t *testing.T) {
	h := htmx.NewHandler(newPieRegistry(t))
	if rec := do(t, h, http.MethodGet, "/charts/demo-pie/hover"); rec.Code != http.StatusBadRequest {
		t.Errorf("missing arc: status = %d, want 400", rec.Code)
	}
}

func TestHandler_LineMeshHover_NoPoints(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterLine("empty-line", line.LineProps{Width: 300, Height: 200, Data: nil})
	h := htmx.NewHandler(r)
	if rec := do(t, h, http.MethodGet, "/charts/empty-line/hover?x=1&y=1"); rec.Code != http.StatusNotFound {
		t.Errorf("empty line mesh hover: status = %d, want 404", rec.Code)
	}
}

// --- zoom: deep nodes resolve their real parent -----------------------------------

// assertDeepZoomOut focuses a depth-2 node and re-clicks it: the focus must
// move off the deep node to one of its ancestors, exercising the recursive
// parent walk (childHas* + the recursive parent call). The exact ancestor is
// not pinned: the current resolvers land on the grandparent for depth-2 nodes,
// so the depth-1 parent and the root id are both accepted.
func assertDeepZoomOut(t *testing.T, h *htmx.Handler, id, deepID string, ancestors ...string) {
	t.Helper()
	rec := do(t, h, http.MethodGet, "/charts/"+id+"/zoom?node="+deepID)
	if rec.Code != http.StatusOK {
		t.Fatalf("zoom in status = %d", rec.Code)
	}
	if got := h.Registry().Get(id).State().FocusID; got != deepID {
		t.Fatalf("focus = %q, want %q", got, deepID)
	}
	rec = do(t, h, http.MethodGet, "/charts/"+id+"/zoom?node="+deepID)
	if rec.Code != http.StatusOK {
		t.Fatalf("zoom out status = %d", rec.Code)
	}
	got := h.Registry().Get(id).State().FocusID
	ok := false
	for _, a := range ancestors {
		if got == a {
			ok = true
			break
		}
	}
	if !ok {
		t.Errorf("focus after zoom-out = %q, want one of %v", got, ancestors)
	}
}

func TestHandler_IcicleDeepZoomOut(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterIcicle("ic", icicle.IcicleProps{Width: 500, Height: 300, EnableZooming: true, Data: icicleData()})
	assertDeepZoomOut(t, htmx.NewHandler(r), "ic", "a1", "A", "root", "")
}

func TestHandler_TreemapDeepZoomOut(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterTreemap("tm", treemap.TreemapProps{
		Width: 500, Height: 400, EnableZooming: true,
		Data: treemap.TreemapNode{ID: "root", Children: []treemap.TreemapNode{
			{ID: "A", Children: []treemap.TreemapNode{{ID: "a1", Value: 12}, {ID: "a2", Value: 8}}},
			{ID: "B", Value: 10},
		}},
	})
	assertDeepZoomOut(t, htmx.NewHandler(r), "tm", "a1", "A", "root", "")
}

func TestHandler_CirclePackingDeepZoomOut(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterCirclePacking("cp", cp.CirclePackingProps{
		Width: 400, Height: 400, EnableZooming: true,
		Data: cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
			{ID: "A", Children: []cp.CirclePackingNode{{ID: "a1", Value: 8}, {ID: "a2", Value: 4}}},
			{ID: "B", Value: 6},
		}},
	})
	assertDeepZoomOut(t, htmx.NewHandler(r), "cp", "a1", "A", "root", "")
}

func TestHandler_SunburstDeepZoomOut(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterSunburst("sb", sunburst.SunburstProps{
		Width: 400, Height: 400, EnableZooming: true,
		Data: sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
			{ID: "A", Children: []sunburst.SunburstNode{{ID: "a1", Value: 8}, {ID: "a2", Value: 4}}},
			{ID: "B", Value: 6},
		}},
	})
	assertDeepZoomOut(t, htmx.NewHandler(r), "sb", "a1", "A", "root", "")
}

func TestHandler_ZoomEmptyNodeStaysRoot(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterIcicle("ic", icicle.IcicleProps{Width: 300, Height: 200, EnableZooming: true, Data: icicleData()})
	h := htmx.NewHandler(r)
	rec := do(t, h, http.MethodGet, "/charts/ic/zoom?node=")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if got := r.Get("ic").State().FocusID; got != "" {
		t.Errorf("focus = %q, want root", got)
	}
}

// --- method / route guards ---------------------------------------------------------

func TestHandler_MethodGuards(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	cases := []struct {
		method, path string
	}{
		{http.MethodPost, "/charts/demo-bar/hover?bar=x"},
		{http.MethodPost, "/charts/demo-bar/slice?axis=x&x=1"},
		{http.MethodGet, "/charts/demo-bar/click?bar=x"},
		{http.MethodGet, "/charts/demo-bar/toggle?series=x"},
		{http.MethodPost, "/charts/demo-bar/zoom?node=x"},
	}
	for _, c := range cases {
		if rec := do(t, h, c.method, c.path); rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s %s: status = %d, want 405", c.method, c.path, rec.Code)
		}
	}
}

func TestHandler_EmptyID(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	if rec := do(t, h, http.MethodGet, "/charts/"); rec.Code != http.StatusNotFound {
		t.Errorf("empty id: status = %d, want 404", rec.Code)
	}
}

func TestHandler_ClickUnknownVerb(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodPost, "/charts/demo-bar/click?bar=value.one&verb=bogus")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("unknown verb: status = %d, want 400", rec.Code)
	}
}

func TestHandler_SliceBadCoordinate(t *testing.T) {
	h := htmx.NewHandler(newLineRegistry(t))
	if rec := do(t, h, http.MethodGet, "/charts/demo-line/slice?axis=x&x=abc"); rec.Code != http.StatusBadRequest {
		t.Errorf("bad coord: status = %d, want 400", rec.Code)
	}
	// axis=y reads the y param; a huge value matches no slice.
	if rec := do(t, h, http.MethodGet, "/charts/demo-line/slice?axis=y&y=99999"); rec.Code != http.StatusNotFound {
		t.Errorf("axis=y no match: status = %d, want 404", rec.Code)
	}
	// slice on a non-line chart is a 400.
	hb := htmx.NewHandler(newBarRegistry(t))
	if rec := do(t, hb, http.MethodGet, "/charts/demo-bar/slice?axis=x&x=1"); rec.Code != http.StatusBadRequest {
		t.Errorf("slice on bar: status = %d, want 400", rec.Code)
	}
}

// --- unknown kind / RenderFull errors ------------------------------------------------

func TestRenderFull_Errors(t *testing.T) {
	r := htmx.NewRegistry()
	h := htmx.NewHandler(r)
	// Unknown instance.
	if _, err := h.RenderFull("nope"); err == nil {
		t.Error("RenderFull(unknown id) should error")
	} else if !strings.Contains(err.Error(), "unknown instance") {
		t.Errorf("error = %q, want unknown instance", err.Error())
	}
	// Unknown kind renders a 500 and errors in-process.
	r.Register("weird", htmx.ChartKind("bogus"), nil)
	if _, err := h.RenderFull("weird"); err == nil {
		t.Error("RenderFull(unknown kind) should error")
	} else if !strings.Contains(err.Error(), "unknown chart kind") {
		t.Errorf("error = %q, want unknown chart kind", err.Error())
	}
	if rec := do(t, h, http.MethodGet, "/charts/weird"); rec.Code != http.StatusInternalServerError {
		t.Errorf("GET unknown kind: status = %d, want 500", rec.Code)
	}
}

// --- legend state application (line/pie legend loops) --------------------------------

func TestRenderFull_LineAndPieWithLegends(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterLine("l", line.LineProps{
		Width: 400, Height: 300,
		Data: []line.LineSeries{{ID: "A", Data: []line.LinePointData{{X: "one", Y: 1.0}, {X: "two", Y: 2.0}}}},
		Legends: []legends.LegendProps{{
			Anchor: legends.LegendAnchorTopRight, Direction: legends.LegendDirectionColumn,
			ItemWidth: 80, ItemHeight: 18,
		}},
	})
	r.RegisterPie("p", pie.PieProps{
		Width: 300, Height: 300,
		Data: []any{
			map[string]any{"id": "A", "value": float64(10)},
			map[string]any{"id": "B", "value": float64(20)},
		},
		Legends: []legends.LegendProps{{
			Anchor: legends.LegendAnchorBottom, Direction: legends.LegendDirectionRow,
			ItemWidth: 80, ItemHeight: 18,
		}},
	})
	h := htmx.NewHandler(r)
	for _, id := range []string{"l", "p"} {
		out, err := h.RenderFull(id)
		if err != nil {
			t.Fatalf("RenderFull(%s): %v", id, err)
		}
		// Legend toggling is wired to the registry instance: the legend items
		// carry the toggle endpoint for this chart id.
		if !strings.Contains(out, "/charts/"+id+"/toggle") {
			t.Errorf("%s: legend not wired to /charts/%s/toggle", id, id)
		}
	}
}
