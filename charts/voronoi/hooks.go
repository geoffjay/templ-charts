package voronoi

import (
	"github.com/geoffjay/templ-charts/charts/scales"
	"github.com/geoffjay/templ-charts/charts/theming"
	"github.com/geoffjay/templ-charts/internal/d3/delaunay"
)

// UseVoronoi mirrors @nivo/voronoi's useVoronoiMesh: it builds linear x/y scales
// from the domains into the chart area (range [0,width]/[0,height], no y-flip,
// matching nivo), positions each datum in pixel space, then constructs the
// Delaunay triangulation and its Voronoi diagram clipped to the chart bounds.
// props.Width/Height are the inner (margin-subtracted) dimensions.
func UseVoronoi(props VoronoiProps) VoronoiResult {
	xScale := scales.NewLinearScaleWithRange(props.XDomain[0], props.XDomain[1], 0, props.Width)
	yScale := scales.NewLinearScaleWithRange(props.YDomain[0], props.YDomain[1], 0, props.Height)

	points := make([]ComputedPoint, len(props.Data))
	pts := make([][2]float64, len(props.Data))
	for i, d := range props.Data {
		x := xScale.Call(d.X)
		y := yScale.Call(d.Y)
		points[i] = ComputedPoint{ID: d.ID, X: x, Y: y, Data: d}
		pts[i] = [2]float64{x, y}
	}

	del := delaunay.NewDelaunayFrom(pts)
	vor := del.Voronoi([4]float64{0, 0, props.Width, props.Height})

	return VoronoiResult{
		Points:   points,
		Delaunay: del,
		Voronoi:  vor,
	}
}

func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}
