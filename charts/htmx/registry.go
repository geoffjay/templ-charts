// Package htmx implements the server-side interactivity layer for templ-charts.
//
// It mirrors an island pattern: a stateful Registry of ChartInstance entries
// keyed by instance id, and an http.Handler
// (Handler) exposing the endpoints the chart components already emit on their
// hx-* attributes:
//
//   - GET  /charts/{id}             full SVG render (hx-target: #chart-<id>)
//   - GET  /charts/{id}/hover?bar=k  bar hover tooltip HTML partial
//   - GET  /charts/{id}/hover?arc=i  pie hover tooltip HTML partial + active arc
//   - GET  /charts/{id}/slice?axis=x&x=v  line slice crosshair + tooltip partial
//   - POST /charts/{id}/click?bar=k&verb=activate  bar activation toggle
//   - POST /charts/{id}/toggle?series=id   re-render full SVG with series toggled
//
// State lives in the Registry; the documented trade-off is that this is fine
// for demos / small apps and would be moved to a session/cookie store for
// horizontal scaling.
package htmx

import (
	"sync"

	"github.com/geoffjay/templ-charts/charts/bar"
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/treemap"
)

// ChartKind enumerates the chart types the registry can serve.
type ChartKind string

const (
	KindBar        ChartKind = "bar"
	KindLine       ChartKind = "line"
	KindPie        ChartKind = "pie"
	KindHeatmap    ChartKind = "heatmap"
	KindIcicle     ChartKind = "icicle"
	KindTreemap    ChartKind = "treemap"
	KindCirclePack ChartKind = "circlepacking"
	KindSunburst   ChartKind = "sunburst"
)

// ChartInstance is one registered chart: its kind, an immutable props
// snapshot, and the mutable per-instance state mutated by the endpoints.
//
// The Props field is stored as any and resolved by the renderer via a type
// switch on ChartKind (bar.BarProps / line.LineProps / pie.PieProps). Storing
// the concrete type keeps the snapshot cheap to clone on each request.
type ChartInstance struct {
	ID    string
	Kind  ChartKind
	Props any // bar.BarProps | line.LineProps | pie.PieProps

	mu         sync.RWMutex
	hiddenIDs  []string
	hoveredKey string  // bar/click activation focus (bar key)
	activeID   string  // pie active arc id (drives radius offset)
	hoverX     float64 // line mesh hover x (chart units)
	hoverY     float64 // line mesh hover y (chart units)
	hasHover   bool    // line mesh: a hover is active (gates crosshair render)
	focusID    string  // hierarchy zoom focus node id ("" = root / full view)
}

// State is a snapshot of the mutable per-instance state. Returned by
// Instance.State so callers (and tests) can inspect without holding the lock.
type State struct {
	HiddenIDs  []string
	HoveredKey string  // bar hovered key (drives bar active highlight)
	ActiveID   string  // pie active arc id (drives radius offset)
	HoverX     float64 // line mesh hover x (chart units, drives crosshair)
	HoverY     float64 // line mesh hover y (chart units, drives crosshair)
	HasHover   bool    // line mesh: a hover is active (gates crosshair render)
	FocusID    string  // hierarchy zoom focus node id ("" = root / full view)
}

// State returns a copy of the instance's mutable state.
func (c *ChartInstance) State() State {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return State{
		HiddenIDs:  append([]string(nil), c.hiddenIDs...),
		HoveredKey: c.hoveredKey,
		ActiveID:   c.activeID,
		HoverX:     c.hoverX,
		HoverY:     c.hoverY,
		HasHover:   c.hasHover,
		FocusID:    c.focusID,
	}
}

// setFocus sets the hierarchy zoom focus node id ("" clears to the full view).
func (c *ChartInstance) setFocus(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.focusID = id
}

// setHidden replaces the hidden-id list (toggle endpoint).
func (c *ChartInstance) setHidden(ids []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hiddenIDs = append([]string(nil), ids...)
}

// toggleHidden flips the membership of id in the hidden list and returns the
// new hidden list. Mirrors nivo's series-toggle behavior.
func (c *ChartInstance) toggleHidden(id string) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i, h := range c.hiddenIDs {
		if h == id {
			c.hiddenIDs = append(c.hiddenIDs[:i], c.hiddenIDs[i+1:]...)
			return append([]string(nil), c.hiddenIDs...)
		}
	}
	c.hiddenIDs = append(c.hiddenIDs, id)
	return append([]string(nil), c.hiddenIDs...)
}

