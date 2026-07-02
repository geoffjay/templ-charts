package chord

import (
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3chord "github.com/geoffjay/templ-charts/internal/d3/chord"
	d3format "github.com/geoffjay/templ-charts/internal/d3/format"
)

// UseChord mirrors @nivo/chord's useChord: it runs the d3-chord layout over the
// input matrix (with padAngle), derives the circle geometry from the inner
// dimensions, then builds the colored group arcs and the ribbon paths. Arc
// paths use charts/arcs (d3-shape Arc, outerRadius=radius, innerRadius=radius*
// innerRadiusRatio); ribbon paths use internal/d3/chord's Ribbon generator at
// radius*(innerRadiusRatio-innerRadiusOffset). props.Width/Height are the inner
// dimensions.
func UseChord(props ChordProps) ChordResult {
	center := [2]float64{props.Width / 2, props.Height / 2}
	radius := math.Min(props.Width, props.Height) / 2
	innerRadius := radius * props.InnerRadiusRatio
	ribbonRadius := radius * (props.InnerRadiusRatio - props.InnerRadiusOffset)

	layout := d3chord.New().PadAngle(props.PadAngle).Compute(props.Data)

	getColor := colors.GetOrdinalColorScale[string](props.Colors, func(id string) string { return id })
	format := valueFormatter(props.ValueFormat)
	arcGen := arcs.CreateArcGenerator(0, 0)

	key := func(i int) string {
		if i >= 0 && i < len(props.Keys) {
			return props.Keys[i]
		}
		return strconv.Itoa(i)
	}

	computedArcs := make([]ComputedArc, len(layout.Groups))
	for i, g := range layout.Groups {
		id := key(g.Index)
		color := getColor(id)
		ca := ComputedArc{
			ID:             id,
			Index:          g.Index,
			Label:          id,
			Value:          g.Value,
			FormattedValue: format(g.Value),
			StartAngle:     g.StartAngle,
			EndAngle:       g.EndAngle,
			Color:          color,
		}
		ca.Path = arcGen.GenerateSvgArc(arcs.Arc{
			StartAngle:  g.StartAngle,
			EndAngle:    g.EndAngle,
			InnerRadius: innerRadius,
			OuterRadius: radius,
		})
		computedArcs[i] = ca
	}

	ribbonGen := d3chord.NewRibbon(ribbonRadius)
	computedRibbons := make([]ComputedRibbon, 0, len(layout.Ribbons))
	for _, rb := range layout.Ribbons {
		src := computedArcs[rb.Source.Index]
		tgt := computedArcs[rb.Target.Index]

		// Ribbon id is stable regardless of source/target ordering.
		ids := []string{src.ID, tgt.ID}
		sort.Strings(ids)

		// Match @nivo/chord getRibbonAngles: draw from the subgroup with the
		// smaller start angle to the other, so the ribbon isn't reversed when
		// the larger-value end flips.
		first, second := rb.Source, rb.Target
		if second.StartAngle < first.StartAngle {
			first, second = second, first
		}
		path := ribbonGen.Call(d3chord.RibbonArg{
			Source: d3chord.RibbonEndpoint{StartAngle: first.StartAngle, EndAngle: first.EndAngle},
			Target: d3chord.RibbonEndpoint{StartAngle: second.StartAngle, EndAngle: second.EndAngle},
		})

		computedRibbons = append(computedRibbons, ComputedRibbon{
			ID:     strings.Join(ids, "."),
			Source: src,
			Target: tgt,
			Path:   path,
			Color:  src.Color,
		})
	}

	legendData := make([]legends.Datum, 0, len(computedArcs))
	for _, a := range computedArcs {
		legendData = append(legendData, legends.Datum{ID: a.ID, Label: a.Label, Color: a.Color})
	}

	return ChordResult{
		Arcs:       computedArcs,
		Ribbons:    computedRibbons,
		Center:     center,
		Radius:     radius,
		LegendData: legendData,
	}
}

func valueFormatter(spec string) func(float64) string {
	if spec == "" {
		return func(v float64) string { return strconv.FormatFloat(v, 'g', -1, 64) }
	}
	return func(v float64) string { return d3format.FormatString(spec, v) }
}

func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}
