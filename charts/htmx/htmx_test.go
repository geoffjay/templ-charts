package htmx_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar"
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/treemap"
)

// newBarRegistry builds a registry with one interactive bar instance for the
// hover/toggle tests.
func newBarRegistry(t *testing.T) *htmx.Registry {
	t.Helper()
	r := htmx.NewRegistry()
	r.RegisterBar("demo-bar", bar.BarProps{
		Width: 500, Height: 300,
		Keys:    []string{"value"},
		IndexBy: "id",
		Data: []bar.BarDatum{
			{"id": "one", "value": float64(10)},
			{"id": "two", "value": float64(20)},
			{"id": "three", "value": float64(30)},
		},
		Legends: []bar.BarLegendProps{
			{LegendProps: legends.LegendProps{
				Anchor:    legends.LegendAnchorTopRight,
				Direction: legends.LegendDirectionColumn,
				ItemWidth: 100, ItemHeight: 20,
			}},
		},
	})
	return r
}

func newLineRegistry(t *testing.T) *htmx.Registry {
	t.Helper()
	r := htmx.NewRegistry()
	r.RegisterLine("demo-line", line.LineProps{
		Width: 500, Height: 300,
		EnableSlices: line.EnableSlicesX,
		Data: []line.LineSeries{
			{ID: "A", Data: []line.LinePointData{
				{X: "one", Y: float64(10)},
				{X: "two", Y: float64(20)},
				{X: "three", Y: float64(30)},
			}},
			{ID: "B", Data: []line.LinePointData{
				{X: "one", Y: float64(5)},
				{X: "two", Y: float64(15)},
				{X: "three", Y: float64(25)},
			}},
		},
	})
	return r
}

func newPieRegistry(t *testing.T) *htmx.Registry {
	t.Helper()
	r := htmx.NewRegistry()
	r.RegisterPie("demo-pie", pie.PieProps{
		Width: 400, Height: 400,
		Data: []any{
			map[string]any{"id": "A", "value": float64(10)},
			map[string]any{"id": "B", "value": float64(20)},
			map[string]any{"id": "C", "value": float64(30)},
		},
	})
	return r
}

func pf(v float64) *float64 { return &v }

func newHeatmapRegistry(t *testing.T) *htmx.Registry {
	t.Helper()
	r := htmx.NewRegistry()
	r.RegisterHeatmap("demo-heatmap", heatmap.HeatMapProps{
		Width: 400, Height: 300,
		Data: []heatmap.HeatMapSerie{
			{ID: "Japan", Data: []heatmap.HeatMapDatum{{X: "Train", Y: pf(10)}, {X: "Car", Y: pf(20)}}},
			{ID: "USA", Data: []heatmap.HeatMapDatum{{X: "Train", Y: pf(30)}, {X: "Car", Y: nil}}},
		},
	})
	return r
}

// do runs a request against the handler and returns the recorded response.
func do(t *testing.T, h http.Handler, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := htmx.NewRegistry()
	if got := r.Get("nope"); got != nil {
		t.Fatalf("Get on empty registry returned %v, want nil", got)
	}
	inst := r.RegisterBar("b", bar.BarProps{Width: 10, Height: 10})
	if inst == nil || inst.ID != "b" {
		t.Fatalf("RegisterBar returned %v, want instance with ID b", inst)
	}
	if r.Get("b") != inst {
		t.Errorf("Get(b) did not return the registered instance")
	}
	if kind := r.Get("b").Kind; kind != htmx.KindBar {
		t.Errorf("Kind = %q, want %q", kind, htmx.KindBar)
	}
}

func TestRegistry_RegisterReplacesState(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterBar("b", bar.BarProps{})
	inst := r.Get("b")
	inst.SetStateForTest(htmx.State{HoveredKey: "x"})

	r.RegisterBar("b", bar.BarProps{})
	inst2 := r.Get("b")
	if s := inst2.State(); s.HoveredKey != "" {
		t.Fatalf("re-register should reset state; got %+v", s)
	}
}

