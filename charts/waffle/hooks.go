package waffle

import (
	"math"
	"strconv"

	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/grid"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
)

// UseWaffle mirrors @nivo/waffle's useWaffle: it assigns each datum a half-open
// cell range based on value/unit, lays out the cell grid via
// grid.GenerateGrid, and merges the two so each cell knows its datum (or
// renders empty). props.Width/Height are the inner dimensions.
func UseWaffle(props WaffleProps) WaffleResult {
	theme := resolveTheme(props.Theme)

	cellCount := props.Rows * props.Columns
	unit := 0.0
	if cellCount > 0 {
		unit = props.Total / float64(cellCount)
	}

	getColor := colors.GetOrdinalColorScale[WaffleDatum](props.Colors, func(d WaffleDatum) string { return d.ID })
	borderGen := colors.GetInheritedColorGenerator(props.BorderColor, theme)
	format := valueFormatter(props.ValueFormat)
	hidden := map[string]bool{}
	for _, id := range props.HiddenIDs {
		hidden[id] = true
	}

	computed := make([]ComputedDatum, 0, len(props.Data))
	pos := 0
	for i, d := range props.Data {
		isHidden := hidden[d.ID]
		startAt := pos
		endAt := startAt
		if !isHidden && unit > 0 {
			endAt = startAt + int(math.Round(d.Value/unit))
		}
		pos = endAt
		color := d.Color
		if color == "" {
			color = getColor(d)
		}
		cd := ComputedDatum{
			ID:             d.ID,
			Label:          labelOr(d.Label, d.ID),
			Value:          d.Value,
			FormattedValue: format(d.Value),
			GroupIndex:     i,
			StartAt:        startAt,
			EndAt:          endAt,
			Color:          color,
			Hidden:         isHidden,
		}
		cd.BorderColor = borderGen(map[string]any{"color": color})
		computed = append(computed, cd)
	}

	// Lay out the empty cell grid (square cells, centered, fill order set by
	// FillDirection — each cell's Index is its fill position).
	g := grid.GenerateGrid(grid.GenerateGridArgs{
		Width:         props.Width,
		Height:        props.Height,
		Columns:       props.Columns,
		Rows:          props.Rows,
		Padding:       props.Padding,
		FillDirection: props.FillDirection,
		Square:        true,
	})

	cells := make([]WaffleCell, 0, len(g.Cells))
	for _, gc := range g.Cells {
		cell := WaffleCell{
			Key:         gc.Key,
			Index:       gc.Index,
			X:           gc.X,
			Y:           gc.Y,
			Width:       gc.Width,
			Height:      gc.Height,
			Color:       props.EmptyColor,
			Opacity:     props.EmptyOpacity,
			BorderColor: borderGen(map[string]any{"color": props.EmptyColor}),
		}
		// Assign to the datum whose [StartAt, EndAt) range contains this index.
		for _, cd := range computed {
			if cd.Hidden {
				continue
			}
			if gc.Index >= cd.StartAt && gc.Index < cd.EndAt {
				cell.Color = cd.Color
				cell.Opacity = 1
				cell.BorderColor = cd.BorderColor
				cell.DatumID = cd.ID
				cell.Label = cd.Label
				cell.HasData = true
				break
			}
		}
		cells = append(cells, cell)
	}

	return WaffleResult{
		Cells:        cells,
		ComputedData: computed,
		GridX:        g.X,
		GridY:        g.Y,
	}
}

func labelOr(label, fallback string) string {
	if label == "" {
		return fallback
	}
	return label
}

// valueFormatter returns a number→string formatter. Empty spec → %g.
func valueFormatter(spec string) func(float64) string {
	if spec == "" {
		return func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	}
	return func(v float64) string { return d3format.FormatString(spec, v) }
}

// resolveTheme returns props.Theme or the default theme when nil.
func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

// themeBackground returns the theme background color (or "" for transparent).
func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}
