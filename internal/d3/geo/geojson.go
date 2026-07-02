package geo

// GeoJSON geometry types this port understands. Coordinates use the standard
// [lon,lat] ordering (degrees). To keep the stream traversal type-safe, a
// Geometry's Coordinates field holds one of these concrete shapes depending on
// Type:
//
//	Point                        → [2]float64
//	MultiPoint / LineString      → [][2]float64
//	MultiLineString / Polygon    → [][][2]float64
//	MultiPolygon                 → [][][][2]float64
//	Sphere / GeometryCollection  → nil (Geometries used for collections)
//
// Callers decoding upstream GeoJSON convert into these shapes (see
// charts/geo). This mirrors d3-geo, which streams any GeoJSON object.
const (
	TypePoint             = "Point"
	TypeMultiPoint        = "MultiPoint"
	TypeLineString        = "LineString"
	TypeMultiLineString   = "MultiLineString"
	TypePolygon           = "Polygon"
	TypeMultiPolygon      = "MultiPolygon"
	TypeSphere            = "Sphere"
	TypeGeometryColl      = "GeometryCollection"
	TypeFeature           = "Feature"
	TypeFeatureCollection = "FeatureCollection"
)

// Geometry is a GeoJSON geometry (or Sphere). See the package constants for the
// concrete type held in Coordinates per Type.
type Geometry struct {
	Type        string
	Coordinates any
	Geometries  []Geometry // GeometryCollection only
}

// Feature is a GeoJSON Feature wrapping a single Geometry.
type Feature struct {
	Type       string
	ID         string
	Properties map[string]any
	Geometry   Geometry
}

// FeatureCollection is a GeoJSON FeatureCollection.
type FeatureCollection struct {
	Type     string
	Features []Feature
}

// Sink receives a stream of projected/clipped geometry primitives. Mirrors
// d3-geo's stream listener contract; sinks that ignore a callback embed noopSink.
type Sink interface {
	Point(x, y float64)
	LineStart()
	LineEnd()
	PolygonStart()
	PolygonEnd()
	Sphere()
}

// noopSink provides no-op implementations so a sink can override only the
// callbacks it cares about (embed it and shadow the rest).
type noopSink struct{}

func (noopSink) Point(x, y float64) {}
func (noopSink) LineStart()         {}
func (noopSink) LineEnd()           {}
func (noopSink) PolygonStart()      {}
func (noopSink) PolygonEnd()        {}
func (noopSink) Sphere()            {}
