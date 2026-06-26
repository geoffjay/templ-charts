// Package tooltip mirrors @nivo/tooltip: the BasicTooltip / Chip /
// TableTooltip HTML templ components (float in a container div, swapped in by
// HTMX), the Crosshair SVG component, and the TooltipPosition /
// TooltipAnchor / CrosshairType enums.
package tooltip

// TooltipPosition is where a tooltip floats relative to the cursor.
type TooltipPosition string

const (
	TooltipPositionCursor TooltipPosition = "cursor"
	TooltipPositionFixed  TooltipPosition = "fixed"
)

// TooltipAnchor is the side of the cursor a fixed tooltip attaches to.
type TooltipAnchor string

const (
	TooltipAnchorTop    TooltipAnchor = "top"
	TooltipAnchorRight  TooltipAnchor = "right"
	TooltipAnchorBottom TooltipAnchor = "bottom"
	TooltipAnchorLeft   TooltipAnchor = "left"
	TooltipAnchorCenter TooltipAnchor = "center"
)

// CrosshairType enumerates nivo's 12 crosshair variants (x/y combinations +
// "bottom-" / "top-" leading-edge variants).
type CrosshairType string

const (
	CrosshairTypeTopLeft     CrosshairType = "top-left"
	CrosshairTypeTopRight    CrosshairType = "top-right"
	CrosshairTypeBottomLeft  CrosshairType = "bottom-left"
	CrosshairTypeBottomRight CrosshairType = "bottom-right"
	CrosshairTypeMiddleLeft  CrosshairType = "middle-left"
	CrosshairTypeMiddleRight CrosshairType = "middle-right"
	CrosshairTypeTop         CrosshairType = "top"
	CrosshairTypeBottom      CrosshairType = "bottom"
	CrosshairTypeMiddle      CrosshairType = "middle"
	CrosshairTypeLeft        CrosshairType = "left"
	CrosshairTypeRight       CrosshairType = "right"
	CrosshairTypeCross       CrosshairType = "cross"
)

// CrosshairProps is the input to the Crosshair SVG component.
type CrosshairProps struct {
	Type   CrosshairType
	Width  float64
	Height float64
	X      float64 // cursor x in chart units
	Y      float64 // cursor y in chart units
	Theme  CrosshairTheme
}

// CrosshairTheme mirrors theming.CrosshairLine (kept local to avoid an import
// cycle in the SVG rendering path; charts/tooltip is leaf-ish).
type CrosshairTheme struct {
	Stroke          string
	StrokeWidth     float64
	StrokeOpacity   float64
	StrokeDasharray string
}

// CrosshairLineGeometry is one line of the crosshair, in chart-local coords.
type CrosshairLineGeometry struct {
	X1, Y1, X2, Y2 float64
}

// ComputeCrosshairLines returns the SVG lines for a crosshair of the given
// type at (x, y) within a width × height chart area. Mirrors @nivo/tooltip's
// crosshair line geometry.
func ComputeCrosshairLines(props CrosshairProps) []CrosshairLineGeometry {
	w, h := props.Width, props.Height
	x, y := props.X, props.Y
	switch props.Type {
	case CrosshairTypeTopLeft:
		return []CrosshairLineGeometry{{0, y, x, y}, {x, 0, x, y}}
	case CrosshairTypeTopRight:
		return []CrosshairLineGeometry{{w, y, x, y}, {x, 0, x, y}}
	case CrosshairTypeBottomLeft:
		return []CrosshairLineGeometry{{0, y, x, y}, {x, h, x, y}}
	case CrosshairTypeBottomRight:
		return []CrosshairLineGeometry{{w, y, x, y}, {x, h, x, y}}
	case CrosshairTypeMiddleLeft:
		return []CrosshairLineGeometry{{0, y, x, y}, {x, 0, x, h}}
	case CrosshairTypeMiddleRight:
		return []CrosshairLineGeometry{{w, y, x, y}, {x, 0, x, h}}
	case CrosshairTypeTop:
		return []CrosshairLineGeometry{{x, 0, x, y}}
	case CrosshairTypeBottom:
		return []CrosshairLineGeometry{{x, y, x, h}}
	case CrosshairTypeLeft:
		return []CrosshairLineGeometry{{0, y, x, y}}
	case CrosshairTypeRight:
		return []CrosshairLineGeometry{{w, y, x, y}}
	case CrosshairTypeMiddle:
		return []CrosshairLineGeometry{{x, 0, x, h}}
	case CrosshairTypeCross:
		return []CrosshairLineGeometry{{0, y, w, y}, {x, 0, x, h}}
	}
	return nil
}
