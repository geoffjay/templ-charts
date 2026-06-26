// Package theming mirrors @nivo/theming: the Theme struct, defaultTheme
// (verbatim from defaults.ts), ExtendDefaultTheme (deep merge + 9-path text
// inheritance), ExtendAxisTheme, borderRadius helpers, the engine bridge
// (svgStyleAttributesMapping/convertStyleAttribute), and SanitizeSvgTextStyle
// / SanitizeHtmlTextStyle.
package theming

// TextStyle mirrors @nivo/theming TextStyle. Fill/FontFamily/FontSize etc.
// are the typed core; Extra carries any additional CSS properties.
type TextStyle struct {
	FontFamily     string
	FontSize       any // number (px) or string (e.g. "11px")
	Fill           string
	OutlineWidth   float64
	OutlineColor   string
	OutlineOpacity float64
	Extra          map[string]any // additional CSS properties (fontWeight, fontStyle, …)
}

// AxisDomainLine is the axis domain line style.
type AxisDomainLine struct {
	Extra map[string]any
}

// AxisTicks is the axis ticks style (line + text).
type AxisTicks struct {
	Line AxisTickLine
	Text TextStyle
}

// AxisTickLine is the axis tick line style.
type AxisTickLine struct {
	Extra map[string]any
}

// AxisLegend is the axis legend style.
type AxisLegend struct {
	Text TextStyle
}

// AxisTheme is the per-axis theme block.
type AxisTheme struct {
	Domain AxisDomain
	Ticks  AxisTicks
	Legend AxisLegend
}

// AxisDomain holds the axis domain line.
type AxisDomain struct {
	Line AxisDomainLine
}

// GridTheme is the grid line style.
type GridTheme struct {
	Line GridLine
}

// GridLine is the grid line style.
type GridLine struct {
	Extra map[string]any
}

// CrosshairLine is the crosshair line style.
type CrosshairLine struct {
	Stroke          string
	StrokeWidth     float64
	StrokeOpacity   float64
	StrokeDasharray string
}

// CrosshairTheme is the crosshair theme block.
type CrosshairTheme struct {
	Line CrosshairLine
}

// LegendsHidden is the legend hidden-item style.
type LegendsHidden struct {
	Symbol LegendsHiddenSymbol
	Text   TextStyle
}

// LegendsHiddenSymbol is the hidden-legend symbol style.
type LegendsHiddenSymbol struct {
	Fill    string
	Opacity float64
}

// LegendsTicks is the legend ticks style.
type LegendsTicks struct {
	Line AxisTickLine
	Text TextStyle
}

// LegendsTitle is the legend title style.
type LegendsTitle struct {
	Text TextStyle
}

// LegendsTheme is the legends theme block.
type LegendsTheme struct {
	Hidden LegendsHidden
	Text   TextStyle
	Title  LegendsTitle
	Ticks  LegendsTicks
}

// LabelsTheme is the labels theme block.
type LabelsTheme struct {
	Text TextStyle
}

// MarkersTheme is the markers theme block.
type MarkersTheme struct {
	LineColor       string
	LineStrokeWidth float64
	Text            TextStyle
}

// DotsTheme is the dots theme block.
type DotsTheme struct {
	Text TextStyle
}

// TooltipTheme is the tooltip theme block (CSS-style values for HTML).
type TooltipTheme struct {
	Container      map[string]any
	Basic          map[string]any
	Chip           map[string]any
	Table          map[string]any
	TableCell      map[string]any
	TableCellValue map[string]any
}

// AnnotationLink is the annotation link style.
type AnnotationLink struct {
	Stroke         string
	StrokeWidth    float64
	OutlineWidth   float64
	OutlineColor   string
	OutlineOpacity float64
	Extra          map[string]any
}

// AnnotationOutline is the annotation outline style.
type AnnotationOutline struct {
	Stroke         string
	StrokeWidth    float64
	OutlineWidth   float64
	OutlineColor   string
	OutlineOpacity float64
	Extra          map[string]any
}

// AnnotationSymbol is the annotation symbol style.
type AnnotationSymbol struct {
	Fill           string
	OutlineWidth   float64
	OutlineColor   string
	OutlineOpacity float64
	Extra          map[string]any
}

