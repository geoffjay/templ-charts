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
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := pc.PCProps{
				Width: 720, Height: 440, Responsive: true,
				Variables: []pc.PCVariable{
					{Key: "temp", Type: pc.PCScaleLinear, Label: "temperature"},
					{Key: "cost", Type: pc.PCScaleLinear, Label: "cost"},
					{Key: "grade", Type: pc.PCScalePoint, Label: "grade"},
				},
				Data: []pc.PCDatum{
					{ID: "batch A", Values: map[string]any{"temp": 20.0, "cost": 5.0, "grade": "B"}},
					{ID: "batch B", Values: map[string]any{"temp": 35.0, "cost": 9.0, "grade": "A"}},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			return render.String(pc.ParallelCoordinates(p))
		},
	})
}
