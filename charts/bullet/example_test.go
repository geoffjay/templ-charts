package bullet_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/bullet"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleBullet renders a small set of KPI bullet rows to an SVG string with
// the charts/render helper.
func ExampleBullet() {
	svg, err := render.String(bullet.Bullet(bullet.BulletProps{
		Width: 400, Height: 200,
		Data: []bullet.BulletItemDatum{
			{ID: "temperature", Ranges: []float64{0, 20, 60, 100}, Measures: []float64{55}, Markers: []float64{72}},
			{ID: "power", Ranges: []float64{0, 30, 70, 120}, Measures: []float64{48, 88}, Markers: []float64{95}},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
