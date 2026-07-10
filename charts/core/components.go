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
	Width           float64 // outer svg width, in px
	Height          float64 // outer svg height, in px
	Margin          Margin  // offsets the inner <g> translation
	Defs            []Def   // bound defs list rendered via the Defs component
	Background      string  // background rect fill; empty uses theme.background
	Role            string  // svg role attribute (defaults to "img")
	AriaLabel       string  // aria-label accessible name (wins over Title)
	AriaLabelledBy  string  // aria-labelledby id reference(s)
	AriaDescribedBy string  // aria-describedby id reference(s)
	IsFocusable     bool    // when true, sets focusable="true" (default false)
	Responsive      bool    // when true, svg scales fluidly to its container via viewBox
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
	X, Y            float64 // dot center position, in svg units
	Size            float64 // dot diameter, in px
	Color           string  // dot fill color
	BorderWidth     float64 // dot border stroke width, in px
	BorderColor     string  // dot border stroke color
	Label           string  // optional text label rendered beside the dot
	LabelTextAnchor string  // label text-anchor ("start" | "middle" | "end")
	LabelYOffset    float64 // vertical offset of the label from the dot center, in px
	LabelFill       string  // label text color
	LabelFontSize   float64 // label font size, in px
	LabelFontFamily string  // label font family
	// Tooltip, when non-empty, is the HTML shown by the client interactivity
	// layer (charts/interact) on hover: it is emitted as a data-tc-tooltip
	// attribute and the dot becomes pointer-events:auto so it is hoverable.
	Tooltip string
	// Animate, when true, emits a SMIL enter animation scaling the dot's radius
	// from 0 to its final value (600ms). AnimateBegin is the per-dot start
	// offset (for staggering). Both default off/empty, so static output is
	// byte-identical. See core.SMILAnimate.
	Animate      bool
	AnimateBegin string
}

// CartesianMarker is a marker spec (axis x/y, value, optional legend).
type CartesianMarker struct {
	Axis              string         // "x" | "y"
	Value             any            // scale value at which the marker line is drawn
	Legend            string         // optional legend text; empty omits the label
	LegendPosition    string         // defaults to "top-right"
	LegendOffsetX     float64        // legend x offset from its anchor, in px (defaults to 14)
	LegendOffsetY     float64        // legend y offset from its anchor, in px (defaults to 14)
	LegendOrientation string         // "horizontal" | "vertical", defaults to "horizontal"
	LineStyle         map[string]any // extra SVG style attributes for the marker line
	TextStyle         map[string]any // extra SVG style attributes for the legend text
	LineColor         string         // marker line stroke color
	LineStrokeWidth   float64        // marker line stroke width, in px
}

// CartesianMarkersProps is the input to the CartesianMarkers component.
type CartesianMarkersProps struct {
	Markers []CartesianMarker // markers to render
	Width   float64           // inner plotting-area width, in px
	Height  float64           // inner plotting-area height, in px
	XScale  func(any) float64 // maps an x-axis value to an x pixel position
	YScale  func(any) float64 // maps a y-axis value to a y pixel position
}
