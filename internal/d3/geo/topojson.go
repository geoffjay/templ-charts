package geo

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Optional TopoJSON decoding. Callers may still supply GeoJSON directly (the
// default input path); this adds an opt-in decoder so a TopoJSON topology can
// be turned into the same geo.Feature / FeatureCollection the rest of the
// package consumes. Ported from topojson-client's feature()/object() (delta
// transform + arc dereferencing). Only the subset needed to render maps is
// implemented: GeometryCollection of Polygon/MultiPolygon/LineString/
// MultiLineString/Point/MultiPoint (nested collections included).

// Topology is a decoded TopoJSON topology. Arcs are shared, delta-encoded
// (when Transform is set) position lists; object geometries reference them by
// index. Decode one with DecodeTopology, then dereference an object into
// GeoJSON with FeatureCollection.
type Topology struct {
	Type      string                  `json:"type"`
	Transform *TopoTransform          `json:"transform,omitempty"`
	Arcs      [][][]float64           `json:"arcs"`
	Objects   map[string]TopoGeometry `json:"objects"`
}

// TopoTransform is the quantization transform: absolute = delta*scale+translate.
type TopoTransform struct {
	Scale     [2]float64 `json:"scale"`
	Translate [2]float64 `json:"translate"`
}

// TopoGeometry is a TopoJSON geometry object. The Arcs/Coordinates shapes vary
// by Type, so they are held as raw JSON and decoded on demand.
type TopoGeometry struct {
	Type        string          `json:"type"`
	ID          any             `json:"id,omitempty"`
	Properties  map[string]any  `json:"properties,omitempty"`
	Arcs        json.RawMessage `json:"arcs,omitempty"`
	Coordinates json.RawMessage `json:"coordinates,omitempty"`
	Geometries  []TopoGeometry  `json:"geometries,omitempty"`
}

// DecodeTopology parses TopoJSON bytes into a Topology.
func DecodeTopology(data []byte) (*Topology, error) {
	var t Topology
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("geo: decode topojson: %w", err)
	}
	if t.Type != "" && t.Type != "Topology" {
		return nil, fmt.Errorf("geo: not a TopoJSON Topology (type %q)", t.Type)
	}
	return &t, nil
}

// DecodeTopoJSON is a one-shot convenience: parse a topology and dereference the
// named object into a FeatureCollection.
func DecodeTopoJSON(data []byte, objectName string) (*FeatureCollection, error) {
	t, err := DecodeTopology(data)
	if err != nil {
		return nil, err
	}
	return t.FeatureCollection(objectName)
}

// FeatureCollection dereferences the named object into a GeoJSON
// FeatureCollection (a single non-collection object yields a one-feature
// collection). Mirrors topojson.feature.
func (t *Topology) FeatureCollection(objectName string) (*FeatureCollection, error) {
	o, ok := t.Objects[objectName]
	if !ok {
		return nil, fmt.Errorf("geo: topojson object %q not found", objectName)
	}
	fc := &FeatureCollection{Type: TypeFeatureCollection}
	if o.Type == "GeometryCollection" {
		fc.Features = make([]Feature, 0, len(o.Geometries))
		for _, geom := range o.Geometries {
			fc.Features = append(fc.Features, t.toFeature(geom))
		}
		return fc, nil
	}
	fc.Features = []Feature{t.toFeature(o)}
	return fc, nil
}

func (t *Topology) toFeature(o TopoGeometry) Feature {
	return Feature{
		Type:       TypeFeature,
		ID:         idString(o.ID),
		Properties: o.Properties,
		Geometry:   t.geometry(o),
	}
}

// geometry converts a TopoJSON geometry object into a geo.Geometry.
func (t *Topology) geometry(o TopoGeometry) Geometry {
	switch o.Type {
	case "GeometryCollection":
		subs := make([]Geometry, len(o.Geometries))
		for i, s := range o.Geometries {
			subs[i] = t.geometry(s)
		}
		return Geometry{Type: TypeGeometryColl, Geometries: subs}
	case "Point":
		return Geometry{Type: TypePoint, Coordinates: t.transformStandalone(decodePoint(o.Coordinates))}
	case "MultiPoint":
		pts := decodePoints(o.Coordinates)
		out := make([][2]float64, len(pts))
		for i, p := range pts {
			out[i] = t.transformStandalone(p)
		}
		return Geometry{Type: TypeMultiPoint, Coordinates: out}
	case "LineString":
		return Geometry{Type: TypeLineString, Coordinates: t.line(decodeArcs1(o.Arcs))}
	case "MultiLineString":
		idx := decodeArcs2(o.Arcs)
		lines := make([][][2]float64, len(idx))
		for i, l := range idx {
			lines[i] = t.line(l)
		}
		return Geometry{Type: TypeMultiLineString, Coordinates: lines}
	case "Polygon":
		return Geometry{Type: TypePolygon, Coordinates: t.polygon(decodeArcs2(o.Arcs))}
	case "MultiPolygon":
		idx := decodeArcs3(o.Arcs)
		polys := make([][][][2]float64, len(idx))
		for i, poly := range idx {
			polys[i] = t.polygon(poly)
		}
		return Geometry{Type: TypeMultiPolygon, Coordinates: polys}
	}
	return Geometry{Type: o.Type}
}

