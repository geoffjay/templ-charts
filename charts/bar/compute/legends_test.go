package compute_test

import (
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar/compute"
)

// legendBars builds a small bar list: keys k1/k2 across indexes A/B, in the
// generation order (all k1 bars, then all k2 bars).
func legendBars() []compute.ComputedBarDatum {
	bars := []compute.ComputedBarDatum{}
	for _, key := range []string{"k1", "k2"} {
		for _, idx := range []string{"A", "B"} {
			bars = append(bars, compute.ComputedBarDatum{
				Key:   key + "." + idx,
				Color: "c-" + key,
				Data: compute.ComputedDatum{
					ID:         key,
					IndexValue: idx,
					Value:      1.0,
					Data:       map[string]any{"name": "n-" + idx},
				},
			})
		}
	}
	return bars
}

func idLabel(d map[string]any) string { return d["id"].(string) }

func indexLabel(d map[string]any) string { return d["indexValue"].(string) }

func legendIDs(items []compute.LegendData) []string {
	ids := make([]string, len(items))
	for i, it := range items {
		ids[i] = it.ID
	}
	return ids
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetLegendData_FromKeys(t *testing.T) {
	cases := []struct {
		name      string
		direction string
		groupMode string
		layout    string
		reverse   bool
		wantIDs   []string
	}{
		// vertical + stacked + column + !reverse → reversed so the legend
		// order matches the visual stack (topmost first).
		{"vertical-stacked-column", "column", "stacked", "vertical", false, []string{"k2", "k1"}},
		{"vertical-stacked-column-reverse", "column", "stacked", "vertical", true, []string{"k1", "k2"}},
		{"vertical-stacked-row", "row", "stacked", "vertical", false, []string{"k1", "k2"}},
		{"vertical-grouped-column", "column", "grouped", "vertical", false, []string{"k1", "k2"}},
		// horizontal + stacked + reverse → reversed.
		{"horizontal-stacked-reverse", "column", "stacked", "horizontal", true, []string{"k2", "k1"}},
		{"horizontal-stacked", "column", "stacked", "horizontal", false, []string{"k1", "k2"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := compute.GetLegendData(legendBars(), "keys", c.direction, c.groupMode, c.layout, c.reverse, idLabel)
			if !equalStrings(legendIDs(got), c.wantIDs) {
				t.Errorf("ids = %v, want %v", legendIDs(got), c.wantIDs)
			}
			// One entry per unique key, carrying the bar color.
			if len(got) != 2 {
				t.Fatalf("len = %d, want 2", len(got))
			}
			for _, it := range got {
				if it.Color != "c-"+it.ID {
					t.Errorf("color for %s = %q, want %q", it.ID, it.Color, "c-"+it.ID)
				}
				if it.Label != it.ID {
					t.Errorf("label = %q, want %q", it.Label, it.ID)
				}
				if it.Hidden {
					t.Errorf("entry %s should not be hidden", it.ID)
				}
			}
		})
	}
}

func TestGetLegendData_FromIndexes(t *testing.T) {
	// One entry per unique index value.
	got := compute.GetLegendData(legendBars(), "indexes", "column", "grouped", "vertical", false, indexLabel)
	if !equalStrings(legendIDs(got), []string{"A", "B"}) {
		t.Errorf("ids = %v, want [A B]", legendIDs(got))
	}
	// Horizontal layout reverses index legends.
	got = compute.GetLegendData(legendBars(), "indexes", "column", "grouped", "horizontal", false, indexLabel)
	if !equalStrings(legendIDs(got), []string{"B", "A"}) {
		t.Errorf("horizontal ids = %v, want [B A]", legendIDs(got))
	}
}

func TestGetLegendData_DefaultColorAndHidden(t *testing.T) {
	bars := []compute.ComputedBarDatum{
		{Key: "k1.A", Color: "", Data: compute.ComputedDatum{ID: "k1", IndexValue: "A", Hidden: true}},
	}
	got := compute.GetLegendData(bars, "keys", "row", "grouped", "vertical", false, idLabel)
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].Color != "#000" {
		t.Errorf("empty color should default to #000, got %q", got[0].Color)
	}
	if !got[0].Hidden {
		t.Errorf("hidden flag should be carried through")
	}
	// Same for indexes.
	got = compute.GetLegendData(bars, "indexes", "row", "grouped", "vertical", false, indexLabel)
	if got[0].Color != "#000" || !got[0].Hidden {
		t.Errorf("indexes entry = %+v, want #000/hidden", got[0])
	}
}

func TestGetLegendData_LabelAccessorSeesDatumFields(t *testing.T) {
	// A custom accessor can reach the nested raw datum under "data".
	label := func(d map[string]any) string {
		return d["data"].(map[string]any)["name"].(string)
	}
	got := compute.GetLegendData(legendBars(), "indexes", "column", "grouped", "vertical", false, label)
	if got[0].Label != "n-A" || got[1].Label != "n-B" {
		t.Errorf("labels = %q/%q, want n-A/n-B", got[0].Label, got[1].Label)
	}
}

func TestGetLegendData_Empty(t *testing.T) {
	got := compute.GetLegendData(nil, "keys", "column", "stacked", "vertical", false, idLabel)
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}
