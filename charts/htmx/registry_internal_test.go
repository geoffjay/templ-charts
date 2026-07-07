package htmx

import "testing"

// Internal tests for the unexported ChartInstance state mutators that the
// HTTP endpoints don't exercise directly.

func TestSetHiddenReplacesList(t *testing.T) {
	inst := &ChartInstance{ID: "x"}
	inst.setHidden([]string{"a", "b"})
	st := inst.State()
	if len(st.HiddenIDs) != 2 || st.HiddenIDs[0] != "a" || st.HiddenIDs[1] != "b" {
		t.Fatalf("hidden = %v, want [a b]", st.HiddenIDs)
	}
	// Replaces, not appends.
	inst.setHidden([]string{"c"})
	st = inst.State()
	if len(st.HiddenIDs) != 1 || st.HiddenIDs[0] != "c" {
		t.Fatalf("hidden after replace = %v, want [c]", st.HiddenIDs)
	}
	// The stored slice is a copy of the input.
	src := []string{"d"}
	inst.setHidden(src)
	src[0] = "mutated"
	if got := inst.State().HiddenIDs[0]; got != "d" {
		t.Errorf("setHidden should copy its input, got %q", got)
	}
}

func TestClearHoverResetsAllHoverState(t *testing.T) {
	inst := &ChartInstance{ID: "x"}
	inst.setHovered("bar.key")
	inst.setHoverXY(10, 20)
	inst.clearHover()
	st := inst.State()
	if st.HoveredKey != "" || st.HoverX != 0 || st.HoverY != 0 || st.HasHover {
		t.Errorf("clearHover left state %+v", st)
	}
}

func TestChartClass(t *testing.T) {
	if got := chartClass(""); got != "tc-chart" {
		t.Errorf("chartClass(\"\") = %q", got)
	}
	if got := chartClass("card wide"); got != "tc-chart card wide" {
		t.Errorf("chartClass(extra) = %q", got)
	}
}
