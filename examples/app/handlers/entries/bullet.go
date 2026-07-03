package entries

import (
	"github.com/geoffjay/templ-charts/charts/bullet"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "bullet",
		Title:       "Bullet",
		Description: "KPI ranges + measure bars + markers on a shared value scale, per-row axis.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/bullet"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(bullet.Bullet(bullet.BulletProps{
    Width: 720, Height: 440,
    Data: []bullet.BulletItemDatum{
        {ID: "temperature", Ranges: []float64{0, 20, 60, 100}, Measures: []float64{55}, Markers: []float64{72}},
        {ID: "power", Ranges: []float64{0, 30, 70, 120}, Measures: []float64{48, 88}, Markers: []float64{95}},
    },
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			p := bullet.BulletProps{
				Width: 720, Height: 440, Responsive: true,
				Data: []bullet.BulletItemDatum{
					{ID: "temperature", Ranges: []float64{0, 20, 60, 100}, Measures: []float64{55}, Markers: []float64{72}},
					{ID: "power", Ranges: []float64{0, 30, 70, 120}, Measures: []float64{48, 88}, Markers: []float64{95}},
				},
				Theme: theme,
			}
			p.Animate = animate
			return render.String(bullet.Bullet(p))
		},
	})
}
