package d3timeformat

import (
	"testing"
	"time"
)

// FuzzFormatString exercises the time-format directive scanner (tokenize) with
// arbitrary spec strings across a spread of time values. Malformed directives,
// trailing '%', and unknown modifiers must never panic. FormatString is
// reachable from public chart Props via caller-supplied time-format strings.
func FuzzFormatString(f *testing.F) {
	specs := []string{
		"", "%Y-%m-%d", "%H:%M:%S", "%", "%-d", "%_H", "%Q", "%%",
		"%j %U %W", "%a %b %e", "%p %I", "%%%Y",
	}
	// (unix seconds, nanoseconds) pairs: zero time, epoch, far future, pre-epoch.
	instants := [][2]int64{
		{0, 0}, {1_600_000_000, 0}, {1 << 40, 0}, {-2_000_000_000, 500},
	}
	for _, s := range specs {
		for _, in := range instants {
			f.Add(s, in[0], in[1])
		}
	}

	f.Fuzz(func(t *testing.T, spec string, sec int64, nsec int64) {
		tm := time.Unix(sec, nsec).UTC()
		_ = FormatString(spec, tm)
	})
}