// AnnotationsTheme is the annotations theme block.
type AnnotationsTheme struct {
	Text    TextStyle
	Link    AnnotationLink
	Outline AnnotationOutline
	Symbol  AnnotationSymbol
}

// Theme is the full theme after inheritance has been applied. Mirrors
// @nivo/theming Theme.
type Theme struct {
	Background  string
	Text        TextStyle
	Axis        AxisTheme
	Grid        GridTheme
	Crosshair   CrosshairTheme
	Legends     LegendsTheme
	Labels      LabelsTheme
	Markers     MarkersTheme
	Dots        DotsTheme
	Tooltip     TooltipTheme
	Annotations AnnotationsTheme
}

// ThemeWithoutInheritance is the input form of a theme (nested text styles
// may omit inherited root text properties). Mirrors ThemeWithoutInheritance.
type ThemeWithoutInheritance = Theme

// PartialTheme is a theme where every field is optional. Used as input to
// ExtendDefaultTheme. Modeled as a pointer-bearing Theme so callers can
// express "unset" via nil pointers.
type PartialTheme struct {
	Background  *string
	Text        *TextStyle
	Axis        *PartialAxisTheme
	Grid        *PartialGridTheme
	Crosshair   *PartialCrosshairTheme
	Legends     *PartialLegendsTheme
	Labels      *PartialLabelsTheme
	Markers     *PartialMarkersTheme
	Dots        *PartialDotsTheme
	Tooltip     *PartialTooltipTheme
	Annotations *PartialAnnotationsTheme
}

// PartialAxisTheme is the partial form of AxisTheme.
type PartialAxisTheme struct {
	Domain *PartialAxisDomain
	Ticks  *PartialAxisTicks
	Legend *PartialAxisLegend
}

// PartialAxisDomain is the partial form of AxisDomain.
type PartialAxisDomain struct {
	Line *AxisDomainLine
}

// PartialAxisTicks is the partial form of AxisTicks.
type PartialAxisTicks struct {
	Line *AxisTickLine
	Text *TextStyle
}

// PartialAxisLegend is the partial form of AxisLegend.
type PartialAxisLegend struct {
	Text *TextStyle
}

// PartialGridTheme is the partial form of GridTheme.
type PartialGridTheme struct {
	Line *GridLine
}

// PartialCrosshairTheme is the partial form of CrosshairTheme.
type PartialCrosshairTheme struct {
	Line *CrosshairLine
}

// PartialLegendsTheme is the partial form of LegendsTheme.
type PartialLegendsTheme struct {
	Hidden *PartialLegendsHidden
	Text   *TextStyle
	Title  *PartialLegendsTitle
	Ticks  *PartialLegendsTicks
}

// PartialLegendsHidden is the partial form of LegendsHidden.
type PartialLegendsHidden struct {
	Symbol *LegendsHiddenSymbol
	Text   *TextStyle
}

// PartialLegendsTitle is the partial form of LegendsTitle.
type PartialLegendsTitle struct {
	Text *TextStyle
}

// PartialLegendsTicks is the partial form of LegendsTicks.
type PartialLegendsTicks struct {
	Line *AxisTickLine
	Text *TextStyle
}

// PartialLabelsTheme is the partial form of LabelsTheme.
type PartialLabelsTheme struct {
	Text *TextStyle
}

// PartialMarkersTheme is the partial form of MarkersTheme.
type PartialMarkersTheme struct {
	LineColor       *string
	LineStrokeWidth *float64
	Text            *TextStyle
}

// PartialDotsTheme is the partial form of DotsTheme.
type PartialDotsTheme struct {
	Text *TextStyle
}

// PartialTooltipTheme is the partial form of TooltipTheme.
type PartialTooltipTheme struct {
	Container      map[string]any
	Basic          map[string]any
	Chip           map[string]any
	Table          map[string]any
	TableCell      map[string]any
	TableCellValue map[string]any
}

// PartialAnnotationsTheme is the partial form of AnnotationsTheme.
type PartialAnnotationsTheme struct {
	Text    *TextStyle
	Link    *AnnotationLink
	Outline *AnnotationOutline
	Symbol  *AnnotationSymbol
}
