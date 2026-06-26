package core

// SvgWrapperProps mirrors @nivo/core's SvgWrapper props. Width/Height are
// the outer svg dimensions; Margin offsets the inner <g> translation; Defs
// is the bound defs list to render via the Defs component; Background, when
// non-empty, sets the background rect fill (overrides theme.background).
//
// Aria attributes and Role match nivo's defaults (role="img", focusable=false
// unless IsFocusable).
type SvgWrapperProps struct {
	Width           float64
	Height          float64
	Margin          Margin
	Defs            []Def
	Background      string
	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	IsFocusable     bool
}

// DotsItemProps mirrors nivo's DotsItem props (minus spring animation, which
// is replaced by an optional SMIL `<animate>`).
type DotsItemProps struct {
	X, Y            float64
	Size            float64
	Color           string
	BorderWidth     float64
	BorderColor     string
	Label           string
	LabelTextAnchor string
	LabelYOffset    float64
	LabelFill       string
	LabelFontSize   float64
	LabelFontFamily string
}

// CartesianMarker is a marker spec (axis x/y, value, optional legend).
type CartesianMarker struct {
	Axis              string // "x" | "y"
	Value             any
	Legend            string
	LegendPosition    string // defaults to "top-right"
	LegendOffsetX     float64
	LegendOffsetY     float64
	LegendOrientation string // "horizontal" | "vertical", defaults to "horizontal"
	LineStyle         map[string]any
	TextStyle         map[string]any
	LineColor         string
	LineStrokeWidth   float64
}

// CartesianMarkersProps is the input to the CartesianMarkers component.
type CartesianMarkersProps struct {
	Markers []CartesianMarker
	Width   float64
	Height  float64
	XScale  func(any) float64
	YScale  func(any) float64
}
