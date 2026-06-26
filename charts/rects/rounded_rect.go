package rects

import (
	"fmt"
	"strings"
)

// RoundedRectProps are the props for the RoundedRect templ component.
type RoundedRectProps struct {
	X, Y, Width, Height float64
	Radius              BorderRadius
	Fill                string
	Stroke              string
	StrokeWidth         float64
	Opacity             float64
	// Animate, when true, emits SMIL <animate> elements growing the rect from
	// zero width (vertical) or zero height (horizontal) to its full size.
	// Mirrors nivo's bar enter animation (600ms ease, fill="freeze").
	Animate    bool
	Horizontal bool // horizontal vs vertical growth direction
	// AnimateDelay/AnimateDuration override the defaults (0s / 0.6s).
	AnimateDelay    string
	AnimateDuration string
}

// roundedRectComputed holds the resolved geometry + path used by the templ
// component. Computed once so the templ render is pure string emission.
type roundedRectComputed struct {
	Path        string
	Fill        string
	Stroke      string
	StrokeWidth float64
	Opacity     float64
	Animate     bool
	Horizontal  bool
	Delay       string
	Duration    string
	// Initial (pre-animation) width/height for SMIL animate.
	InitialW, InitialH float64
	FinalW, FinalH     float64
	X, Y               float64
}

func computeRoundedRect(p RoundedRectProps) roundedRectComputed {
	corners := p.Radius.Resolved(p.Width, p.Height)
	path := BuildRoundedRectPath(p.X, p.Y, p.Width, p.Height,
		corners.TopLeft, corners.TopRight, corners.BottomRight, corners.BottomLeft)
	c := roundedRectComputed{
		Path:        path,
		Fill:        p.Fill,
		Stroke:      p.Stroke,
		StrokeWidth: p.StrokeWidth,
		Opacity:     p.Opacity,
		Animate:     p.Animate,
		Horizontal:  p.Horizontal,
		X:           p.X,
		Y:           p.Y,
		FinalW:      p.Width,
		FinalH:      p.Height,
		Delay:       p.AnimateDelay,
		Duration:    p.AnimateDuration,
	}
	if c.Duration == "" {
		c.Duration = "0.6s"
	}
	if c.Delay == "" {
		c.Delay = "0s"
	}
	// Pre-animation rect: zero on the growth axis.
	if p.Horizontal {
		c.InitialW = 0
		c.InitialH = p.Height
	} else {
		c.InitialW = p.Width
		c.InitialH = 0
	}
	return c
}

// smilAnimateWidth returns the <animate> markup for width growth, or empty.
func smilAnimateWidth(c roundedRectComputed) string {
	if !c.Animate {
		return ""
	}
	return fmt.Sprintf(
		`<animate attributeName="width" from="%s" to="%s" begin="%s" dur="%s" fill="freeze"/>`,
		fmtR(c.InitialW), fmtR(c.FinalW), c.Delay, c.Duration,
	)
}

// smilAnimateHeight returns the <animate> markup for height growth, or empty.
func smilAnimateHeight(c roundedRectComputed) string {
	if !c.Animate {
		return ""
	}
	return fmt.Sprintf(
		`<animate attributeName="height" from="%s" to="%s" begin="%s" dur="%s" fill="freeze"/>`,
		fmtR(c.InitialH), fmtR(c.FinalH), c.Delay, c.Duration,
	)
}

// Render returns the SVG string for a rounded rect. This is the non-templ
// rendering path used by chart packages that compose SVG strings directly.
// The templ component below wraps it for component-style composition.
func Render(p RoundedRectProps) string {
	c := computeRoundedRect(p)
	var b strings.Builder
	b.WriteString(`<path d="`)
	b.WriteString(c.Path)
	b.WriteString(`" fill="`)
	b.WriteString(c.Fill)
	b.WriteString(`"`)
	if c.Opacity > 0 && c.Opacity < 1 {
		b.WriteString(fmt.Sprintf(` opacity="%s"`, fmtR(c.Opacity)))
	}
	if c.StrokeWidth > 0 && c.Stroke != "" {
		b.WriteString(fmt.Sprintf(` stroke="%s" stroke-width="%s"`, c.Stroke, fmtR(c.StrokeWidth)))
	}
	b.WriteString(">")
	if c.Animate {
		if c.Horizontal {
			b.WriteString(smilAnimateWidth(c))
		} else {
			b.WriteString(smilAnimateHeight(c))
		}
	}
	b.WriteString("</path>")
	return b.String()
}
