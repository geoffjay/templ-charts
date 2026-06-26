package d3timeformat

import (
	"testing"
	"time"
)

func mustParse(t *testing.T, layout, value string) time.Time {
	tm, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("time.Parse(%q, %q): %v", layout, value, err)
	}
	return tm
}

// eq checks a format spec against an expected string, using a plain
// t.Errorf format that doesn't conflict with the spec's own %X tokens.
func eq(t *testing.T, spec, want string, tm time.Time) {
	t.Helper()
	got := FormatString(spec, tm)
	if got != want {
		t.Errorf("spec=%q got=%q want=%q", spec, got, want)
	}
}

func TestYearMonthDay(t *testing.T) {
	tm := mustParse(t, "2006-01-02 15:04:05", "2024-03-15 14:30:45")
	eq(t, "%Y", "2024", tm)
	eq(t, "%y", "24", tm)
	eq(t, "%Y-%m-%d", "2024-03-15", tm)
	eq(t, "%-m/%-d/%Y", "3/15/2024", tm)
	eq(t, "%b %d, %Y", "Mar 15, 2024", tm)
	eq(t, "%B %d, %Y", "March 15, 2024", tm)
	eq(t, "%a %b %e", "Fri Mar 15", tm)
	eq(t, "%A", "Friday", tm)
}

func TestTimeComponents(t *testing.T) {
	tm := mustParse(t, "2006-01-02 15:04:05", "2024-03-15 14:30:45")
	eq(t, "%H:%M:%S", "14:30:45", tm)
	eq(t, "%I:%M %p", "02:30 PM", tm)
	eq(t, "%-I:%M %p", "2:30 PM", tm)
	// 2024 is a leap year: Jan(31)+Feb(29)+15 = day 75.
	eq(t, "%j", "075", tm)
}

func TestWeekNumber(t *testing.T) {
	// 2024-01-07 is the first Sunday of 2024 -> %U = 01.
	tm := mustParse(t, "2006-01-02", "2024-01-07")
	eq(t, "%U", "01", tm)
	// 2024-01-01 IS a Monday, so it's the start of %W week 1.
	tm = mustParse(t, "2006-01-02", "2024-01-01")
	eq(t, "%W", "01", tm)
	// 2024-01-07 is a Sunday; %W counts Monday-baseline weeks so a Sunday
	// after the first Monday is still week 1.
	tm = mustParse(t, "2006-01-02", "2024-01-07")
	eq(t, "%W", "01", tm)
}

func TestLiteralPercent(t *testing.T) {
	tm := mustParse(t, "2006-01-02", "2024-01-01")
	eq(t, "100%% complete", "100% complete", tm)
}

func TestUnknownDirective(t *testing.T) {
	tm := mustParse(t, "2006-01-02", "2024-01-01")
	// unknown directive %Q falls through to the directive char 'Q'.
	eq(t, "%Q", "Q", tm)
}

func TestEmptySpec(t *testing.T) {
	if f := Format(""); f != nil {
		t.Errorf("Format(\"\") should return nil")
	}
}

func TestTrailingPercent(t *testing.T) {
	tm := mustParse(t, "2006-01-02", "2024-01-01")
	eq(t, "foo%", "foo%", tm)
}

func TestOffset(t *testing.T) {
	tm := time.Date(2024, 3, 15, 14, 30, 45, 0, time.UTC)
	eq(t, "%Z", "+0000", tm)
	loc := time.FixedZone("IST", 5*3600+30*60)
	tm2 := time.Date(2024, 3, 15, 14, 30, 45, 0, loc)
	eq(t, "%Z", "+0530", tm2)
}

func TestTwoDigitYear(t *testing.T) {
	tm2024 := mustParse(t, "2006", "2024")
	eq(t, "%y", "24", tm2024)
	tm2000 := mustParse(t, "2006", "2000")
	eq(t, "%y", "00", tm2000)
}
