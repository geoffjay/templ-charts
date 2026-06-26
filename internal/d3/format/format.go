// Package d3format provides a pragmatic port of d3-format covering the
// specifier patterns most commonly used by chart authors and nivo defaults.
//
// d3-format's full grammar (fill, align, sign, symbol, zero, comma, width,
// precision, trim, type, plus locale) is large; this implementation supports
// the subset that matters for templ-charts v1 charts:
//
//	[[fill]align][sign][symbol][0][width][,.][precision][~][type]
//
// where:
//   - sign:   "-" (default, suppress + on positives), "+" (always), "(" (accounting)
//   - symbol: "$" (currency)
//   - zero:   "0" pads with leading zeros to width
//   - width:  integer field width
//   - comma:  "," inserts grouping separators (every 3 digits)
//   - precision: integer ".N"
//   - trim:   "~" trims trailing zeros
//   - type:   "b" binary, "c" char, "d" decimal int, "e"/"E" exponent,
//     "f"/"F" fixed-point, "g"/"G" general (sig figs), "o" octal,
//     "p" percent (sig figs), "r" rounded (sig figs, no prefix),
//     "s" SI-prefix, "x"/"X" hex, "%" percent × 100
//
// Unsupported features fall back to fmt-style behavior rather than failing.
// Behavior approximates d3-format; tests cover the common cases.
package d3format

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Format returns a formatter function for the given specifier.
// On parse failure returns nil; callers should fall back to a default.
func Format(spec string) func(v float64) string {
	f, ok := parseSpec(spec)
	if !ok {
		return nil
	}
	return func(v float64) string { return f.format(v) }
}

// FormatString applies the specifier to a single value, returning
// the formatted string or "<value>" on parse failure (matches nivo's
// fallback to `${value}`).
func FormatString(spec string, v float64) string {
	f := Format(spec)
	if f == nil {
		return fmt.Sprintf("%v", v)
	}
	return f(v)
}

type formatter struct {
	fill      byte
	align     byte // '<' '>' '^' '=' (0 = none)
	sign      byte // '-' '+' ' ' '('
	symbol    byte // '$' or 0
	zero      bool
	width     int
	comma     bool
	precision int  // -1 = not set
	trim      bool // '~' trims trailing zeros
	typ       byte // b,c,d,e,E,f,F,g,G,o,p,r,s,x,X,%, or 0 (=g-like)
}

const (
	defaultFill = ' '
	defaultType = 'g' // d3 default type when none specified
)

// parseSpec parses a d3-format specifier into a formatter struct.
// Grammar: [[fill]align][sign][symbol][0][width][,][.precision][~][type]
func parseSpec(spec string) (*formatter, bool) {
	if spec == "" {
		return &formatter{fill: defaultFill, sign: '-', precision: -1, typ: defaultType}, true
	}
	f := &formatter{fill: defaultFill, sign: '-', precision: -1, typ: defaultType}

	i := 0
	n := len(spec)

	// fill+align: a single char followed by '<' '>' '^' '='
	if n >= 2 {
		switch spec[1] {
		case '<', '>', '^', '=':
			f.fill = spec[0]
			f.align = spec[1]
			i = 2
		}
	}
	// align alone (no fill)
	if i == 0 && n >= 1 && (spec[0] == '<' || spec[0] == '>' || spec[0] == '^' || spec[0] == '=') {
		f.align = spec[0]
		i = 1
	}

	// sign
	if i < n && (spec[i] == '-' || spec[i] == '+' || spec[i] == ' ' || spec[i] == '(') {
		f.sign = spec[i]
		i++
	}
	// symbol
	if i < n && (spec[i] == '$') {
		f.symbol = spec[i]
		i++
	}
	// zero padding
	if i < n && spec[i] == '0' {
		f.zero = true
		i++
	}
	// width (digits)
	for i < n && spec[i] >= '0' && spec[i] <= '9' {
		f.width = f.width*10 + int(spec[i]-'0')
		i++
	}
	// comma / dot (d3 allows comma before or after dot, or comma alone)
	if i < n && spec[i] == ',' {
		f.comma = true
		i++
	}
	if i < n && spec[i] == '.' {
		i++
		j := i
		f.precision = 0
		for i < n && spec[i] >= '0' && spec[i] <= '9' {
			f.precision = f.precision*10 + int(spec[i]-'0')
			i++
		}
		if i == j {
			return nil, false
		}
	}
	if i < n && spec[i] == ',' {
		f.comma = true
		i++
	}
	// trim trailing zeros (~)
	if i < n && spec[i] == '~' {
		f.trim = true
		i++
	}
	// type
	if i < n {
		switch spec[i] {
		case 'b', 'c', 'd', 'e', 'E', 'f', 'F', 'g', 'G', 'o', 'p', 'r', 's', 'x', 'X', '%':
			f.typ = spec[i]
			i++
		default:
			return nil, false
		}
	}
	if i != n {
		return nil, false
	}
	return f, true
}

