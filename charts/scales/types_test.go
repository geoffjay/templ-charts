package scales

import (
	"testing"
	"time"
)

func TestScaleSpecTypes(t *testing.T) {
	cases := []struct {
		spec ScaleSpec
		want ScaleType
	}{
		{ScaleLinearSpec{}, ScaleTypeLinear},
		{ScaleLogSpec{}, ScaleTypeLog},
		{ScaleSymlogSpec{}, ScaleTypeSymlog},
		{ScalePointSpec{}, ScaleTypePoint},
		{ScaleBandSpec{}, ScaleTypeBand},
		{ScaleTimeSpec{}, ScaleTypeTime},
	}
	for _, c := range cases {
		if got := c.spec.ScaleType(); got != c.want {
			t.Errorf("%T.ScaleType() = %v, want %v", c.spec, got, c.want)
		}
	}
}

func TestCreatePrecisionMethod(t *testing.T) {
	base := time.Date(2024, 3, 15, 10, 30, 45, 123456789, time.UTC)
	cases := []struct {
		precision TimePrecision
		want      time.Time
	}{
		{TimePrecisionMillisecond, base},
		{TimePrecisionSecond, time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)},
		{TimePrecisionMinute, time.Date(2024, 3, 15, 10, 30, 0, 0, time.UTC)},
		{TimePrecisionHour, time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)},
		{TimePrecisionDay, time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)},
		{TimePrecisionMonth, time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
		{TimePrecisionYear, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, c := range cases {
		got := CreatePrecisionMethod(c.precision)(base)
		if !got.Equal(c.want) {
			t.Errorf("precision %s: got %v, want %v", c.precision, got, c.want)
		}
	}
}

func TestCreateDateNormalizer_TimeInput(t *testing.T) {
	norm := CreateDateNormalizer("native", TimePrecisionDay, true)
	in := time.Date(2024, 3, 15, 10, 30, 45, 0, time.UTC)
	want := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if got := norm(in); !got.Equal(want) {
		t.Errorf("normalize(time) = %v, want %v", got, want)
	}
}

func TestCreateDateNormalizer_StringInput(t *testing.T) {
	// A non-native format triggers string parsing (RFC3339 and YYYY-MM-DD).
	norm := CreateDateNormalizer("%Y-%m-%d", TimePrecisionDay, true)
	want := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if got := norm("2024-03-15"); !got.Equal(want) {
		t.Errorf("normalize(date string) = %v, want %v", got, want)
	}
	if got := norm("2024-03-15T10:30:00Z"); !got.Equal(want) {
		t.Errorf("normalize(rfc3339 string) = %v, want day-truncated %v", got, want)
	}
	// Unparseable strings fall through to the zero time.
	if got := norm("bogus"); !got.IsZero() {
		t.Errorf("normalize(bogus) = %v, want zero time", got)
	}
	// Native format ignores strings entirely.
	native := CreateDateNormalizer("native", TimePrecisionDay, true)
	if got := native("2024-03-15"); !got.IsZero() {
		t.Errorf("native normalize(string) = %v, want zero time", got)
	}
	// Unsupported types yield the zero time.
	if got := native(42); !got.IsZero() {
		t.Errorf("normalize(int) = %v, want zero time", got)
	}
}

func TestParseTimeFormat(t *testing.T) {
	if got, err := parseTimeFormat("2024-03-15T10:00:00Z", "%Y", false); err != nil || got.IsZero() {
		t.Errorf("rfc3339 parse = %v, %v", got, err)
	}
	if got, err := parseTimeFormat("2024-03-15", "%Y-%m-%d", false); err != nil || got.IsZero() {
		t.Errorf("date parse = %v, %v", got, err)
	}
	if _, err := parseTimeFormat("bogus", "%Y", false); err == nil {
		t.Error("bogus string should error")
	} else if err.Error() != "time parse error" {
		t.Errorf("error = %q, want time parse error", err.Error())
	}
}

func TestTimePrecisionsOrder(t *testing.T) {
	if len(TimePrecisions) != 7 {
		t.Fatalf("TimePrecisions len = %d, want 7", len(TimePrecisions))
	}
	if TimePrecisions[0] != TimePrecisionMillisecond || TimePrecisions[6] != TimePrecisionYear {
		t.Errorf("TimePrecisions order unexpected: %v", TimePrecisions)
	}
}