func TestRegistry_IDsSorted(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterBar("z", bar.BarProps{})
	r.RegisterBar("a", bar.BarProps{})
	r.RegisterBar("m", bar.BarProps{})
	got := r.IDs()
	want := []string{"a", "m", "z"}
	if len(got) != len(want) {
		t.Fatalf("IDs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("IDs[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestRegistry_Unregister(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterBar("b", bar.BarProps{})
	r.Unregister("b")
	if r.Get("b") != nil {
		t.Errorf("Unregister did not remove the instance")
	}
	// no-op on absent id
	r.Unregister("nope")
}

func TestHandler_FullRender_Bar(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-bar")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "image/svg+xml") {
		t.Errorf("Content-Type = %q, want image/svg+xml", ct)
	}
	body := rec.Body.String()
	if !strings.HasPrefix(body, "<svg") {
		t.Errorf("body should start with <svg, got %q", body[:min(40, len(body))])
	}
	// Interactive render emits hx-* attrs because ChartID is set. templ
	// HTML-encodes the attribute value, so the quote chars become &#34;.
	if !strings.Contains(body, `/charts/demo-bar/hover?bar=value.`) {
		t.Errorf("body missing hx-get hover attrs")
	}
}

func TestHandler_FullRender_NotFound(t *testing.T) {
	h := htmx.NewHandler(htmx.NewRegistry())
	rec := do(t, h, http.MethodGet, "/charts/nope")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandler_FullRender_MethodNotAllowed(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodPost, "/charts/demo-bar")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}

func TestHandler_BarHover(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-bar/hover?bar=value.one")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "nivo-tooltip-basic") {
		t.Errorf("hover body missing tooltip class; got %q", body)
	}
	// BasicTooltip renders "id: formattedValue".
	if !strings.Contains(body, "value.one") && !strings.Contains(body, "one") {
		t.Errorf("hover body missing label; got %q", body)
	}
}

func TestHandler_BarHover_NotFound(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-bar/hover?bar=does.not.exist")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandler_BarHover_MissingParam(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-bar/hover")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHandler_Toggle_HidesSeries(t *testing.T) {
	r := newBarRegistry(t)
	h := htmx.NewHandler(r)

	// Before toggle, the "value" series renders bars.
	rec0 := do(t, h, http.MethodGet, "/charts/demo-bar")
	before := rec0.Body.String()
	// 3 indexes × 1 key = 3 bars.
	if got := strings.Count(before, "<rect"); got < 4 { // +1 background rect
		t.Fatalf("expected at least 4 rects (3 bars + bg), got %d", got)
	}

	// Toggle "value" off: re-render should have no bars with that key.
	rec := do(t, h, http.MethodPost, "/charts/demo-bar/toggle?series=value")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	after := rec.Body.String()
	if strings.Contains(after, `/charts/demo-bar/hover?bar=value.`) {
		t.Errorf("after toggling 'value' off, body should not emit hover attrs for value bars")
	}
	// State updated.
	if s := r.Get("demo-bar").State(); !contains(s.HiddenIDs, "value") {
		t.Errorf("after toggle, hidden ids = %v, want to contain 'value'", s.HiddenIDs)
	}

	// Toggle back on.
	rec2 := do(t, h, http.MethodPost, "/charts/demo-bar/toggle?series=value")
	after2 := rec2.Body.String()
	if !strings.Contains(after2, `/charts/demo-bar/hover?bar=value.`) {
		t.Errorf("after toggling 'value' back on, body should emit hover attrs for value bars")
	}
	if s := r.Get("demo-bar").State(); contains(s.HiddenIDs, "value") {
		t.Errorf("after second toggle, hidden ids = %v, want empty", s.HiddenIDs)
	}
}

func TestHandler_Toggle_MissingSeries(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodPost, "/charts/demo-bar/toggle")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHandler_Click_TogglesHoveredKey(t *testing.T) {
	r := newBarRegistry(t)
	h := htmx.NewHandler(r)
	rec := do(t, h, http.MethodPost, "/charts/demo-bar/click?bar=value.one&verb=activate")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if s := r.Get("demo-bar").State(); s.HoveredKey != "value.one" {
		t.Errorf("after click, HoveredKey = %q, want 'value.one'", s.HoveredKey)
	}
	// Click again toggles off.
	rec2 := do(t, h, http.MethodPost, "/charts/demo-bar/click?bar=value.one&verb=activate")
	if s := r.Get("demo-bar").State(); s.HoveredKey != "" {
		t.Errorf("after second click, HoveredKey = %q, want empty", s.HoveredKey)
	}
	_ = rec2
}

func TestHandler_Click_UnsupportedKind(t *testing.T) {
	h := htmx.NewHandler(newLineRegistry(t))
	rec := do(t, h, http.MethodPost, "/charts/demo-line/click?bar=x")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHandler_PieHover(t *testing.T) {
	r := newPieRegistry(t)
	h := htmx.NewHandler(r)
	rec := do(t, h, http.MethodGet, "/charts/demo-pie/hover?arc=A")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	body := rec.Body.String()
	if !strings.Contains(body, "nivo-tooltip-basic") {
		t.Errorf("pie hover body missing tooltip class; got %q", body)
	}
	// ActiveID set so the next full render pops the arc.
	if s := r.Get("demo-pie").State(); s.ActiveID != "A" {
		t.Errorf("after hover, ActiveID = %q, want 'A'", s.ActiveID)
	}
}

func TestHandler_PieHover_NotFound(t *testing.T) {
	h := htmx.NewHandler(newPieRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-pie/hover?arc=ZZ")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandler_PieToggle(t *testing.T) {
	r := newPieRegistry(t)
	h := htmx.NewHandler(r)
	rec := do(t, h, http.MethodPost, "/charts/demo-pie/toggle?series=A")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if !strings.HasPrefix(rec.Body.String(), "<svg") {
		t.Errorf("pie toggle body should be full SVG, got %q", rec.Body.String()[:min(40, len(rec.Body.String()))])
	}
	if s := r.Get("demo-pie").State(); !contains(s.HiddenIDs, "A") {
		t.Errorf("after toggle, hidden = %v, want to contain 'A'", s.HiddenIDs)
	}
}

func TestHandler_LineSlice(t *testing.T) {
	h := htmx.NewHandler(newLineRegistry(t))
	// First find a slice coordinate by probing a full render's slices layer.
	// We can compute it: slices are per-unique-x. The first slice's x is the
	// x-scale output for "one". Rather than hardcode, use the rendered slice
	// rect's hx-get URL.
	rec0 := do(t, h, http.MethodGet, "/charts/demo-line")
	body := rec0.Body.String()
	// templ HTML-encodes the attribute value, so & becomes &amp;.
	search := `/slice?axis=x&amp;x=`
	idx := strings.Index(body, search)
	if idx < 0 {
		// fall back: find the slice path with a raw & (some templ versions
		// may not encode inside attribute values).
		search = `/slice?axis=x&x=`
		idx = strings.Index(body, search)
	}
	if idx < 0 {
		t.Fatalf("full line render missing slice hx-get; has slice=%v", strings.Contains(body, "slice"))
	}
	// Extract the x value: between "x=" and the next quote.
	rest := body[idx+len(search):]
	end := strings.IndexByte(rest, '"')
	if end < 0 {
		t.Fatalf("could not find end of x value in %q", rest[:min(40, len(rest))])
	}
	xVal := rest[:end]

	rec := do(t, h, http.MethodGet, "/charts/demo-line/slice?axis=x&x="+xVal)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", ct)
	}
	body2 := rec.Body.String()
	if !strings.Contains(body2, "nivo-tooltip-table") {
		t.Errorf("slice tooltip body missing table class; got %q", body2)
	}
}

func TestHandler_LineSlice_NotFound(t *testing.T) {
	h := htmx.NewHandler(newLineRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-line/slice?axis=x&x=9999")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandler_UnknownAction_NotFound(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-bar/bogus")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandler_RootPath_NotFound(t *testing.T) {
	h := htmx.NewHandler(newBarRegistry(t))
	rec := do(t, h, http.MethodGet, "/other")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

// --- helpers ---

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

func TestHeatmapFullRenderEmitsHoverAttrs(t *testing.T) {
	h := htmx.NewHandler(newHeatmapRegistry(t))
	out, err := h.RenderFull("demo-heatmap")
	if err != nil {
		t.Fatalf("RenderFull: %v", err)
	}
	// heatmap hover is routed through the client interactivity layer
	// (charts/interact): the 3 data cells emit a data-tc-tooltip, the nil cell
	// (USA.Car) does not, and no per-cell server round-trip is emitted.
	if got := strings.Count(out, "data-tc-tooltip"); got != 3 {
		t.Errorf("expected 3 hoverable data cells (data-tc-tooltip), got %d", got)
	}
	if strings.Contains(out, "hover?cell=") {
		t.Errorf("heatmap hover should be client-side now, not a server round-trip")
	}
}

func TestHeatmapHoverReturnsTooltipAndOOB(t *testing.T) {
	h := htmx.NewHandler(newHeatmapRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-heatmap/hover?cell=Japan.Train")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Japan - Train") {
		t.Errorf("tooltip missing the cell label %q", "Japan - Train")
	}
	if !strings.Contains(body, `hx-swap-oob="innerHTML:#chart-demo-heatmap"`) {
		t.Errorf("response missing the OOB chart swap")
	}
	if !strings.Contains(body, "<svg") {
		t.Errorf("OOB swap missing the re-rendered svg")
	}
}

func TestHeatmapHoverMissingCellIs400(t *testing.T) {
	h := htmx.NewHandler(newHeatmapRegistry(t))
	rec := do(t, h, http.MethodGet, "/charts/demo-heatmap/hover")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for missing cell", rec.Code)
	}
}

// --- hierarchy zoom (icicle/treemap/circle-packing/sunburst) ---

func icicleData() icicle.IcicleNode {
	return icicle.IcicleNode{ID: "root", Children: []icicle.IcicleNode{
		{ID: "A", Children: []icicle.IcicleNode{{ID: "a1", Value: 8}, {ID: "a2", Value: 4}}},
		{ID: "B", Children: []icicle.IcicleNode{{ID: "b1", Value: 6}}},
		{ID: "C", Value: 10},
	}}
}

// assertZoomRoundTrip drives the shared zoom handler contract: focusing a child
// changes the SVG, and zooming back to root restores the unfocused render.
func assertZoomRoundTrip(t *testing.T, h *htmx.Handler, id, childID string) {
	t.Helper()
	base, err := h.RenderFull(id)
	if err != nil {
		t.Fatalf("RenderFull(%s): %v", id, err)
	}
	rec := do(t, h, http.MethodGet, "/charts/"+id+"/zoom?node="+childID)
	if rec.Code != http.StatusOK {
		t.Fatalf("zoom status = %d, want 200", rec.Code)
	}
	focused := rec.Body.String()
	if !strings.Contains(focused, "<svg") {
		t.Errorf("zoom body should be a full SVG")
	}
	if focused == base {
		t.Errorf("focusing %q did not change the render", childID)
	}
	if inst := h.Registry().Get(id); inst == nil || inst.State().FocusID != childID {
		t.Errorf("focus state = %q, want %q", h.Registry().Get(id).State().FocusID, childID)
	}
	// Re-clicking the focused node zooms out to its parent — for a depth-1 node
	// that is the root (full view), restoring the original render.
	rec2 := do(t, h, http.MethodGet, "/charts/"+id+"/zoom?node="+childID)
	if rec2.Code != http.StatusOK {
		t.Fatalf("zoom-out status = %d, want 200", rec2.Code)
	}
	if got := rec2.Body.String(); got != base {
		t.Errorf("zooming back to root did not restore the unfocused render")
	}
}

func TestHandler_IcicleZoom(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterIcicle("demo-icicle", icicle.IcicleProps{
		Width: 500, Height: 300, EnableZooming: true, Data: icicleData(),
	})
	assertZoomRoundTrip(t, htmx.NewHandler(r), "demo-icicle", "A")
}

func TestHandler_TreemapZoom(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterTreemap("demo-treemap", treemap.TreemapProps{
		Width: 500, Height: 400, EnableZooming: true,
		Data: treemap.TreemapNode{ID: "root", Children: []treemap.TreemapNode{
			{ID: "A", Children: []treemap.TreemapNode{{ID: "a1", Value: 12}, {ID: "a2", Value: 8}}},
			{ID: "B", Children: []treemap.TreemapNode{{ID: "b1", Value: 10}}},
			{ID: "C", Value: 15},
		}},
	})
	assertZoomRoundTrip(t, htmx.NewHandler(r), "demo-treemap", "A")
}

func TestHandler_CirclePackingZoom(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterCirclePacking("demo-cp", cp.CirclePackingProps{
		Width: 400, Height: 400, EnableZooming: true,
		Data: cp.CirclePackingNode{ID: "root", Children: []cp.CirclePackingNode{
			{ID: "A", Children: []cp.CirclePackingNode{{ID: "a1", Value: 8}, {ID: "a2", Value: 4}}},
			{ID: "B", Children: []cp.CirclePackingNode{{ID: "b1", Value: 6}}},
			{ID: "C", Value: 10},
		}},
	})
	assertZoomRoundTrip(t, htmx.NewHandler(r), "demo-cp", "A")
}

func TestHandler_SunburstZoom(t *testing.T) {
	r := htmx.NewRegistry()
	r.RegisterSunburst("demo-sb", sunburst.SunburstProps{
		Width: 400, Height: 400, EnableZooming: true,
		Data: sunburst.SunburstNode{ID: "root", Children: []sunburst.SunburstNode{
			{ID: "A", Children: []sunburst.SunburstNode{{ID: "a1", Value: 8}, {ID: "a2", Value: 4}}},
			{ID: "B", Children: []sunburst.SunburstNode{{ID: "b1", Value: 6}}},
			{ID: "C", Value: 10},
		}},
	})
	assertZoomRoundTrip(t, htmx.NewHandler(r), "demo-sb", "A")
}

func TestHandler_ZoomUnsupportedKind(t *testing.T) {
	rec := do(t, htmx.NewHandler(newBarRegistry(t)), http.MethodGet, "/charts/demo-bar/zoom?node=x")
	if rec.Code != http.StatusBadRequest {
		t.Errorf("zoom on a bar should be 400, got %d", rec.Code)
	}
}
