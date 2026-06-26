package core

import (
	"reflect"
	"strings"
)

// Def is the union type for SVG defs configurations. The Type field
// discriminates between gradient and pattern defs. Mirrors @nivo/core's
// gradientTypes/patternTypes unions.
type Def struct {
	// ID is the def identifier referenced by url(#ID) in fill/stroke attrs.
	ID string

	// Type is one of: "linearGradient", "patternDots", "patternLines",
	// "patternSquares".
	Type string

	// LinearGradient fields.
	Colors []GradientStop

	// Pattern fields shared by dots/lines/squares.
	Background string
	Color      string
	Size       float64
	Padding    float64
	Stagger    bool
	Spacing    float64
	Rotation   float64
	LineWidth  float64

	// Spreading extra options (x1/x2/y1/y2 for gradients, etc.) via Raw.
	Raw map[string]any
}

// GradientStop is a single color stop in a linearGradient.
type GradientStop struct {
	Offset  float64
	Color   string
	Opacity float64
}

// DefType constants.
const (
	DefTypeLinearGradient string = "linearGradient"
	DefTypePatternDots    string = "patternDots"
	DefTypePatternLines   string = "patternLines"
	DefTypePatternSquares string = "patternSquares"
)

// LinearGradientDef constructs a linearGradient def config. Mirrors nivo's
// linearGradientDef factory.
func LinearGradientDef(id string, colors []GradientStop, options map[string]any) Def {
	return Def{ID: id, Type: DefTypeLinearGradient, Colors: colors, Raw: options}
}

// PatternDotsDefDefaultProps mirrors nivo's PatternDotsDefaultProps.
var PatternDotsDefDefaultProps = struct {
	Color, Background string
	Size, Padding     float64
	Stagger           bool
}{
	Color: "#000000", Background: "#ffffff", Size: 4, Padding: 4, Stagger: false,
}

// PatternDotsDef constructs a patternDots def config. Mirrors patternDotsDef.
func PatternDotsDef(id string, options map[string]any) Def {
	return Def{ID: id, Type: DefTypePatternDots, Raw: options}
}

// PatternLinesDefDefaultProps mirrors nivo's PatternLinesDefaultProps.
var PatternLinesDefDefaultProps = struct {
	Spacing, Rotation float64
	Background, Color string
	LineWidth         float64
}{
	Spacing: 5, Rotation: 0, Background: "#000000", Color: "#ffffff", LineWidth: 2,
}

// PatternLinesDef constructs a patternLines def config. Mirrors patternLinesDef.
func PatternLinesDef(id string, options map[string]any) Def {
	return Def{ID: id, Type: DefTypePatternLines, Raw: options}
}

// PatternSquaresDefDefaultProps mirrors nivo's PatternSquaresDefaultProps.
var PatternSquaresDefDefaultProps = struct {
	Color, Background string
	Size, Padding     float64
	Stagger           bool
}{
	Color: "#000000", Background: "#ffffff", Size: 4, Padding: 4, Stagger: false,
}

// PatternSquaresDef constructs a patternSquares def config.
func PatternSquaresDef(id string, options map[string]any) Def {
	return Def{ID: id, Type: DefTypePatternSquares, Raw: options}
}

// SvgDefsAndFill is the bound output of BindDefs: the list of defs to render
// and the per-node fill overrides keyed by node index. Mirrors nivo's
// bindDefs return value (mutated nodes + bound defs list).
type SvgDefsAndFill struct {
	Defs []Def
	// FillByNodeIndex maps a node index to a fill override ("url(#id)").
	// Charts apply this to each node's fill attr when rendering.
	FillByNodeIndex map[int]string
}

// DefRule is a rule that conditionally applies a def to a node. Match is
// either "*" (match all), a func(any) bool, or a map[string]any compared by
// equality against a subset of the node's (or node.data's) properties.
type DefRule struct {
	ID    string
	Match any // "*" | func(any) bool | map[string]any
}

