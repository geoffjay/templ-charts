package core

import (
	"fmt"
	"strings"
	"time"

	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
	d3timeformat "github.com/geoffjay/templ-charts/internal/d3/timeformat"
)

// ValueFormat formats a value of type V into a string. nivo supports three
// forms: a user-supplied func, a string spec prefixed with "time:" (parsed
// as a d3-time-format spec), or a plain d3-format spec. A nil/empty format
// falls back to "%v" formatting.
type ValueFormat[V any] any // func(V) string | string

// GetValueFormatter mirrors nivo's getValueFormatter: a user func is returned
// as-is; a "time:"-prefixed string becomes a d3-time-format func; any other
// string is parsed as a d3-format spec; nil yields a default "%v" formatter.
//
// The returned func accepts any value and coerces it to the type expected by
// the chosen formatter (float64 for d3-format, time.Time for d3-time-format).
func GetValueFormatter[V any](format ValueFormat[V]) func(V) string {
	switch f := any(format).(type) {
	case nil:
		return func(v V) string { return fmt.Sprintf("%v", v) }
	case func(V) string:
		return f
	case string:
		s := strings.TrimSpace(f)
		if s == "" {
			return func(v V) string { return fmt.Sprintf("%v", v) }
		}
		if strings.HasPrefix(s, "time:") {
			spec := s[len("time:"):]
			tf := d3timeformat.Format(spec)
			return func(v V) string {
				if t, ok := toTime(v); ok {
					return tf(t)
				}
				return fmt.Sprintf("%v", v)
			}
		}
		nf := d3format.Format(s)
		return func(v V) string {
			return nf(toFloat64(v))
		}
	default:
		return func(v V) string { return fmt.Sprintf("%v", v) }
	}
}

// toTime extracts a time.Time from common date-ish values.
func toTime(v any) (time.Time, bool) {
	switch x := v.(type) {
	case time.Time:
		return x, true
	case *time.Time:
		if x == nil {
			return time.Time{}, false
		}
		return *x, true
	case string:
		if t, err := time.Parse(time.RFC3339, x); err == nil {
			return t, true
		}
		return time.Time{}, false
	case nil:
		return time.Time{}, false
	default:
		return time.Time{}, false
	}
}