func (f *formatter) format(v float64) string {
	// Handle NaN, ±Inf
	if math.IsNaN(v) {
		return "NaN"
	}
	if math.IsInf(v, 1) {
		return "∞"
	}
	if math.IsInf(v, -1) {
		return "-∞"
	}

	neg := v < 0 || (v == 0 && math.Signbit(v))
	abs := math.Abs(v)

	var body string
	switch f.typ {
	case 'c':
		body = string(rune(int(abs)))
		if neg {
			body = "-" + body
		}
		return pad(f, body)

	case 'd':
		iv := int64(math.Round(abs))
		body = strconv.FormatInt(iv, 10)
		body = f.withSign(body, neg)
		body = f.withSymbol(body)
		body = f.applyGrouping(body, neg)
		body = f.applyZeroPad(body, neg)
		return pad(f, body)

	case 'b':
		iv := int64(math.Round(abs))
		body = strconv.FormatInt(iv, 2)
		body = f.withSign(body, neg)
		body = f.applyZeroPad(body, neg)
		return pad(f, body)

	case 'o':
		iv := int64(math.Round(abs))
		body = strconv.FormatInt(iv, 8)
		body = f.withSign(body, neg)
		body = f.applyZeroPad(body, neg)
		return pad(f, body)

	case 'x', 'X':
		iv := int64(math.Round(abs))
		body = strconv.FormatInt(iv, 16)
		if f.typ == 'X' {
			body = strings.ToUpper(body)
		}
		body = f.withSign(body, neg)
		body = f.applyZeroPad(body, neg)
		return pad(f, body)

	case 'p', '%':
		return f.formatPercent(v, neg)

	case 'e', 'E':
		return f.formatExponent(abs, neg)

	case 'g', 'G', 0:
		// default (no type) = g-like
		return f.formatGeneral(abs, neg)

	case 'f', 'F':
		return f.formatFixed(abs, neg)

	case 'r':
		// rounded: significant digits, no SI prefix
		return f.formatGeneral(abs, neg)

	case 's':
		// SI suffix
		return f.formatSI(abs, neg)
	}

	// shouldn't reach here
	return fmt.Sprintf("%v", v)
}

func (f *formatter) withSign(body string, neg bool) string {
	if neg {
		if f.sign == '(' {
			return "(" + body + ")"
		}
		return "-" + body
	}
	switch f.sign {
	case '+':
		return "+" + body
	case ' ':
		return " " + body
	case '(':
		return body
	}
	return body
}

func (f *formatter) withSymbol(body string) string {
	if f.symbol == '$' {
		return "$" + body
	}
	return body
}

