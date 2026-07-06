package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	pc "github.com/geoffjay/templ-charts/charts/parallelcoordinates"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "parallel-coordinates",
		Title:       "Parallel coordinates",
		Description: "One axis per variable; each record a polyline across linear/point scales.",
		Snippet: `import (
    pc "github.com/geoffjay/templ-charts/charts/parallelcoordinates"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(pc.ParallelCoordinates(pc.PCProps{
    Width: 720, Height: 440,
    Variables: []pc.PCVariable{
        {Key: "temp", Type: pc.PCScaleLinear, Label: "temperature"},
        {Key: "cost", Type: pc.PCScaleLinear, Label: "cost"},
        {Key: "grade", Type: pc.PCScalePoint, Label: "grade"},
    },
    Data: []pc.PCDatum{
        {ID: "batch A", Values: map[string]any{"temp": 20.0, "cost": 5.0, "grade": "B"}},
        {ID: "batch B", Values: map[string]any{"temp": 35.0, "cost": 9.0, "grade": "A"}},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := pc.PCProps{
				Width: 720, Height: 440, Responsive: true,
				Variables: []pc.PCVariable{
					{Key: "temp", Type: pc.PCScaleLinear, Label: "temperature"},
					{Key: "pressure", Type: pc.PCScaleLinear, Label: "pressure"},
					{Key: "yield", Type: pc.PCScaleLinear, Label: "yield %"},
					{Key: "cost", Type: pc.PCScaleLinear, Label: "cost"},
					{Key: "grade", Type: pc.PCScalePoint, Label: "grade"},
				},
				Data: []pc.PCDatum{
					{ID: "batch A", Values: map[string]any{"temp": 20.0, "pressure": 1.8, "yield": 74.0, "cost": 5.0, "grade": "B"}},
					{ID: "batch B", Values: map[string]any{"temp": 35.0, "pressure": 2.4, "yield": 91.0, "cost": 9.0, "grade": "A"}},
					{ID: "batch C", Values: map[string]any{"temp": 28.0, "pressure": 2.1, "yield": 83.0, "cost": 7.5, "grade": "A"}},
					{ID: "batch D", Values: map[string]any{"temp": 42.0, "pressure": 3.2, "yield": 68.0, "cost": 11.0, "grade": "C"}},
					{ID: "batch E", Values: map[string]any{"temp": 24.0, "pressure": 1.5, "yield": 79.0, "cost": 4.2, "grade": "B"}},
					{ID: "batch F", Values: map[string]any{"temp": 31.0, "pressure": 2.7, "yield": 88.0, "cost": 8.4, "grade": "A"}},
					{ID: "batch G", Values: map[string]any{"temp": 38.0, "pressure": 2.9, "yield": 72.0, "cost": 10.2, "grade": "C"}},
					{ID: "batch H", Values: map[string]any{"temp": 26.0, "pressure": 1.9, "yield": 86.0, "cost": 6.1, "grade": "B"}},
					{ID: "batch I", Values: map[string]any{"temp": 33.0, "pressure": 2.2, "yield": 94.0, "cost": 8.9, "grade": "A"}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(pc.ParallelCoordinates(p))
		},
	})
}
