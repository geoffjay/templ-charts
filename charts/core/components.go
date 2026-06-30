package core

// SvgWrapperProps mirrors @nivo/core's SvgWrapper props. Width/Height are
// the outer svg dimensions; Margin offsets the inner <g> translation; Defs
// is the bound defs list to render via the Defs component; Background, when
// non-empty, sets the background rect fill (overrides theme.background).
//
// Aria attributes and Role match nivo's defaults (role="img", focusable=false
// unless IsFocusable).
//
// Responsive, when true, makes the svg scale fluidly to its container: the
// viewBox (and thus aspect ratio) is preserved while an inline
// width:100%;height:auto style overrides the fixed pixel dimensions. The
// Width/Height are still emitted as the intrinsic size so the aspect ratio is
// well-defined and non-CSS contexts get a sensible fallback.
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
	Responsive      bool
	// Title / Desc, when non-empty, render <title>/<desc> as the first
	// children of the <svg>. These are the SVG-native accessibility mechanism:
	// for role="img", assistive tech derives the accessible name from <title>
	// and the description from <desc>. They complement the aria-* attributes
	// (an explicit aria-label still wins for the name) and double as the
	// browser's native hover tooltip for <title>.
	Title string
	Desc  string
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
	// Tooltip, when non-empty, is the HTML shown by the client interactivity
	// layer (charts/interact) on hover: it is emitted as a data-tc-tooltip
	// attribute and the dot becomes pointer-events:auto so it is hoverable.
	Tooltip string
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
