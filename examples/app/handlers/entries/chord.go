package entries

import (
	"github.com/geoffjay/templ-charts/charts/chord"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "chord",
		Title:       "Chord",
		Description: "Radial flow diagram (d3-chord): entity arcs sized by total flow, ribbons per sub-flow.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/chord"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(chord.Chord(chord.ChordProps{
    Width: 520, Height: 520,
    Keys: []string{"Tokyo", "Osaka", "Kyoto"},
    Data: [][]float64{
        {0, 15834, 6987},
        {12345, 0, 5432},
        {5678, 4321, 0},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := chord.ChordProps{
				Width: 520, Height: 520, Responsive: true,
				Keys: []string{"Tokyo", "Osaka", "Kyoto", "Sapporo", "Nagoya", "Fukuoka", "Sendai", "Hiroshima"},
				Data: [][]float64{
					// Asymmetric passenger flows (thousands) between cities.
					{0, 15834, 6987, 4210, 9640, 3875, 5120, 2440},
					{12345, 0, 8432, 1560, 7210, 4980, 1830, 3660},
					{5678, 9321, 0, 980, 3540, 1720, 1140, 2080},
					{5230, 1210, 760, 0, 1480, 620, 2890, 410},
					{8470, 6650, 2980, 1120, 0, 2340, 1690, 2760},
					{4210, 5540, 1490, 540, 2010, 0, 880, 3980},
					{6120, 1440, 930, 3240, 1370, 720, 0, 590},
					{1980, 4130, 2450, 360, 3120, 4420, 640, 0},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(chord.Chord(p))
		},
	})
}
