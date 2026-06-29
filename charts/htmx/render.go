package htmx

import (
	"context"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
)

// renderFull renders the full SVG for the instance, applying the current
// state (hidden ids, pie active id) to a clone of the props snapshot. The
// returned string is a complete <svg>…</svg> document.
func renderFull(inst *ChartInstance) (string, error) {
	st := inst.State()
	switch inst.Kind {
	case KindBar:
		props := inst.Props.(bar.BarProps)
		applyBarState(&props, inst.ID, st)
		return renderComponent(bar.Bar(props))
	case KindLine:
		props := inst.Props.(line.LineProps)
		applyLineState(&props, inst.ID, st)
		return renderComponent(line.Line(props))
	case KindPie:
		props := inst.Props.(pie.PieProps)
		applyPieState(&props, inst.ID, st)
		return renderComponent(pie.Pie(props))
	}
	return "", errUnknownKind
}

// applyBarState overlays the per-instance state onto a bar.BarProps clone.
func applyBarState(props *bar.BarProps, id string, st State) {
	props.ChartID = id
	props.IsInteractive = true
	props.InitialHiddenIDs = st.HiddenIDs
	props.HoveredKey = st.HoveredKey
	if props.Legends != nil {
		for i := range props.Legends {
			props.Legends[i].ChartID = id
			if props.Legends[i].Toggle == "" {
				props.Legends[i].Toggle = "1"
			}
		}
	}
}

// applyLineState overlays the per-instance state onto a line.LineProps clone.
func applyLineState(props *line.LineProps, id string, st State) {
	props.ChartID = id
	props.IsInteractive = true
	props.InitialHiddenIDs = st.HiddenIDs
	props.HoverX = st.HoverX
	props.HoverY = st.HoverY
	props.HasHover = st.HasHover
	if props.Legends != nil {
		for i := range props.Legends {
			props.Legends[i].ChartID = id
			if props.Legends[i].Toggle == "" {
				props.Legends[i].Toggle = "1"
			}
		}
	}
}

// applyPieState overlays the per-instance state onto a pie.PieProps clone.
func applyPieState(props *pie.PieProps, id string, st State) {
	props.ChartID = id
	props.IsInteractive = true
	props.InitialHiddenIDs = st.HiddenIDs
	props.ActiveID = st.ActiveID
	if props.Legends != nil {
		for i := range props.Legends {
			props.Legends[i].ChartID = id
			if props.Legends[i].Toggle == "" {
				props.Legends[i].Toggle = "1"
			}
		}
	}
}

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) (string, error) {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

// errUnknownKind is returned by renderFull when the instance Kind is not one
// of the supported chart kinds.
var errUnknownKind = htmxError("htmx: unknown chart kind")

// errUnknownInstance is returned when an instance id is not registered.
var errUnknownInstance = htmxError("htmx: unknown instance")

// htmxError is a simple string error type.
type htmxError string

func (e htmxError) Error() string { return string(e) }
