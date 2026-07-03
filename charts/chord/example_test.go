package chord_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/chord"
	"github.com/geoffjay/templ-charts/charts/render"
)

// ExampleChord renders a small chord diagram from a square flow matrix to an
// SVG string with the charts/render helper.
func ExampleChord() {
	svg, err := render.String(chord.Chord(chord.ChordProps{
		Width: 520, Height: 520,
		Keys: []string{"Tokyo", "Osaka", "Kyoto"},
		Data: [][]float64{
			{0, 15834, 6987},
			{12345, 0, 5432},
			{5678, 4321, 0},
		},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
