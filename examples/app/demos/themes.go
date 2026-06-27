package demos

import (
	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/htmx"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// ThemeDemo is one themed chart tile on the /themes page. Each tile renders
// one chart kind under one named theme.
type ThemeDemo struct {
	ID    string
	Title string
	Theme *theming.Theme
	Kind  htmx.ChartKind
	Props any
}

// ThemeGroup holds one theme name + its built Theme, used to build the per-
// kind tiles.
type ThemeGroup struct {
	Name  string
	Theme *theming.Theme
}

// ThemeGroups returns the dark/light/custom theme variants shown on the
// /themes page.
func ThemeGroups() []ThemeGroup {
	bg := "rgb(24, 32, 38)"
	lightText := "#d1d9e0"
	return []ThemeGroup{
		{Name: "default", Theme: &theming.DefaultTheme},
		{Name: "dark", Theme: darkTheme(bg, lightText)},
		{Name: "custom", Theme: customTheme()},
	}
}

// ThemesDemos returns the cartesian product of (theme group) × (chart kind)
// for the /themes page: one ThemeDemo per (theme, kind) pair.
func ThemesDemos() []ThemeDemo {
	groups := ThemeGroups()
	var out []ThemeDemo
	for _, g := range groups {
		// Bar
		bp := bar.BarProps{
			Width: commonChartWidth, Height: commonChartHeight,
			IndexBy: "country",
			Keys:    []string{"hot dogs", "burgers", "sandwich", "kebab"},
			Data:    barData(), Margin: defaultMargin(),
			Theme: g.Theme,
		}
		out = append(out, ThemeDemo{
			ID: "themes-bar-" + g.Name, Title: "bar / " + g.Name,
			Theme: g.Theme, Kind: htmx.KindBar, Props: bp,
		})
		// Line
		lp := line.LineProps{
			Width: commonChartWidth, Height: commonChartHeight,
			Margin: defaultMargin(), Curve: core.CurveMonotoneX,
			Data: lineData(), Theme: g.Theme,
		}
		out = append(out, ThemeDemo{
			ID: "themes-line-" + g.Name, Title: "line / " + g.Name,
			Theme: g.Theme, Kind: htmx.KindLine, Props: lp,
		})
		// Pie
		pp := pie.PieProps{
			Width: commonChartWidth, Height: commonChartHeight,
			Margin: defaultMargin(), InnerRadius: 0.5, PadAngle: 0.5,
			Data: pieData(), Theme: g.Theme,
		}
		out = append(out, ThemeDemo{
			ID: "themes-pie-" + g.Name, Title: "pie / " + g.Name,
			Theme: g.Theme, Kind: htmx.KindPie, Props: pp,
		})
	}
	return out
}

// darkTheme builds a dark-background theme by extending the default theme
// with a dark background + light text.
func darkTheme(background, textFill string) *theming.Theme {
	gridLine := theming.GridLine{Extra: map[string]any{"stroke": "rgba(255,255,255,0.15)", "strokeWidth": float64(1)}}
	t := theming.ExtendDefaultTheme(theming.DefaultTheme, &theming.PartialTheme{
		Background: &background,
		Text:       &theming.TextStyle{Fill: textFill, FontFamily: "sans-serif"},
		Grid:       &theming.PartialGridTheme{Line: &gridLine},
	})
	return &t
}

// customTheme builds a warm custom theme with a cream background and serif
// text.
func customTheme() *theming.Theme {
	bg := "#fff8f0"
	t := theming.ExtendDefaultTheme(theming.DefaultTheme, &theming.PartialTheme{
		Background: &bg,
		Text:       &theming.TextStyle{Fill: "#7a4a2b", FontFamily: "Georgia, serif"},
	})
	return &t
}

// ThemesLegendDemosNotUsed is a placeholder so legends import isn't dropped
// if future refactors remove direct usage. Kept intentional: the themes page
// doesn't use legends, but the import documents the available option.
var _ = legends.LegendAnchorTopRight