// applyGrouping inserts ',' every 3 digits in the integer portion. The
// fractional part (if any) is preserved untouched.
func (f *formatter) applyGrouping(body string, neg bool) string {
	if !f.comma {
		return body
	}
	neg = strings.HasPrefix(body, "-")
	body = strings.TrimPrefix(body, "-")
	intPart, fracPart := body, ""
	if dot := strings.IndexByte(body, '.'); dot >= 0 {
		intPart, fracPart = body[:dot], body[dot+1:]
	}
	if len(intPart) > 3 {
		var b strings.Builder
		first := len(intPart) % 3
		if first > 0 {
			b.WriteString(intPart[:first])
			if len(intPart) > first {
				b.WriteByte(',')
			}
		}
		for i := first; i < len(intPart); i += 3 {
			b.WriteString(intPart[i : i+3])
			if i+3 < len(intPart) {
				b.WriteByte(',')
			}
		}
		intPart = b.String()
	}
	out := intPart
	if fracPart != "" || f.precision > 0 {
		out += "." + fracPart
	}
	if neg {
		return "-" + out
	}
	return out
}

// applyZeroPad pads the integer portion with leading zeros to f.width.
// Only applies when f.zero is true and width is set. Recognizes decimal
// and hex (a-f, A-F) digits as digit characters so it works for hex formats
// where the body may start with a non-decimal digit.
func (f *formatter) applyZeroPad(body string, neg bool) string {
	if !f.zero || f.width == 0 {
		return body
	}
	dot := strings.IndexByte(body, '.')
	intLen := len(body)
	if dot >= 0 {
		intLen = dot
	}
	isDigit := func(c byte) bool {
		return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
	}
	digitCount := 0
	for i := 0; i < intLen; i++ {
		if isDigit(body[i]) {
			digitCount++
		}
	}
	padN := f.width - digitCount
	if padN <= 0 {
		return body
	}
	pad := strings.Repeat("0", padN)
	insertAt := 0
	for insertAt < intLen {
		c := body[insertAt]
		if isDigit(c) {
			break
		}
		insertAt++
	}
	return body[:insertAt] + pad + body[insertAt:]
}

// pad right/center/lefts the body to f.width according to f.align.
func pad(f *formatter, body string) string {
	if f.width == 0 || len(body) >= f.width {
		return body
	}
	gap := f.width - len(body)
	fill := f.fill
	if fill == 0 {
		fill = ' '
	}
	switch f.align {
	case '<':
		return body + strings.Repeat(string(fill), gap)
	case '>':
		return strings.Repeat(string(fill), gap) + body
	case '^':
		left := gap / 2
		right := gap - left
		return strings.Repeat(string(fill), left) + body + strings.Repeat(string(fill), right)
	case '=':
		for i := 0; i < len(body); i++ {
			if body[i] >= '0' && body[i] <= '9' {
				return body[:i] + strings.Repeat(string(fill), gap) + body[i:]
			}
		}
		return strings.Repeat(string(fill), gap) + body
	}
	// default: right-align (d3 default)
	return strings.Repeat(string(fill), gap) + body
}

// formatFixed handles 'f'/'F' type.
func (f *formatter) formatFixed(abs float64, neg bool) string {
	prec := f.precision
	if prec < 0 {
		prec = 6
	}
	body := strconv.FormatFloat(abs, 'f', prec, 64)
	body = f.applyGrouping(body, neg)
	body = f.withSign(body, neg)
	body = f.withSymbol(body)
	body = f.applyZeroPad(body, neg)
	if f.trim {
		body = trimZeros(body)
	}
	return pad(f, body)
}

// formatExponent handles 'e'/'E' type.
func (f *formatter) formatExponent(abs float64, neg bool) string {
	prec := f.precision
	if prec < 0 {
		prec = 6
	}
	t := byte('e')
	if f.typ == 'E' {
		t = 'E'
	}
	body := strconv.FormatFloat(abs, t, prec, 64)
	body = f.withSign(body, neg)
	body = f.withSymbol(body)
	body = f.applyZeroPad(body, neg)
	if f.trim {
		body = trimZeros(body)
	}
	return pad(f, body)
}

