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
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := chord.ChordProps{
				Width: 520, Height: 520, Responsive: true,
				Keys: []string{"Tokyo", "Osaka", "Kyoto"},
				Data: [][]float64{
					{0, 15834, 6987},
					{12345, 0, 5432},
					{5678, 4321, 0},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(chord.Chord(p))
		},
	})
}
