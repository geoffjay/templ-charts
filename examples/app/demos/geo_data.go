package demos

import (
	_ "embed"
	"encoding/json"
	"hash/fnv"
	"sort"

	"github.com/geoffjay/templ-charts/charts/geo"
)

// worldCountriesJSON is a minified copy of the Natural Earth world-countries
// GeoJSON that @nivo/geo's own demos ship (public domain; coordinates rounded
// to 3 decimals to keep the asset ~190KB). 176 country features with ISO-3166
// alpha-3 ids ("USA", "FRA", …) and a "name" property.
//
//go:embed world_countries.json
var worldCountriesJSON []byte

// rawFeatureCollection / rawFeature decode the embedded GeoJSON. Coordinates
// are decoded lazily per geometry type so they land in the concrete shapes
// internal/d3/geo streams.
type rawFeatureCollection struct {
	Features []rawFeature `json:"features"`
}

type rawFeature struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Properties struct {
		Name string `json:"name"`
	} `json:"properties"`
	Geometry rawGeometry `json:"geometry"`
}

type rawGeometry struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
}

var (
	loadedWorld     []geo.Feature
	loadedWorldName map[string]string
)

// worldCountries parses (once) the embedded GeoJSON into geo.Features. On a
// decode error it returns nil, so the geo page degrades gracefully rather than
// panicking.
func worldCountries() ([]geo.Feature, map[string]string) {
	if loadedWorld != nil {
		return loadedWorld, loadedWorldName
	}
	var fc rawFeatureCollection
	if err := json.Unmarshal(worldCountriesJSON, &fc); err != nil {
		return nil, nil
	}
	features := make([]geo.Feature, 0, len(fc.Features))
	names := make(map[string]string, len(fc.Features))
	for _, rf := range fc.Features {
		g, ok := decodeGeometry(rf.Geometry)
		if !ok {
			continue
		}
		features = append(features, geo.Feature{Type: "Feature", ID: rf.ID, Geometry: g})
		names[rf.ID] = rf.Properties.Name
	}
	loadedWorld, loadedWorldName = features, names
	return loadedWorld, loadedWorldName
}

func decodeGeometry(rg rawGeometry) (geo.Geometry, bool) {
	switch rg.Type {
	case geo.TypePolygon:
		var c [][][2]float64
		if json.Unmarshal(rg.Coordinates, &c) != nil {
			return geo.Geometry{}, false
		}
		return geo.Geometry{Type: geo.TypePolygon, Coordinates: c}, true
	case geo.TypeMultiPolygon:
		var c [][][][2]float64
		if json.Unmarshal(rg.Coordinates, &c) != nil {
			return geo.Geometry{}, false
		}
		return geo.Geometry{Type: geo.TypeMultiPolygon, Coordinates: c}, true
	default:
		return geo.Geometry{}, false
	}
}

// syntheticChoroplethData assigns each country a deterministic pseudo-value in
// [0, 1_000_000] derived from its id (stable across runs), matching the style
// of nivo's choropleth demo (which uses random values). A handful of ids are
// left out so the demo also shows the unknownColor path.
func syntheticChoroplethData(features []geo.Feature) []geo.ChoroplethDatum {
	skip := map[string]bool{"ATA": true, "GRL": true} // Antarctica, Greenland → unknown
	data := make([]geo.ChoroplethDatum, 0, len(features))
	for _, f := range features {
		if skip[f.ID] {
			continue
		}
		h := fnv.New32a()
		_, _ = h.Write([]byte(f.ID))
		data = append(data, geo.ChoroplethDatum{
			ID:    f.ID,
			Value: float64(h.Sum32() % 1_000_001),
		})
	}
	sort.Slice(data, func(i, j int) bool { return data[i].ID < data[j].ID })
	return data
}
