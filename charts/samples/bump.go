package samples

import "github.com/geoffjay/templ-charts/charts/bump"

// Bump returns ready-to-render bump data: four framework series with their
// rank (Y) across five years (X).
func Bump() []bump.BumpSerie {
	return []bump.BumpSerie{
		{ID: "React", Data: []bump.BumpDatum{{X: "2018", Y: 1}, {X: "2019", Y: 1}, {X: "2020", Y: 2}, {X: "2021", Y: 1}, {X: "2022", Y: 1}}},
		{ID: "Vue", Data: []bump.BumpDatum{{X: "2018", Y: 2}, {X: "2019", Y: 3}, {X: "2020", Y: 1}, {X: "2021", Y: 2}, {X: "2022", Y: 3}}},
		{ID: "Svelte", Data: []bump.BumpDatum{{X: "2018", Y: 3}, {X: "2019", Y: 2}, {X: "2020", Y: 3}, {X: "2021", Y: 4}, {X: "2022", Y: 2}}},
		{ID: "Angular", Data: []bump.BumpDatum{{X: "2018", Y: 4}, {X: "2019", Y: 4}, {X: "2020", Y: 4}, {X: "2021", Y: 3}, {X: "2022", Y: 4}}},
	}
}
