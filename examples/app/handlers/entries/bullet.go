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
					{ID: "revenue", Title: "revenue ($k)", Ranges: []float64{0, 150, 320, 500}, Measures: []float64{284}, Markers: []float64{400}},
					{ID: "latency", Title: "latency (ms)", Ranges: []float64{0, 80, 160, 240}, Measures: []float64{132}, Markers: []float64{100}},
					{ID: "cpu", Title: "cpu load (%)", Ranges: []float64{0, 40, 75, 100}, Measures: []float64{62, 81}, Markers: []float64{90}},
					{ID: "throughput", Title: "throughput (rps)", Ranges: []float64{0, 250, 600, 900}, Measures: []float64{510}, Markers: []float64{750}},
					{ID: "uptime", Title: "uptime (days)", Ranges: []float64{0, 90, 180, 365}, Measures: []float64{221}, Markers: []float64{300, 180}},
				},
				Theme: theme,
			}
			p.Animate = animate
			return render.String(bullet.Bullet(p))
		},
	})
}