// setActive sets the pie active-arc id ("" clears).
func (c *ChartInstance) setActive(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.activeID = id
}

// setHovered sets the bar hovered key (used by click activation focus).
func (c *ChartInstance) setHovered(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hoveredKey = key
}

// setHoverXY sets the line mesh hover coordinates (chart units).
func (c *ChartInstance) setHoverXY(x, y float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hoverX = x
	c.hoverY = y
	c.hasHover = true
}

// clearHover resets all hover state (bar hovered key + line hover coords).
func (c *ChartInstance) clearHover() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hoveredKey = ""
	c.hoverX = 0
	c.hoverY = 0
	c.hasHover = false
}

// Registry maps instance IDs to *ChartInstance. The zero value is not usable;
// use NewRegistry.
type Registry struct {
	mu        sync.RWMutex
	instances map[string]*ChartInstance
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{instances: map[string]*ChartInstance{}}
}

// Register stores a chart instance under id. The props snapshot is taken as-is
// (callers should pass a fully-populated bar.BarProps / line.LineProps /
// pie.PieProps). Returns the created *ChartInstance so callers can seed
// InitialHiddenIDs etc.
//
// Register is idempotent: re-registering the same id replaces the previous
// instance (and resets its state).
func (r *Registry) Register(id string, kind ChartKind, props any) *ChartInstance {
	inst := &ChartInstance{ID: id, Kind: kind, Props: props}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.instances[id] = inst
	return inst
}

// RegisterBar is a convenience wrapper for Register(id, KindBar, props).
func (r *Registry) RegisterBar(id string, props bar.BarProps) *ChartInstance {
	return r.Register(id, KindBar, props)
}

// RegisterLine is a convenience wrapper for Register(id, KindLine, props).
func (r *Registry) RegisterLine(id string, props line.LineProps) *ChartInstance {
	return r.Register(id, KindLine, props)
}

// RegisterPie is a convenience wrapper for Register(id, KindPie, props).
func (r *Registry) RegisterPie(id string, props pie.PieProps) *ChartInstance {
	return r.Register(id, KindPie, props)
}

// RegisterHeatmap is a convenience wrapper for Register(id, KindHeatmap, props).
func (r *Registry) RegisterHeatmap(id string, props heatmap.HeatMapProps) *ChartInstance {
	return r.Register(id, KindHeatmap, props)
}

// RegisterIcicle is a convenience wrapper for Register(id, KindIcicle, props).
func (r *Registry) RegisterIcicle(id string, props icicle.IcicleProps) *ChartInstance {
	return r.Register(id, KindIcicle, props)
}

// RegisterTreemap is a convenience wrapper for Register(id, KindTreemap, props).
func (r *Registry) RegisterTreemap(id string, props treemap.TreemapProps) *ChartInstance {
	return r.Register(id, KindTreemap, props)
}

// RegisterCirclePacking is a convenience wrapper for
// Register(id, KindCirclePack, props).
func (r *Registry) RegisterCirclePacking(id string, props cp.CirclePackingProps) *ChartInstance {
	return r.Register(id, KindCirclePack, props)
}

// RegisterSunburst is a convenience wrapper for Register(id, KindSunburst, props).
func (r *Registry) RegisterSunburst(id string, props sunburst.SunburstProps) *ChartInstance {
	return r.Register(id, KindSunburst, props)
}

// Get returns the instance for id, or nil if not registered.
func (r *Registry) Get(id string) *ChartInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.instances[id]
}

// Unregister removes an instance (no-op if absent).
func (r *Registry) Unregister(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.instances, id)
}

// IDs returns the sorted list of registered instance IDs (handy for demos).
func (r *Registry) IDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.instances))
	for id := range r.instances {
		out = append(out, id)
	}
	// simple insertion sort (registry is small)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// SetStateForTest overwrites the mutable state in one shot. It is intended
// for tests that need to seed state without driving it through the HTTP
// endpoints.
func (c *ChartInstance) SetStateForTest(s State) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hiddenIDs = append([]string(nil), s.HiddenIDs...)
	c.hoveredKey = s.HoveredKey
	c.activeID = s.ActiveID
	c.hoverX = s.HoverX
	c.hoverY = s.HoverY
	c.hasHover = s.HasHover
	c.focusID = s.FocusID
}
