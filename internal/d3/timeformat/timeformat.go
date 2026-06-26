// Package d3timeformat provides a minimal port of d3-time-format covering
// the most common strftime-style directives used on chart axes.
//
// Supported directives (matching d3-time-format semantics for the subset
// we need):
//
//	%Y  - four-digit year (e.g. "2024")
//	%y  - two-digit year (e.g. "24")
//	%m  - zero-padded month (01..12)
//	%-m - month without padding (1..12)
//	%b  - abbreviated month name (Jan..Dec)
//	%B  - full month name (January..December)
//	%d  - zero-padded day of month (01..31)
//	%-d - day of month without padding (1..31)
//	%e  - space-padded day of month (" 1".."31")
//	%a  - abbreviated weekday (Sun..Sat)
//	%A  - full weekday (Sunday..Saturday)
//	%H  - zero-padded hour, 24h (00..23)
//	%-H - hour 24h without padding
//	%I  - zero-padded hour, 12h (01..12)
//	%p  - AM/PM
//	%M  - zero-padded minute (00..59)
//	%S  - zero-padded second (00..60)
//	%j  - zero-padded day of year (001..366)
//	%U  - week number (Sunday-first, 00..53)
//	%W  - week number (Monday-first, 00..53)
//	%Z  - timezone offset (e.g. "+0000")
//	%%  - literal percent
//
// Any other character is emitted verbatim. Unknown directives fall back
// to emitting the directive character (e.g. "%Q" -> "Q").
//
// Parsing (Parse) is not implemented in v1 — chart axes only need
// formatting for tick labels. Time-scale input is expected as a time.Time.
package d3timeformat

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Format returns a function that formats a time.Time according to the
// given specifier. Returns nil if the specifier is empty.
func Format(spec string) func(t time.Time) string {
	if spec == "" {
		return nil
	}
	tokens := tokenize(spec)
	return func(t time.Time) string {
		var b strings.Builder
		for _, tok := range tokens {
			b.WriteString(tok.format(t))
		}
		return b.String()
	}
}

// FormatString applies the specifier to a single time.
func FormatString(spec string, t time.Time) string {
	f := Format(spec)
	if f == nil {
		return t.Format(time.RFC3339)
	}
	return f(t)
}

// token is a parsed piece of a format string: either a literal string or
// a directive.
type token struct {
	literal   string
	directive byte
	pad       byte // '-' = no pad, '0' = zero (default), '_' = space, 0 = default
}

// tokenize parses a format spec into a list of tokens.
func tokenize(spec string) []token {
	var toks []token
	i := 0
	for i < len(spec) {
		c := spec[i]
		if c != '%' {
			// accumulate literal run
			j := i
			for j < len(spec) && spec[j] != '%' {
				j++
			}
			toks = append(toks, token{literal: spec[i:j]})
			i = j
			continue
		}
		// directive: %X or %-X or %_X or %0X (modifier then X)
		i++ // skip '%'
		if i >= len(spec) {
			// trailing '%' — emit literal '%'
			toks = append(toks, token{literal: "%"})
			break
		}
		tok := token{pad: 0}
		switch spec[i] {
		case '-', '_', '0':
			switch spec[i] {
			case '-':
				tok.pad = '-'
			case '_':
				tok.pad = '_'
			case '0':
				tok.pad = '0'
			}
			i++
			if i >= len(spec) {
				// trailing modifier — emit literal percent + modifier
				toks = append(toks, token{literal: "%" + string(spec[i-1])})
				break
			}
		}
		tok.directive = spec[i]
		i++
		toks = append(toks, tok)
	}
	return toks
}

// format applies a single token to a time.
func (t token) format(tm time.Time) string {
	if t.literal != "" {
		return t.literal
	}
	switch t.directive {
	case '%':
		return "%"
	case 'Y':
		return pad4(tm.Year())
	case 'y':
		return padN(tm.Year()%100, 2, t.pad)
	case 'm':
		return padN(int(tm.Month()), 2, t.pad)
	case 'b':
		return tm.Month().String()[:3]
	case 'B':
		return tm.Month().String()
	case 'd':
		return padN(tm.Day(), 2, t.pad)
	case 'e':
		// space-padded day
		return padN(tm.Day(), 2, '_')
	case 'a':
		return tm.Weekday().String()[:3]
	case 'A':
		return tm.Weekday().String()
	case 'H':
		return padN(tm.Hour(), 2, t.pad)
	case 'I':
		h := tm.Hour() % 12
		if h == 0 {
			h = 12
		}
		return padN(h, 2, t.pad)
	case 'p':
		if tm.Hour() < 12 {
			return "AM"
		}
		return "PM"
	case 'M':
		return padN(tm.Minute(), 2, t.pad)
	case 'S':
		return padN(tm.Second(), 2, t.pad)
	case 'j':
		return padN(tm.YearDay(), 3, t.pad)
	case 'U':
		return padN(weekNumber(tm, time.Sunday), 2, t.pad)
	case 'W':
		return padN(weekNumber(tm, time.Monday), 2, t.pad)
	case 'Z':
		_, off := tm.Zone()
		return formatOffset(off)
	}
	// unknown directive — emit the directive char (matches d3's tolerant fallthrough)
	return string(t.directive)
}

// padN pads a non-negative integer to `width` with leading zeros, honoring
// the modifier: '-' = no pad, '_' = space pad, '0' or unset = zero pad.
func padN(v, width int, mod byte) string {
	s := strconv.Itoa(v)
	if mod == '-' || len(s) >= width {
		return s
	}
	pad := byte('0')
	if mod == '_' {
		pad = ' '
	}
	return strings.Repeat(string(pad), width-len(s)) + s
}

// pad4 always zero-pads a year to 4 digits (d3 %Y behavior).
func pad4(v int) string {
	if v < 0 {
		return "-" + padN(-v, 4, '0')
	}
	return padN(v, 4, '0')
}

// weekNumber returns the week number of tm using the given first-day-of-week
// (Sunday or Monday). Weeks before the first `firstDay` of the year are 0.
func weekNumber(tm time.Time, firstDay time.Weekday) int {
	// ISO-ish week numbering: week starts on firstDay; week 0 is the days
	// before the first firstDay-of-year.
	yday := tm.YearDay() - 1 // 0-indexed
	wday := (int(tm.Weekday()) - int(firstDay) + 7) % 7
	return (yday - wday + 7) / 7
}

// formatOffset renders a seconds offset as ±HHMM.
func formatOffset(offSeconds int) string {
	sign := "+"
	if offSeconds < 0 {
		sign = "-"
		offSeconds = -offSeconds
	}
	h := offSeconds / 3600
	m := (offSeconds % 3600) / 60
	return fmt.Sprintf("%s%02d%02d", sign, h, m)
}
