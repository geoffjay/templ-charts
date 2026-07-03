package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/funnel"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "funnel",
		Title:       "Funnel",
		Description: "Ordered parts as smooth/linear trapezoids with separators and labels.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/funnel"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(funnel.Funnel(funnel.FunnelProps{
    Width: 720, Height: 440,
    Data: []funnel.FunnelDatum{
        {ID: "sent", Label: "Sent", Value: 60000},
        {ID: "viewed", Label: "Viewed", Value: 38000},
        {ID: "clicked", Label: "Clicked", Value: 22000},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := funnel.FunnelProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []funnel.FunnelDatum{
					{ID: "sent", Label: "Sent", Value: 60000},
					{ID: "viewed", Label: "Viewed", Value: 38000},
					{ID: "clicked", Label: "Clicked", Value: 22000},
				},
				Theme: theme,
			}
			if palette != "" {
				p.Colors = colors.Scheme(palette)
			}
			p.Animate = animate
			return render.String(funnel.Funnel(p))
		},
	})
}
