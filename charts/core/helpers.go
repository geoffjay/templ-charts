package core

import "math"

// opacityOr1 returns o if o > 0, else 1.
func opacityOr1(o float64) float64 {
	if o <= 0 {
		return 1
	}
	return o
}

// patternNum returns v if non-zero, else fallback. Used by pattern templ
// components to apply defaults from PatternDots/Lines/SquaresDefDefaultProps.
func patternNum(v, fallback float64) float64 {
	if v == 0 {
		return fallback
	}
	return v
}

// patternStr returns v if non-empty, else fallback.
func patternStr(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

// textAnchorOr returns v if non-empty, else fallback.
func textAnchorOr(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

// computeMarkerLabel mirrors nivo's CartesianMarkersItem.computeLabel: given
// axis, width, height, position, offsets, and orientation, returns the
// legend label x, y, rotation, and textAnchor.
func computeMarkerLabel(axis string, width, height float64, position string, offsetX, offsetY float64, orientation string) (x, y, rotation float64, textAnchor string) {
	rotation = 0
	textAnchor = "start"
	if orientation == "vertical" {
		rotation = -90
	}
	if axis == "x" {
		switch position {
		case "top-left":
			x = -offsetX
			y = offsetY
			textAnchor = "end"
		case "top":
			y = -offsetY
			if orientation == "horizontal" {
				textAnchor = "middle"
			} else {
				textAnchor = "start"
			}
		case "top-right":
			x = offsetX
			y = offsetY
			if orientation == "horizontal" {
				textAnchor = "start"
			} else {
				textAnchor = "end"
			}
		case "right":
			x = offsetX
			y = height / 2
			if orientation == "horizontal" {
				textAnchor = "start"
			} else {
				textAnchor = "middle"
			}
		case "bottom-right":
			x = offsetX
			y = height - offsetY
			textAnchor = "start"
		case "bottom":
			y = height + offsetY
			if orientation == "horizontal" {
				textAnchor = "middle"
			} else {
				textAnchor = "end"
			}
		case "bottom-left":
			y = height - offsetY
			x = -offsetX
			if orientation == "horizontal" {
				textAnchor = "end"
			} else {
				textAnchor = "start"
			}
		case "left":
			x = -offsetX
			y = height / 2
			if orientation == "horizontal" {
				textAnchor = "end"
			} else {
				textAnchor = "middle"
			}
		}
	} else {
		switch position {
		case "top-left":
			x = offsetX
			y = -offsetY
			textAnchor = "start"
		case "top":
			x = width / 2
			y = -offsetY
			if orientation == "horizontal" {
				textAnchor = "middle"
			} else {
				textAnchor = "start"
			}
		case "top-right":
			x = width - offsetX
			y = -offsetY
			if orientation == "horizontal" {
				textAnchor = "end"
			} else {
				textAnchor = "start"
			}
		case "right":
			x = width + offsetX
			if orientation == "horizontal" {
				textAnchor = "start"
			} else {
				textAnchor = "middle"
			}
		case "bottom-right":
			x = width - offsetX
			y = offsetY
			textAnchor = "end"
		case "bottom":
			x = width / 2
			y = offsetY
			if orientation == "horizontal" {
				textAnchor = "middle"
			} else {
				textAnchor = "end"
			}
		case "bottom-left":
			x = offsetX
			y = offsetY
			if orientation == "horizontal" {
				textAnchor = "start"
			} else {
				textAnchor = "end"
			}
		case "left":
			x = -offsetX
			if orientation == "horizontal" {
				textAnchor = "end"
			} else {
				textAnchor = "middle"
			}
		}
	}
	return x, y, rotation, textAnchor
}
