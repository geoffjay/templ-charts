package geo

import "math"

// Projection fitting, ported from d3-geo src/projection/fit.js. Each method
// auto-derives scale + translate so a GeoJSON object fills a given box, removing
// the manual scale/center fiddling geo demos need. obj is any GeoObject
// (Geometry/Feature/FeatureCollection).

// fit is the shared driver: reset to scale 150 / translate (0,0) with any
// rectangular postclip temporarily removed, measure the projected bounds of
// obj, hand them to fitBounds (which sets the final scale/translate), then
// restore the postclip.
func (p *Projection) fit(fitBounds func(b [2][2]float64), obj GeoObject) *Projection {
	savedPost := p.postclip
	p.postclip = identityStream
	p.Scale(150).Translate(0, 0)
	b := NewPath(p).Bounds(obj)
	fitBounds(b)
	p.postclip = savedPost
	return p
}

// FitExtent scales and translates the projection so obj fills the box
// [[x0,y0],[x1,y1]] (centered on the shorter axis). Mirrors projection.fitExtent.
func (p *Projection) FitExtent(x0, y0, x1, y1 float64, obj GeoObject) *Projection {
	return p.fit(func(b [2][2]float64) {
		w := x1 - x0
		h := y1 - y0
		k := math.Min(w/(b[1][0]-b[0][0]), h/(b[1][1]-b[0][1]))
		x := x0 + (w-k*(b[1][0]+b[0][0]))/2
		y := y0 + (h-k*(b[1][1]+b[0][1]))/2
		p.Scale(150*k).Translate(x, y)
	}, obj)
}

// FitSize is FitExtent with the box [[0,0],[w,h]]. Mirrors projection.fitSize.
func (p *Projection) FitSize(w, h float64, obj GeoObject) *Projection {
	return p.FitExtent(0, 0, w, h, obj)
}

// FitWidth fits obj to the given width, deriving height from the aspect ratio.
// Mirrors projection.fitWidth.
func (p *Projection) FitWidth(width float64, obj GeoObject) *Projection {
	return p.fit(func(b [2][2]float64) {
		k := width / (b[1][0] - b[0][0])
		x := (width - k*(b[1][0]+b[0][0])) / 2
		y := -k * b[0][1]
		p.Scale(150*k).Translate(x, y)
	}, obj)
}

// FitHeight fits obj to the given height, deriving width from the aspect ratio.
// Mirrors projection.fitHeight.
func (p *Projection) FitHeight(height float64, obj GeoObject) *Projection {
	return p.fit(func(b [2][2]float64) {
		k := height / (b[1][1] - b[0][1])
		x := -k * b[0][0]
		y := (height - k*(b[1][1]+b[0][1])) / 2
		p.Scale(150*k).Translate(x, y)
	}, obj)
}
