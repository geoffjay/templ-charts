package theming

// Engine identifies a rendering backend: svg, css, or canvas. templ-charts
// only emits SVG, but the bridge tables are kept for parity with @nivo/theming.
type Engine string

const (
	EngineSVG    Engine = "svg"
	EngineCSS    Engine = "css"
	EngineCanvas Engine = "canvas"
)

// TextAlign enumerates horizontal text alignment values.
type TextAlign string

const (
	TextAlignStart  TextAlign = "start"
	TextAlignCenter TextAlign = "center"
	TextAlignEnd    TextAlign = "end"
)

// TextBaseline enumerates vertical text baseline values.
type TextBaseline string

const (
	TextBaselineTop    TextBaseline = "top"
	TextBaselineCenter TextBaseline = "center"
	TextBaselineBottom TextBaseline = "bottom"
)

// EngineStyleAttributesMapping maps textAlign/textBaseline to engine-specific
// attribute values. Mirrors @nivo/theming bridge.ts.
type EngineStyleAttributesMapping struct {
	TextAlign    map[TextAlign]string
	TextBaseline map[TextBaseline]string
}

// SVGStyleAttributesMapping maps to SVG text-anchor / dominant-baseline values.
var SVGStyleAttributesMapping = EngineStyleAttributesMapping{
	TextAlign: map[TextAlign]string{
		TextAlignStart:  "start",
		TextAlignCenter: "middle",
		TextAlignEnd:    "end",
	},
	TextBaseline: map[TextBaseline]string{
		TextBaselineTop:    "text-before-edge",
		TextBaselineCenter: "middle",
		TextBaselineBottom: "text-after-edge",
	},
}

// CSSStyleAttributesMapping maps to CSS text-align / vertical-align values.
var CSSStyleAttributesMapping = EngineStyleAttributesMapping{
	TextAlign: map[TextAlign]string{
		TextAlignStart:  "left",
		TextAlignCenter: "center",
		TextAlignEnd:    "right",
	},
	TextBaseline: map[TextBaseline]string{
		TextBaselineTop:    "top",
		TextBaselineCenter: "middle",
		TextBaselineBottom: "bottom",
	},
}

// CanvasStyleAttributesMapping maps to canvas textAlign/textBaseline values.
var CanvasStyleAttributesMapping = EngineStyleAttributesMapping{
	TextAlign: map[TextAlign]string{
		TextAlignStart:  "left",
		TextAlignCenter: "center",
		TextAlignEnd:    "right",
	},
	TextBaseline: map[TextBaseline]string{
		TextBaselineTop:    "top",
		TextBaselineCenter: "middle",
		TextBaselineBottom: "bottom",
	},
}

// StyleAttributesMapping maps an Engine to its attribute mapping.
var StyleAttributesMapping = map[Engine]EngineStyleAttributesMapping{
	EngineSVG:    SVGStyleAttributesMapping,
	EngineCSS:    CSSStyleAttributesMapping,
	EngineCanvas: CanvasStyleAttributesMapping,
}

// ConvertStyleAttribute converts a textAlign/textBaseline value for the given
// engine. Mirrors @nivo/theming convertStyleAttribute.
func ConvertStyleAttribute(engine Engine, attr string, value string) string {
	m, ok := StyleAttributesMapping[engine]
	if !ok {
		return value
	}
	switch attr {
	case "textAlign":
		if v, ok := m.TextAlign[TextAlign(value)]; ok {
			return v
		}
	case "textBaseline":
		if v, ok := m.TextBaseline[TextBaseline(value)]; ok {
			return v
		}
	}
	return value
}

// SanitizeSvgTextStyle strips outline* fields from a TextStyle for SVG output
// (SVG text doesn't render outlines the way HTML does). Mirrors nivo's
// sanitizeSvgTextStyle.
func SanitizeSvgTextStyle(style TextStyle) map[string]any {
	out := map[string]any{}
	if style.FontFamily != "" {
		out["fontFamily"] = style.FontFamily
	}
	if style.FontSize != nil {
		out["fontSize"] = style.FontSize
	}
	if style.Fill != "" {
		out["fill"] = style.Fill
	}
	for k, v := range style.Extra {
		out[k] = v
	}
	return out
}

// SanitizeHtmlTextStyle strips outline* fields and rewrites fill as color for
// HTML tooltip output. Mirrors nivo's sanitizeHtmlTextStyle.
func SanitizeHtmlTextStyle(style TextStyle) map[string]any {
	out := map[string]any{}
	if style.FontFamily != "" {
		out["fontFamily"] = style.FontFamily
	}
	if style.FontSize != nil {
		out["fontSize"] = style.FontSize
	}
	if style.Fill != "" {
		out["color"] = style.Fill
	}
	for k, v := range style.Extra {
		out[k] = v
	}
	return out
}