// formatGeneral handles 'g'/'G'/none/'r'. Uses significant digits (default 6).
func (f *formatter) formatGeneral(abs float64, neg bool) string {
	sigfigs := f.precision
	if sigfigs < 0 {
		sigfigs = 6
	}
	t := byte('g')
	if f.typ == 'G' {
		t = 'G'
	}
	body := strconv.FormatFloat(abs, t, sigfigs, 64)
	body = f.applyGrouping(body, neg)
	body = f.withSign(body, neg)
	body = f.withSymbol(body)
	body = f.applyZeroPad(body, neg)
	if f.trim {
		body = trimZeros(body)
	}
	return pad(f, body)
}

// formatPercent handles 'p' and '%' types. Both multiply by 100; 'p' uses
// significant figures while '%' uses fixed precision (default 6).
func (f *formatter) formatPercent(v float64, neg bool) string {
	scaled := v * 100
	abs := math.Abs(scaled)
	var body string
	if f.typ == 'p' {
		sigfigs := f.precision
		if sigfigs < 0 {
			sigfigs = 6
		}
		body = strconv.FormatFloat(abs, 'g', sigfigs, 64)
	} else { // '%'
		prec := f.precision
		if prec < 0 {
			prec = 6
		}
		body = strconv.FormatFloat(abs, 'f', prec, 64)
	}
	body = f.applyGrouping(body, neg)
	body = f.withSign(body, neg)
	body = f.withSymbol(body)
	body = f.applyZeroPad(body, neg)
	if f.trim {
		body = trimZeros(body)
	}
	return pad(f, body+"%")
}

// formatSI handles 's' type (SI prefixes: k, M, G, T, P, m, µ, n, p, f).
func (f *formatter) formatSI(abs float64, neg bool) string {
	sigfigs := f.precision
	if sigfigs < 0 {
		sigfigs = 6
	}
	prefix, scaled := siPrefix(abs)
	body := strconv.FormatFloat(scaled, 'g', sigfigs, 64)
	if f.trim {
		body = trimZeros(body)
	}
	body = f.applyGrouping(body, neg)
	body = f.withSign(body, neg)
	body = f.withSymbol(body)
	if prefix != "" {
		body += prefix
	}
	return pad(f, body)
}

// siPrefix returns the SI prefix string and scaled value for the given
// absolute value. Returns "" and abs for values near 1 (no prefix).
func siPrefix(abs float64) (string, float64) {
	if abs == 0 || math.IsInf(abs, 0) || math.IsNaN(abs) {
		return "", abs
	}
	exp := int(math.Floor(math.Log10(abs)/3) * 3)
	prefixes := []struct {
		minExp int
		sym    string
	}{
		{-15, "f"},
		{-12, "p"},
		{-9, "n"},
		{-6, "µ"},
		{-3, "m"},
		{0, ""},
		{3, "k"},
		{6, "M"},
		{9, "G"},
		{12, "T"},
		{15, "P"},
	}
	for i := len(prefixes) - 1; i >= 0; i-- {
		if exp >= prefixes[i].minExp {
			return prefixes[i].sym, abs / math.Pow10(prefixes[i].minExp)
		}
	}
	return "f", abs / math.Pow10(-15)
}

// trimZeros trims trailing zeros from the fractional part of a numeric
// string (including the dot itself if all fractional digits are removed).
// Any non-digit suffix after the fractional run (e.g. "%", "k") is
// preserved.
func trimZeros(s string) string {
	dot := strings.IndexByte(s, '.')
	if dot < 0 {
		return s
	}
	// Find the end of the fractional digits (the start of any suffix).
	fracEnd := len(s)
	for i := len(s) - 1; i > dot; i-- {
		c := s[i]
		if c < '0' || c > '9' {
			fracEnd = i + 1
			break
		}
	}
	// Trim trailing zeros from [dot+1, fracEnd).
	newFracEnd := fracEnd
	for newFracEnd > dot+1 && s[newFracEnd-1] == '0' {
		newFracEnd--
	}
	if newFracEnd == dot+1 {
		newFracEnd-- // drop the dot too
	}
	return s[:newFracEnd] + s[fracEnd:]
}
