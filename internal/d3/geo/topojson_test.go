package geo

import (
	"math"
	"testing"
)

// TopoJSON decoding tests. The decoded coordinates match topojson-client@3's
// feature() output for the same topology.

const topoSample = `{
  "type":"Topology",
  "transform":{"scale":[0.01,0.01],"translate":[-1,-1]},
  "arcs":[[[0,0],[100,0],[0,100]],[[100,100],[-100,0],[0,-100]]],
  "objects":{"shapes":{"type":"GeometryCollection","geometries":[
    {"type":"Polygon","arcs":[[0,1]],"id":"P"},
    {"type":"LineString","arcs":[0],"id":"L"}
  ]}}
}`

func TestTopoJSON_Decode(t *testing.T) {
	fc, err := DecodeTopoJSON([]byte(topoSample), "shapes")
	if err != nil {
		t.Fatalf("DecodeTopoJSON: %v", err)
	}
	if len(fc.Features) != 2 {
		t.Fatalf("features = %d, want 2", len(fc.Features))
	}

	poly := fc.Features[0]
	if poly.ID != "P" || poly.Geometry.Type != TypePolygon {
		t.Fatalf("feature0 = %q/%q, want P/Polygon", poly.ID, poly.Geometry.Type)
	}
	wantPoly := [][][2]float64{{{-1, -1}, {0, -1}, {0, 0}, {-1, 0}, {-1, -1}}}
	rings, ok := poly.Geometry.Coordinates.([][][2]float64)
	if !ok {
		t.Fatalf("polygon coords type %T", poly.Geometry.Coordinates)
	}
	assertCoords3(t, "poly", rings, wantPoly)

	line := fc.Features[1]
	if line.ID != "L" || line.Geometry.Type != TypeLineString {
		t.Fatalf("feature1 = %q/%q, want L/LineString", line.ID, line.Geometry.Type)
	}
	wantLine := [][2]float64{{-1, -1}, {0, -1}, {0, 0}}
	pts, ok := line.Geometry.Coordinates.([][2]float64)
	if !ok {
		t.Fatalf("line coords type %T", line.Geometry.Coordinates)
	}
	assertCoords2(t, "line", pts, wantLine)
}

// TestTopoJSON_DecodedIsRenderable confirms the decoded features flow straight
// into the path generator (the whole point of the decoder).
func TestTopoJSON_DecodedIsRenderable(t *testing.T) {
	fc, err := DecodeTopoJSON([]byte(topoSample), "shapes")
	if err != nil {
		t.Fatal(err)
	}
	p := GeoEquirectangular().FitSize(200, 200, fc)
	out := NewPath(p).Feature(fc.Features[0])
	if out == "" || out[0] != 'M' {
		t.Errorf("decoded polygon should render a path, got %q", out)
	}
}

func assertCoords2(t *testing.T, label string, got, want [][2]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: %d points, want %d", label, len(got), len(want))
	}
	for i := range want {
		if math.Abs(got[i][0]-want[i][0]) > tol || math.Abs(got[i][1]-want[i][1]) > tol {
			t.Errorf("%s[%d] = %v, want %v", label, i, got[i], want[i])
		}
	}
}

func assertCoords3(t *testing.T, label string, got, want [][][2]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: %d rings, want %d", label, len(got), len(want))
	}
	for i := range want {
		assertCoords2(t, label, got[i], want[i])
	}
}
