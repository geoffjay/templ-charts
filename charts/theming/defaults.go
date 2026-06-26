package theming

// DefaultTheme is @nivo/theming's defaultTheme, transcribed verbatim from
// defaults.ts. All field values match the upstream defaults.
var DefaultTheme = Theme{
	Background: "transparent",
	Text: TextStyle{
		FontFamily:     "sans-serif",
		FontSize:       11,
		Fill:           "#333333",
		OutlineWidth:   0,
		OutlineColor:   "#ffffff",
		OutlineOpacity: 1,
	},
	Axis: AxisTheme{
		Domain: AxisDomain{
			Line: AxisDomainLine{
				Extra: map[string]any{"stroke": "transparent", "strokeWidth": float64(1)},
			},
		},
		Ticks: AxisTicks{
			Line: AxisTickLine{
				Extra: map[string]any{"stroke": "#777777", "strokeWidth": float64(1)},
			},
			Text: TextStyle{},
		},
		Legend: AxisLegend{
			Text: TextStyle{FontSize: 12},
		},
	},
	Grid: GridTheme{
		Line: GridLine{
			Extra: map[string]any{"stroke": "#dddddd", "strokeWidth": float64(1)},
		},
	},
	Legends: LegendsTheme{
		Hidden: LegendsHidden{
			Symbol: LegendsHiddenSymbol{Fill: "#333333", Opacity: 0.6},
			Text:   TextStyle{Fill: "#333333", OutlineOpacity: 0, Extra: map[string]any{"opacity": 0.6}},
		},
		Text: TextStyle{},
		Ticks: LegendsTicks{
			Line: AxisTickLine{
				Extra: map[string]any{"stroke": "#777777", "strokeWidth": float64(1)},
			},
			Text: TextStyle{FontSize: 10},
		},
		Title: LegendsTitle{Text: TextStyle{}},
	},
	Labels: LabelsTheme{Text: TextStyle{}},
	Markers: MarkersTheme{
		LineColor:       "#000000",
		LineStrokeWidth: 1,
		Text:            TextStyle{},
	},
	Dots: DotsTheme{Text: TextStyle{}},
	Tooltip: TooltipTheme{
		Container: map[string]any{
			"background":   "white",
			"color":        "inherit",
			"fontSize":     "inherit",
			"borderRadius": "2px",
			"boxShadow":    "0 1px 2px rgba(0, 0, 0, 0.25)",
			"padding":      "5px 9px",
		},
		Basic: map[string]any{
			"whiteSpace": "pre",
			"display":    "flex",
			"alignItems": "center",
		},
		Chip: map[string]any{
			"marginRight": 7,
		},
		Table:          map[string]any{},
		TableCell:      map[string]any{"padding": "3px 5px"},
		TableCellValue: map[string]any{"fontWeight": "bold"},
	},
	Crosshair: CrosshairTheme{
		Line: CrosshairLine{
			Stroke:          "#000000",
			StrokeWidth:     1,
			StrokeOpacity:   0.75,
			StrokeDasharray: "6 6",
		},
	},
	Annotations: AnnotationsTheme{
		Text: TextStyle{
			FontSize:       13,
			OutlineWidth:   2,
			OutlineColor:   "#ffffff",
			OutlineOpacity: 1,
		},
		Link: AnnotationLink{
			Stroke:         "#000000",
			StrokeWidth:    1,
			OutlineWidth:   2,
			OutlineColor:   "#ffffff",
			OutlineOpacity: 1,
		},
		Outline: AnnotationOutline{
			Stroke:         "#000000",
			StrokeWidth:    2,
			OutlineWidth:   2,
			OutlineColor:   "#ffffff",
			OutlineOpacity: 1,
			Extra:          map[string]any{"fill": "none"},
		},
		Symbol: AnnotationSymbol{
			Fill:           "#000000",
			OutlineWidth:   2,
			OutlineColor:   "#ffffff",
			OutlineOpacity: 1,
		},
	},
}
