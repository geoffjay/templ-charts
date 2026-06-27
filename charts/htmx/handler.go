package htmx

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/tooltip"
)

// Handler is the http.Handler serving the htmx endpoints. Mount it at any
// prefix; routes are written as absolute paths under /charts/. If you mount
// it at a non-root prefix, strip the prefix before delegating (or use a
// sub-mount with http.StripPrefix).
//
//	Route                       Method  Action
//	/charts/{id}                GET     Full SVG render (hx-aware markup)
//	/charts/{id}/hover          GET     Bar/Pie hover tooltip HTML partial
//	/charts/{id}/slice          GET     Line slice tooltip HTML partial
//	/charts/{id}/click          POST    Bar activation toggle (re-render SVG)
//	/charts/{id}/toggle         POST    Series toggle (re-render SVG)
//
// Unknown instances return 404; unknown verbs return 400.
type Handler struct {
	registry *Registry
}

// NewHandler returns a Handler bound to the given registry.
func NewHandler(r *Registry) *Handler {
	return &Handler{registry: r}
}

// Registry returns the handler's underlying registry (handy for tests / demos
// that want to inspect or mutate state directly).
func (h *Handler) Registry() *Registry { return h.registry }

// RenderFull renders the complete SVG for the instance with the current state
// applied. It is the in-process equivalent of GET /charts/{id}: the page
// handlers use it to render the initial SVG inline without a round-trip.
func (h *Handler) RenderFull(id string) (string, error) {
	inst := h.registry.Get(id)
	if inst == nil {
		return "", errUnknownInstance
	}
	return renderFull(inst)
}

// ServeHTTP dispatches the request to the appropriate endpoint handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// All routes start with "/charts/".
	if !strings.HasPrefix(path, "/charts/") {
		http.NotFound(w, r)
		return
	}
	rest := strings.TrimPrefix(path, "/charts/")
	// Split into id + optional action: "<id>" or "<id>/<action>".
	id, action, _ := strings.Cut(rest, "/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	inst := h.registry.Get(id)
	if inst == nil {
		http.NotFound(w, r)
		return
	}

	switch action {
	case "":
		// Full render: GET only.
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleFull(w, r, inst)
	case "hover":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleHover(w, r, inst)
	case "slice":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleSlice(w, r, inst)
	case "click":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleClick(w, r, inst)
	case "toggle":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleToggle(w, r, inst)
	default:
		http.NotFound(w, r)
	}
}

