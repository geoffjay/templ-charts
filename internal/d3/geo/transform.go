package geo

// transform is a (lon,lat)↔(x,y) mapping with an optional inverse. Raw
// projections and rotations are transforms; compose chains them. Mirrors the
// function-with-.invert convention in d3-geo.
type transform struct {
	forward func(lambda, phi float64) (float64, float64)
	invert  func(x, y float64) (float64, float64)
}

// compose returns a∘b: apply a, then feed its output into b. The inverse is
// b⁻¹∘a⁻¹ when both parts are invertible. Ported from d3-geo src/compose.js.
func compose(a, b transform) transform {
	t := transform{
		forward: func(x, y float64) (float64, float64) {
			x, y = a.forward(x, y)
			return b.forward(x, y)
		},
	}
	if a.invert != nil && b.invert != nil {
		t.invert = func(x, y float64) (float64, float64) {
			x, y = b.invert(x, y)
			return a.invert(x, y)
		}
	}
	return t
}
