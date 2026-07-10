package geo

import "testing"

// A minimal but valid quantized topology with a transform, one arc, and a
// single Polygon object — the happy-path seed the fuzzer mutates from.
const seedTopology = `{
	"type": "Topology",
	"transform": {"scale": [1, 1], "translate": [0, 0]},
	"arcs": [[[0, 0], [1, 0], [0, 1], [-1, 0], [0, -1]]],
	"objects": {
		"land": {"type": "GeometryCollection", "geometries": [
			{"type": "Polygon", "arcs": [[0]]}
		]}
	}
}`

func topologyCorpus() []string {
	return []string{
		seedTopology,
		`{}`,
		`{"type":"Topology","arcs":[],"objects":{}}`,
		`{"arcs":[[[1]]],"objects":{"x":{"type":"Polygon","arcs":[[-1]]}}}`,
		`{"objects":{"x":{"type":"GeometryCollection","geometries":[{"type":"GeometryCollection","geometries":[]}]}}}`,
		`{"arcs":[[]],"objects":{"x":{"type":"LineString","arcs":[0]}}}`,
		`not json`,
		``,
	}
}

// FuzzDecodeTopology feeds arbitrary bytes to the TopoJSON decoder. The decoder
// does heavy manual slice indexing after json.Unmarshal (arc dereferencing,
// negative-index bit-complement, delta accumulation, recursive geometry
// descent) — it must return an error rather than panic on malformed input.
func FuzzDecodeTopology(f *testing.F) {
	for _, s := range topologyCorpus() {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = DecodeTopology(data)
	})
}

// FuzzDecodeTopoJSON additionally fuzzes the object-name selector path.
func FuzzDecodeTopoJSON(f *testing.F) {
	names := []string{"land", "", "missing"}
	for _, s := range topologyCorpus() {
		for _, n := range names {
			f.Add([]byte(s), n)
		}
	}
	f.Fuzz(func(t *testing.T, data []byte, objectName string) {
		_, _ = DecodeTopoJSON(data, objectName)
	})
}
