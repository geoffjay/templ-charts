package core

import (
	"fmt"
	"math"
)

// patternDotsParams holds the computed geometry for a patternDots def.
type patternDotsParams struct {
	FullSize, Radius, HalfPadding float64
	Color, Background             string
}

func computePatternDots(d Def) patternDotsParams {
	size := patternNum(d.Size, PatternDotsDefDefaultProps.Size)
	padding := patternNum(d.Padding, PatternDotsDefDefaultProps.Padding)
	fullSize := size + padding
	if d.Stagger {
		fullSize = size*2 + padding*2
	}
	return patternDotsParams{
		FullSize:    fullSize,
		Radius:      size / 2,
		HalfPadding: padding / 2,
		Color:       patternStr(d.Color, PatternDotsDefDefaultProps.Color),
		Background:  patternStr(d.Background, PatternDotsDefDefaultProps.Background),
	}
}

// patternSquaresParams holds the computed geometry for a patternSquares def.
type patternSquaresParams struct {
	FullSize, Size, HalfPadding float64
	Color, Background           string
}

func computePatternSquares(d Def) patternSquaresParams {
	size := patternNum(d.Size, PatternSquaresDefDefaultProps.Size)
	padding := patternNum(d.Padding, PatternSquaresDefDefaultProps.Padding)
	fullSize := size + padding
	if d.Stagger {
		fullSize = size*2 + padding*2
	}
	return patternSquaresParams{
		FullSize:    fullSize,
		Size:        size,
		HalfPadding: padding / 2,
		Color:       patternStr(d.Color, PatternSquaresDefDefaultProps.Color),
		Background:  patternStr(d.Background, PatternSquaresDefDefaultProps.Background),
	}
}

// patternLinesParams holds the computed geometry for a patternLines def.
type patternLinesParams struct {
	Width, Height, LineWidth float64
	Path                     string
	Color, Background        string
}

func computePatternLines(d Def) patternLinesParams {
	spacing := patternNum(d.Spacing, PatternLinesDefDefaultProps.Spacing)
	rotationIn := patternNum(d.Rotation, PatternLinesDefDefaultProps.Rotation)
	background := patternStr(d.Background, PatternLinesDefDefaultProps.Background)
	color := patternStr(d.Color, PatternLinesDefDefaultProps.Color)
	lineWidth := patternNum(d.LineWidth, PatternLinesDefDefaultProps.LineWidth)

	rotation := math.Mod(rotationIn, 360)
	spacing = math.Abs(spacing)
	if rotation > 180 {
		rotation -= 360
	} else if rotation > 90 {
		rotation -= 180
	} else if rotation < -180 {
		rotation += 360
	} else if rotation < -90 {
		rotation += 180
	}
	width := spacing
	height := spacing
	path := ""
	if rotation == 0 {
		path = fmt.Sprintf("M 0 0 L %s 0 M 0 %s L %s %s",
			fmtFloat(width), fmtFloat(height), fmtFloat(width), fmtFloat(height))
	} else if rotation == 90 {
		path = fmt.Sprintf("M 0 0 L 0 %s M %s 0 L %s %s",
			fmtFloat(height), fmtFloat(width), fmtFloat(width), fmtFloat(height))
	} else {
		width = math.Abs(spacing / math.Sin(degreesToRadians(rotation)))
		height = spacing / math.Sin(degreesToRadians(90-rotation))
		if rotation > 0 {
			path = fmt.Sprintf("M 0 %s L %s %s M %s %s L %s %s M %s 0 L %s %s",
				fmtFloat(-height), fmtFloat(width*2), fmtFloat(height),
				fmtFloat(-width), fmtFloat(-height), fmtFloat(width), fmtFloat(height),
				fmtFloat(-width), fmtFloat(width), fmtFloat(height*2))
		} else {
			path = fmt.Sprintf("M %s %s L %s %s M %s %s L %s %s M 0 %s L %s 0",
				fmtFloat(-width), fmtFloat(height), fmtFloat(width), fmtFloat(-height),
				fmtFloat(-width), fmtFloat(height*2), fmtFloat(width*2), fmtFloat(-height),
				fmtFloat(height*2), fmtFloat(width*2))
		}
	}
	return patternLinesParams{
		Width: width, Height: height, LineWidth: lineWidth, Path: path,
		Color: color, Background: background,
	}
}

// markerLayout holds computed legend label coordinates for a cartesian marker.
type markerLayout struct {
	X, Y, Rotation float64
	TextAnchor     string
}

func computeMarkerLayout(m CartesianMarker, width, height float64) markerLayout {
	position := m.LegendPosition
	if position == "" {
		position = "top-right"
	}
	orientation := m.LegendOrientation
	if orientation == "" {
		orientation = "horizontal"
	}
	offsetX := m.LegendOffsetX
	if offsetX == 0 {
		offsetX = 14
	}
	offsetY := m.LegendOffsetY
	if offsetY == 0 {
		offsetY = 14
	}
	x, y, rotation, textAnchor := computeMarkerLabel(m.Axis, width, height, position, offsetX, offsetY, orientation)
	return markerLayout{X: x, Y: y, Rotation: rotation, TextAnchor: textAnchor}
}

// markerLine holds the computed line coordinates for a cartesian marker.
type markerLine struct {
	X, Y, X2, Y2 float64
}

func computeMarkerLine(m CartesianMarker, width, height float64, xScale, yScale func(any) float64) markerLine {
	if m.Axis == "y" {
		return markerLine{X: 0, Y: yScale(m.Value), X2: width, Y2: 0}
	}
	return markerLine{X: xScale(m.Value), Y: 0, X2: 0, Y2: height}
}