// BindDefs mirrors @nivo/core lib/defs.js bindDefs: given base defs, nodes,
// and rules, returns the expanded defs list (with inheritance-generated
// variants) and per-node fill overrides.
//
// dataKey, colorKey (default "color"), and targetKey (default "fill") mirror
// nivo's options. Node values are reached via getPath (string path) or
// directly when dataKey is empty.
func BindDefs(
	defs []Def,
	nodes []any,
	rules []DefRule,
	dataKey string,
	colorKey string,
	targetKey string,
) SvgDefsAndFill {
	if colorKey == "" {
		colorKey = "color"
	}
	if targetKey == "" {
		targetKey = "fill"
	}
	out := SvgDefsAndFill{FillByNodeIndex: map[int]string{}}
	if len(defs) == 0 || len(nodes) == 0 {
		out.Defs = append([]Def(nil), defs...)
		return out
	}

	boundDefs := append([]Def(nil), defs...)
	generated := map[string]bool{}

	for i, node := range nodes {
		for _, rule := range rules {
			if !isMatchingDef(rule.Match, node, dataKey) {
				continue
			}
			def := findDef(defs, rule.ID)
			if def == nil {
				break
			}
			switch def.Type {
			case DefTypePatternDots, DefTypePatternLines, DefTypePatternSquares:
				bg := def.Background
				fg := def.Color
				inheritedID := def.ID
				needVariant := false
				if bg == "inherit" || fg == "inherit" {
					nc := nodeColor(node, colorKey)
					if bg == "inherit" {
						inheritedID = inheritedID + ".bg." + nc
						bg = nc
						needVariant = true
					}
					if fg == "inherit" {
						inheritedID = inheritedID + ".fg." + nc
						fg = nc
						needVariant = true
					}
				}
				if needVariant {
					out.FillByNodeIndex[i] = "url(#" + inheritedID + ")"
					if !generated[inheritedID] {
						v := *def
						v.ID = inheritedID
						v.Background = bg
						v.Color = fg
						boundDefs = append(boundDefs, v)
						generated[inheritedID] = true
					}
				} else {
					out.FillByNodeIndex[i] = "url(#" + def.ID + ")"
				}
			case DefTypeLinearGradient:
				hasInherit := false
				for _, stop := range def.Colors {
					if stop.Color == "inherit" {
						hasInherit = true
						break
					}
				}
				if !hasInherit {
					out.FillByNodeIndex[i] = "url(#" + def.ID + ")"
					break
				}
				nc := nodeColor(node, colorKey)
				inheritedID := def.ID
				colors := make([]GradientStop, len(def.Colors))
				for j, stop := range def.Colors {
					if stop.Color == "inherit" {
						inheritedID = inheritedID + "." + itoa(j) + "." + nc
						colors[j] = GradientStop{Offset: stop.Offset, Color: nc, Opacity: stop.Opacity}
					} else {
						colors[j] = stop
					}
				}
				out.FillByNodeIndex[i] = "url(#" + inheritedID + ")"
				if !generated[inheritedID] {
					v := *def
					v.ID = inheritedID
					v.Colors = colors
					boundDefs = append(boundDefs, v)
					generated[inheritedID] = true
				}
			}
			break
		}
	}
	out.Defs = boundDefs
	return out
}

func findDef(defs []Def, id string) *Def {
	for i := range defs {
		if defs[i].ID == id {
			return &defs[i]
		}
	}
	return nil
}

func nodeColor(node any, colorKey string) string {
	v, _ := getPath(reflect.ValueOf(node), splitPath(colorKey))
	return toString(v)
}

// isMatchingDef mirrors nivo's isMatchingDef: "*" matches all; a func matches
// by predicate; a map matches by equality of a picked subset of node props.
func isMatchingDef(match any, node any, dataKey string) bool {
	switch m := match.(type) {
	case nil:
		return false
	case string:
		return m == "*"
	case func(any) bool:
		return m(node)
	case map[string]any:
		data := node
		if dataKey != "" {
			if v, ok := getPath(reflect.ValueOf(node), splitPath(dataKey)); ok {
				data = v
			}
		}
		// equality check: for each key in m, data[key] must equal m[key].
		for k, want := range m {
			got, _ := getPath(reflect.ValueOf(data), splitPath(k))
			if !equalAny(got, want) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func splitPath(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ".")
}

func equalAny(a, b any) bool { return toString(a) == toString(b) }

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	neg := i < 0
	if neg {
		i = -i
	}
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
