package calendar_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/geoffjay/templ-charts/charts/calendar"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/internal/golden"
)

func date(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func sampleData() []calendar.CalendarDatum {
	return []calendar.CalendarDatum{
		{Day: "2024-01-05", Value: 10},
		{Day: "2024-01-20", Value: 50},
		{Day: "2024-02-14", Value: 90},
		{Day: "2024-03-01", Value: 30},
		{Day: "2024-03-28", Value: 70},
	}
}

func render(t *testing.T, props calendar.CalendarProps) string {
	t.Helper()
	var b strings.Builder
	if err := calendar.Calendar(props).Render(context.Background(), &b); err != nil {
		t.Fatalf("Calendar.Render: %v", err)
	}
	return b.String()
}

func baseProps() calendar.CalendarProps {
	return calendar.CalendarProps{
		Width: 800, Height: 220,
		Margin: core.Margin{Top: 40, Right: 40, Bottom: 10, Left: 40},
		From:   date("2024-01-01"), To: date("2024-03-31"),
		Data: sampleData(),
	}
}

func TestCalendar_RendersSVG(t *testing.T) {
	out := render(t, baseProps())
	if !strings.HasPrefix(out, "<svg") || !strings.Contains(out, "</svg>") {
		t.Fatalf("expected a full <svg>…</svg>")
	}
}

func TestCalendar_DayCellCount(t *testing.T) {
	// 2024-01-01 .. 2024-03-31 inclusive = 91 days → 91 day <rect> + 1
	// background rect = 92.
	out := render(t, baseProps())
	if got := strings.Count(out, "<rect"); got != 92 {
		t.Errorf("rect count = %d, want 92 (91 days + 1 background)", got)
	}
}

func TestCalendar_MonthLegends(t *testing.T) {
	out := render(t, baseProps())
	for _, m := range []string{">Jan<", ">Feb<", ">Mar<"} {
		if !strings.Contains(out, m) {
			t.Errorf("missing month legend %q", m)
		}
	}
}

func TestCalendar_DerivesRangeFromData(t *testing.T) {
	// Without explicit From/To, the range comes from the data days.
	p := baseProps()
	p.From, p.To = time.Time{}, time.Time{}
	out := render(t, p)
	// Data spans 2024-01-05 .. 2024-03-28 = 84 days → 84 + 1 background.
	if got := strings.Count(out, "<rect"); got != 85 {
		t.Errorf("derived-range rect count = %d, want 85", got)
	}
}

func TestCalendar_Golden(t *testing.T) {
	golden.Assert(t, "calendar-quarter", render(t, baseProps()))
}