func (t *Topology) polygon(rings [][]int) [][][2]float64 {
	out := make([][][2]float64, len(rings))
	for i, r := range rings {
		out[i] = t.ring(r)
	}
	return out
}

// ring is line() padded to at least 4 points (a closed ring), matching
// topojson-client.
func (t *Topology) ring(arcIdx []int) [][2]float64 {
	pts := t.line(arcIdx)
	for len(pts) < 4 && len(pts) > 0 {
		pts = append(pts, pts[0])
	}
	return pts
}

// line concatenates the referenced arcs, dropping the vertex shared between
// consecutive arcs, matching topojson-client.
func (t *Topology) line(arcIdx []int) [][2]float64 {
	var pts [][2]float64
	for _, ai := range arcIdx {
		t.decodeArcInto(ai, &pts)
	}
	if len(pts) < 2 && len(pts) > 0 {
		pts = append(pts, pts[0])
	}
	return pts
}

// decodeArcInto appends arc i's positions to pts (popping the previously shared
// endpoint first); a negative index references the arc reversed (~i).
func (t *Topology) decodeArcInto(i int, pts *[][2]float64) {
	if len(*pts) > 0 {
		*pts = (*pts)[:len(*pts)-1]
	}
	idx := i
	reversed := false
	if i < 0 {
		idx = ^i // ~i == -i-1
		reversed = true
	}
	if idx < 0 || idx >= len(t.Arcs) {
		return
	}
	a := t.Arcs[idx]
	start := len(*pts)
	var acc [2]float64
	for k, p := range a {
		*pts = append(*pts, t.transformPos(p, k, &acc))
	}
	if reversed {
		reverseSub(*pts, start, len(*pts))
	}
}

// transformPos applies the delta transform for a position at index k within an
// arc (the accumulator resets at k==0).
func (t *Topology) transformPos(p []float64, k int, acc *[2]float64) [2]float64 {
	if len(p) < 2 {
		return [2]float64{}
	}
	if t.Transform == nil {
		return [2]float64{p[0], p[1]}
	}
	if k == 0 {
		acc[0], acc[1] = 0, 0
	}
	acc[0] += p[0]
	acc[1] += p[1]
	return [2]float64{
		acc[0]*t.Transform.Scale[0] + t.Transform.Translate[0],
		acc[1]*t.Transform.Scale[1] + t.Transform.Translate[1],
	}
}

// transformStandalone applies the transform to a standalone (non-arc) point,
// which is absolute-quantized (no delta accumulation).
func (t *Topology) transformStandalone(p []float64) [2]float64 {
	if len(p) < 2 {
		return [2]float64{}
	}
	if t.Transform == nil {
		return [2]float64{p[0], p[1]}
	}
	return [2]float64{
		p[0]*t.Transform.Scale[0] + t.Transform.Translate[0],
		p[1]*t.Transform.Scale[1] + t.Transform.Translate[1],
	}
}

func reverseSub(a [][2]float64, from, to int) {
	for i, j := from, to-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func idString(id any) string {
	switch v := id.(type) {
	case nil:
		return ""
	case string:
		return v
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// --- raw arc/coordinate decoders -------------------------------------------

func decodeArcs1(raw json.RawMessage) []int {
	var v []int
	_ = json.Unmarshal(raw, &v)
	return v
}

func decodeArcs2(raw json.RawMessage) [][]int {
	var v [][]int
	_ = json.Unmarshal(raw, &v)
	return v
}

func decodeArcs3(raw json.RawMessage) [][][]int {
	var v [][][]int
	_ = json.Unmarshal(raw, &v)
	return v
}

func decodePoint(raw json.RawMessage) []float64 {
	var v []float64
	_ = json.Unmarshal(raw, &v)
	return v
}

func decodePoints(raw json.RawMessage) [][]float64 {
	var v [][]float64
	_ = json.Unmarshal(raw, &v)
	return v
}
