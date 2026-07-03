package htmx

// MountProps configures Mount, the container that wires a server-rendered
// chart SVG into the htmx interactivity layer.
type MountProps struct {
	// ID is the registered chart instance id — the same id passed to
	// Registry.Register and used by the /charts/{id} routes. Required; it is
	// what couples the container (#chart-<ID>), the tooltip target
	// (#tooltip-<ID>), and the per-element hover wiring emitted inside the SVG.
	ID string

	// SVG is the pre-rendered chart markup to inject, typically obtained from
	// Handler.RenderFull(ID) (or by rendering the chart component directly for
	// a static embed). Required.
	SVG string

	// Interactive wires the hover/tooltip plumbing: a mouseleave reset on the
	// chart container plus the sibling tooltip swap target. It should match the
	// Interactive state of the registered instance. Leave false for a static
	// (zero-JS) embed.
	Interactive bool

	// Class is an optional extra CSS class appended to the chart container, so
	// callers can attach their own styling (card chrome, sizing, layout)
	// without the library shipping any opinionated CSS.
	Class string
}

// chartClass returns the chart container's class list: the stable library hook
// "tc-chart", with any caller-supplied class appended.
func chartClass(extra string) string {
	if extra == "" {
		return "tc-chart"
	}
	return "tc-chart " + extra
}
