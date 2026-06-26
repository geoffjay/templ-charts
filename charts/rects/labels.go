package rects

// LabelAnchor is one of the 9 label anchor positions for a rect (bar).
type LabelAnchor string

const (
	AnchorTopStart    LabelAnchor = "top-start"
	AnchorTop         LabelAnchor = "top"
	AnchorTopEnd      LabelAnchor = "top-end"
	AnchorMiddleStart LabelAnchor = "middle-start"
	AnchorMiddle      LabelAnchor = "middle"
	AnchorMiddleEnd   LabelAnchor = "middle-end"
	AnchorBottomStart LabelAnchor = "bottom-start"
	AnchorBottom      LabelAnchor = "bottom"
	AnchorBottomEnd   LabelAnchor = "bottom-end"
)

// LabelPosition is the x/y/textAnchor/dominantBaseline computed for a rect
// label by RectLabelPosition.
type LabelPosition struct {
	X, Y             float64
	TextAnchor       string
	DominantBaseline string
}

// RectLabelPosition computes the label position for a rect given an anchor and
// a label offset. Mirrors @nivo/bar's getLabelGenerator × labelPosition
// matrix. For vertical bars the rect grows upward from y+height; for
// horizontal bars it grows rightward from x. The anchor is interpreted in
// the rect's own orientation:
//
//   - "top*" anchors sit at the outer edge (y for vertical, x+w for horizontal)
//   - "bottom*" anchors sit at the inner edge (y+h for vertical, x for horizontal)
//   - "middle*" anchors center on the cross axis
//   - "*-start" left-aligns (text-anchor=start), "*-end" right-aligns
func RectLabelPosition(r Rect, anchor LabelAnchor, offset float64, horizontal bool) LabelPosition {
	if horizontal {
		switch anchor {
		case AnchorTopStart, AnchorTop, AnchorTopEnd:
			// "top" for horizontal bar = right end (x+w)
			x := r.X + r.Width + offset
			y := r.Y + r.Height/2
			return LabelPosition{X: x, Y: y, TextAnchor: anchorHoriz(anchor), DominantBaseline: "central"}
		case AnchorMiddleStart, AnchorMiddle, AnchorMiddleEnd:
			x := r.X + r.Width/2
			y := r.Y + r.Height/2
			return LabelPosition{X: x, Y: y, TextAnchor: anchorHoriz(anchor), DominantBaseline: "central"}
		case AnchorBottomStart, AnchorBottom, AnchorBottomEnd:
			x := r.X - offset
			y := r.Y + r.Height/2
			return LabelPosition{X: x, Y: y, TextAnchor: anchorHoriz(anchor), DominantBaseline: "central"}
		}
	}
	// vertical
	switch anchor {
	case AnchorTopStart, AnchorTop, AnchorTopEnd:
		y := r.Y - offset
		x := r.X + r.Width/2
		return LabelPosition{X: x, Y: y, TextAnchor: anchorHoriz(anchor), DominantBaseline: "auto"}
	case AnchorMiddleStart, AnchorMiddle, AnchorMiddleEnd:
		x := r.X + r.Width/2
		y := r.Y + r.Height/2
		return LabelPosition{X: x, Y: y, TextAnchor: anchorHoriz(anchor), DominantBaseline: "central"}
	case AnchorBottomStart, AnchorBottom, AnchorBottomEnd:
		y := r.Y + r.Height + offset
		x := r.X + r.Width/2
		return LabelPosition{X: x, Y: y, TextAnchor: anchorHoriz(anchor), DominantBaseline: "hanging"}
	}
	return LabelPosition{X: r.X + r.Width/2, Y: r.Y + r.Height/2, TextAnchor: "middle", DominantBaseline: "central"}
}

// anchorHoriz maps a 9-anchor to its horizontal text-anchor component.
func anchorHoriz(a LabelAnchor) string {
	switch a {
	case AnchorTopStart, AnchorMiddleStart, AnchorBottomStart:
		return "start"
	case AnchorTopEnd, AnchorMiddleEnd, AnchorBottomEnd:
		return "end"
	default:
		return "middle"
	}
}
