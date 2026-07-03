package geo_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/geoffjay/templ-charts/charts/geo"
	"github.com/geoffjay/templ-charts/charts/render"
)

// exampleFeatures builds a tiny inline GeoJSON dataset with two small polygon
// features, avoiding the bundled world-countries file.
func exampleFeatures() []geo.Feature {
	return []geo.Feature{
		{
			Type: "Feature",
			ID:   "AAA",
			Geometry: geo.Geometry{
				Type: geo.TypePolygon,
				Coordinates: [][][2]float64{{
					{0, 0}, {10, 0}, {10, 10}, {0, 10}, {0, 0},
				}},
			},
		},
		{
			Type: "Feature",
			ID:   "BBB",
			Geometry: geo.Geometry{
				Type: geo.TypePolygon,
				Coordinates: [][][2]float64{{
					{20, 20}, {30, 20}, {30, 30}, {20, 30}, {20, 20},
				}},
			},
		},
	}
}

// ExampleGeoMap renders a GeoMap from a tiny inline GeoJSON dataset to an SVG
// string with the charts/render helper.
func ExampleGeoMap() {
	svg, err := render.String(geo.GeoMap(geo.GeoMapProps{
		Features:  exampleFeatures(),
		FillColor: "#a7c6da",
		GeoBase: geo.GeoBase{
			Width: 400, Height: 300,
			ProjectionScale: 60,
			BorderWidth:     0.4,
			BorderColor:     "#152238",
		},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}

// ExampleChoropleth renders a Choropleth binding a tiny {id,value} dataset onto
// inline features to an SVG string with the charts/render helper.
func ExampleChoropleth() {
	svg, err := render.String(geo.Choropleth(geo.ChoroplethProps{
		Features: exampleFeatures(),
		Data: []geo.ChoroplethDatum{
			{ID: "AAA", Value: 10},
			{ID: "BBB", Value: 90},
		},
		GeoBase: geo.GeoBase{
			Width: 400, Height: 300,
			ProjectionScale: 60,
			BorderWidth:     0.4,
			BorderColor:     "#152238",
		},
	}))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.HasPrefix(svg, "<svg"))
	// Output: true
}
