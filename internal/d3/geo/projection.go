package geo

import "math"

// Projection wraps a raw spherical projection with the standard d3-geo
// machinery: three-axis rotation, scale, translate, center, adaptive
// resampling, and a preclip stage (antimeridian by default). It is a stream
// transform — Stream(sink) returns a Sink that projects incoming (lon,lat)
// degrees into plane coordinates. Ported from d3-geo src/projection/index.js
// (angle/reflect and clipExtent omitted — see NOTES.md).
type Projection struct {
	project transform // raw projection (radians → unit plane)

	k    float64 // scale
	x, y float64 // translate

	lambda, phi                       float64 // center, radians
	deltaLambda, deltaPhi, deltaGamma float64 // rotate, radians
	delta2                            float64 // precision²

	preclip func(Sink) Sink

	rotate                 transform
	projectTransform       transform
	projectRotateTransform transform
	projectResample        func(Sink) Sink
}

func newProjection(project transform) *Projection {
	p := &Projection{
		project: project,
		k:       150,
		x:       480,
		y:       250,
		delta2:  0.5,
		preclip: clipAntimeridian,
	}
	p.recenter()
	return p
}

func (p *Projection) recenter() {
	cx, cy := p.project.forward(p.lambda, p.phi)
	dx := p.x - p.k*cx
	dy := p.y + p.k*cy

	p.rotate = rotateRadians(p.deltaLambda, p.deltaPhi, p.deltaGamma)
	p.projectTransform = transform{
		forward: func(lambda, phi float64) (float64, float64) {
			px, py := p.project.forward(lambda, phi)
			return px*p.k + dx, dy - py*p.k
		},
	}
	if p.project.invert != nil {
		p.projectTransform.invert = func(x, y float64) (float64, float64) {
			return p.project.invert((x-dx)/p.k, (dy-y)/p.k)
		}
	}
	p.projectRotateTransform = compose(p.rotate, p.projectTransform)
	p.projectResample = resample(p.projectTransform, p.delta2)
}

// Stream returns the projecting stream feeding sink. Chain (outermost first):
// radians → rotate → preclip → resample+project → sink.
func (p *Projection) Stream(sink Sink) Sink {
	return transformRadians(transformRotate(p.rotate)(p.preclip(p.projectResample(sink))))
}

// Project maps (lon,lat) degrees to plane coordinates directly (no clip or
// resample), matching calling a d3 projection as a function.
func (p *Projection) Project(lon, lat float64) (float64, float64) {
	return p.projectRotateTransform.forward(lon*radians, lat*radians)
}

// Invert maps plane coordinates back to (lon,lat) degrees, or NaNs if the raw
// projection has no inverse.
func (p *Projection) Invert(x, y float64) (lon, lat float64) {
	if p.projectRotateTransform.invert == nil {
		return math.NaN(), math.NaN()
	}
	lon, lat = p.projectRotateTransform.invert(x, y)
	return lon * degrees, lat * degrees
}

func (p *Projection) Scale(k float64) *Projection { p.k = k; p.recenter(); return p }

func (p *Projection) Translate(x, y float64) *Projection {
	p.x, p.y = x, y
	p.recenter()
	return p
}

func (p *Projection) Rotate(lambda, phi, gamma float64) *Projection {
	p.deltaLambda = lambda * radians
	p.deltaPhi = phi * radians
	p.deltaGamma = gamma * radians
	p.recenter()
	return p
}

func (p *Projection) Center(lon, lat float64) *Projection {
	p.lambda = lon * radians
	p.phi = lat * radians
	p.recenter()
	return p
}

func (p *Projection) Precision(delta float64) *Projection {
	p.delta2 = delta * delta
	p.recenter()
	return p
}

// ScaleValue / TranslateValue expose the current parameters (getters).
func (p *Projection) ScaleValue() float64                { return p.k }
func (p *Projection) TranslateValue() (float64, float64) { return p.x, p.y }
