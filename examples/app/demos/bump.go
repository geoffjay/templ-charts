package demos

import (
	"github.com/geoffjay/templ-charts/charts/bump"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
)

// BumpDemo is one bump tile on the /bump page (static render).
type BumpDemo struct {
	ID          string
	Title       string
	Description string
	Props       bump.BumpProps
}

// BumpDemos returns the bump demos for the /bump page.
func BumpDemos() []BumpDemo {
	data := []bump.BumpSerie{
		{ID: "React", Data: []bump.BumpDatum{
			{X: "2018", Y: 1}, {X: "2019", Y: 1}, {X: "2020", Y: 2}, {X: "2021", Y: 1}, {X: "2022", Y: 1},
		}},
		{ID: "Vue", Data: []bump.BumpDatum{
			{X: "2018", Y: 2}, {X: "2019", Y: 3}, {X: "2020", Y: 1}, {X: "2021", Y: 2}, {X: "2022", Y: 3},
		}},
		{ID: "Svelte", Data: []bump.BumpDatum{
			{X: "2018", Y: 3}, {X: "2019", Y: 2}, {X: "2020", Y: 3}, {X: "2021", Y: 4}, {X: "2022", Y: 2},
		}},
		{ID: "Angular", Data: []bump.BumpDatum{
			{X: "2018", Y: 4}, {X: "2019", Y: 4}, {X: "2020", Y: 4}, {X: "2021", Y: 3}, {X: "2022", Y: 4},
		}},
	}
	return []BumpDemo{
		{
			ID:          "bump-basic",
			Title:       "Smooth ranking over time",
			Description: "Four series ranked across five years, smooth (curveBumpX) interpolation with end labels and per-point hover.",
			Props: bump.BumpProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 30, Right: 100, Bottom: 40, Left: 100},
				Data:        data,
				Interactive: true,
			},
		},
		{
			ID:          "bump-linear-legend",
			Title:       "Linear interpolation + legend",
			Description: "interpolation=linear, start+end labels off, with a category legend.",
			Props: bump.BumpProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:        core.Margin{Top: 30, Right: 160, Bottom: 40, Left: 60},
				Data:          data,
				Interpolation: bump.InterpolationLinear,
				StartLabel:    bump.BoolPtr(false),
				EndLabel:      bump.BoolPtr(false),
				Legends: []legends.LegendProps{
					{Anchor: legends.LegendAnchorRight, Direction: legends.LegendDirectionColumn, TranslateX: 120},
				},
			},
		},
		{
			ID:          "bump-mesh",
			Title:       "Voronoi-mesh hover",
			Description: "useMesh: hover anywhere resolves to the nearest rank point via an accurate Voronoi mesh (internal/d3/delaunay), rather than needing to land on a point.",
			Props: bump.BumpProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 30, Right: 100, Bottom: 40, Left: 100},
				Data:        data,
				Interactive: true,
				UseMesh:     true,
			},
		},
	}
}