// handleFull re-renders the complete SVG for the instance and writes it with
// Content-Type image/svg+xml.
func (h *Handler) handleFull(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	out, err := renderFull(inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	_, _ = w.Write([]byte(out))
}

// handleToggle flips the membership of the `series` query param in the
// instance's hidden-id list, then re-renders the full SVG (same body as
// handleFull — the legend swaps the chart container's innerHTML).
func (h *Handler) handleToggle(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	series := r.URL.Query().Get("series")
	if series == "" {
		http.Error(w, "missing series", http.StatusBadRequest)
		return
	}
	inst.toggleHidden(series)
	out, err := renderFull(inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	_, _ = w.Write([]byte(out))
}

// handleClick toggles bar activation. v1 supports verb=activate which simply
// records the hovered key as the active key (no visual change in v1 bar);
// future verbs (select, etc.) can extend this. The response is the full SVG.
func (h *Handler) handleClick(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	if inst.Kind != KindBar {
		http.Error(w, "click not supported for this chart kind", http.StatusBadRequest)
		return
	}
	key := r.URL.Query().Get("bar")
	verb := r.URL.Query().Get("verb")
	if verb == "" {
		verb = "activate"
	}
	switch verb {
	case "activate":
		// Toggle the hovered key: if it's the current hovered key, clear;
		// otherwise set it. This mirrors nivo's bar activation toggle.
		st := inst.State()
		if st.HoveredKey == key {
			inst.setHovered("")
		} else {
			inst.setHovered(key)
		}
	default:
		http.Error(w, "unknown verb: "+verb, http.StatusBadRequest)
		return
	}
	out, err := renderFull(inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	_, _ = w.Write([]byte(out))
}

// handleHover produces an HTML tooltip fragment for the bar or pie arc under
// the cursor. For pie, it also sets ActiveID so subsequent full renders keep
// the arc popped (radius offset). The fragment is swapped into #tooltip-<id>
// by htmx.
func (h *Handler) handleHover(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	q := r.URL.Query()
	switch inst.Kind {
	case KindBar:
		key := q.Get("bar")
		if key == "" {
			http.Error(w, "missing bar", http.StatusBadRequest)
			return
		}
		html, ok := barHoverTooltip(inst, key)
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeHTML(w, html)
	case KindPie:
		arcID := q.Get("arc")
		if arcID == "" {
			http.Error(w, "missing arc", http.StatusBadRequest)
			return
		}
		html, ok := pieHoverTooltip(inst, arcID, true)
		if !ok {
			http.NotFound(w, r)
			return
		}
		writeHTML(w, html)
	default:
		http.Error(w, "hover not supported for this chart kind", http.StatusBadRequest)
	}
}

// handleSlice produces an HTML table-tooltip fragment for a line slice. The
// `axis` param is "x" or "y"; `x` (or `y`) is the slice coordinate (float).
func (h *Handler) handleSlice(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	if inst.Kind != KindLine {
		http.Error(w, "slice not supported for this chart kind", http.StatusBadRequest)
		return
	}
	q := r.URL.Query()
	axis := q.Get("axis")
	if axis == "" {
		axis = "x"
	}
	coordStr := q.Get("x")
	if axis == "y" {
		coordStr = q.Get("y")
	}
	coord, err := strconv.ParseFloat(coordStr, 64)
	if err != nil {
		http.Error(w, "invalid coordinate", http.StatusBadRequest)
		return
	}
	html, ok := lineSliceTooltip(inst, axis, coord)
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeHTML(w, html)
}

// writeHTML writes an HTML fragment with the right content type.
func writeHTML(w http.ResponseWriter, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

// barHoverTooltip finds the bar with the given Key in the instance's computed
// bars and renders a BasicTooltip HTML fragment. Returns ("", false) if no
// bar matches.
func barHoverTooltip(inst *ChartInstance, key string) (string, bool) {
	props := inst.Props.(bar.BarProps)
	applyBarState(&props, inst.ID, inst.State())
	result := bar.UseBar(props)
	for _, b := range result.BarsWithValue {
		if b.Key == key {
			html, err := renderComponent(tooltip.BasicTooltip(tooltip.BasicTooltipProps{
				ID:             result.GetTooltipLabel(b.Data),
				FormattedValue: b.Data.FormattedValue,
				Color:          b.Color,
				EnableChip:     true,
			}))
			if err != nil {
				return "", false
			}
			return html, true
		}
	}
	return "", false
}

// pieHoverTooltip finds the arc with the given ID in the instance's computed
// data and renders a BasicTooltip HTML fragment. When setActive is true, it
// also updates the instance's ActiveID state (so the next full render pops
// the arc). Returns ("", false) if no arc matches.
func pieHoverTooltip(inst *ChartInstance, arcID string, setActive bool) (string, bool) {
	props := inst.Props.(pie.PieProps)
	st := inst.State()
	applyPieState(&props, inst.ID, st)
	result := pie.UsePie(props)
	for _, d := range result.DataWithArc {
		if d.ID == arcID {
			if setActive {
				// Toggle active off when hovering the same arc again is the
				// nivo mouseleave behaviour; v1 sets active on enter and the
				// demo clears it via a separate mouseleave htmx call. For the
				// simple case we just set it.
				inst.setActive(arcID)
			}
			html, err := renderComponent(tooltip.BasicTooltip(tooltip.BasicTooltipProps{
				ID:             d.Label,
				FormattedValue: d.FormattedValue,
				Color:          d.Color,
				EnableChip:     true,
			}))
			if err != nil {
				return "", false
			}
			return html, true
		}
	}
	return "", false
}

// lineSliceTooltip finds the slice at the given axis coordinate in the
// instance's computed slices and renders a TableTooltip HTML fragment with
// one row per series point in the slice. Returns ("", false) if no slice
// matches.
func lineSliceTooltip(inst *ChartInstance, axis string, coord float64) (string, bool) {
	props := inst.Props.(line.LineProps)
	applyLineState(&props, inst.ID, inst.State())
	result := line.UseLine(props)
	// Slices are only computed when EnableSlices is set; otherwise we fall
	// back to the nearest point along the axis.
	slices := result.Slices
	if len(slices) == 0 {
		return "", false
	}
	for _, s := range slices {
		var sCoord float64
		if axis == "y" {
			sCoord = s.Y
		} else {
			sCoord = s.X
		}
		if sameFloat(sCoord, coord) {
			rows := make([]tooltip.TableTooltipRow, 0, len(s.Points))
			for _, p := range s.Points {
				rows = append(rows, tooltip.TableTooltipRow{
					Label: p.SeriesID,
					Value: p.Data.YFormatted,
					Color: p.SeriesColor,
				})
			}
			html, err := renderComponent(tooltip.TableTooltip(tooltip.TableTooltipProps{
				Rows: rows,
			}))
			if err != nil {
				return "", false
			}
			return html, true
		}
	}
	return "", false
}

// sameFloat compares two floats with a small tolerance (the slice coord is
// formatted with %g in the hx-get URL, so round-trip may differ in the last
// digit).
func sameFloat(a, b float64) bool {
	return absFloat(a-b) < 1e-6
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
