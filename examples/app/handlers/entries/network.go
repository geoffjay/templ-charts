package entries

import (
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/network"
	"github.com/geoffjay/templ-charts/charts/render"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func init() {
	register(ChartEntry{
		Slug:        "network",
		Title:       "Network",
		Description: "Force-directed node/link graph (d3-force), deterministic fixed-tick layout.",
		Snippet: `import (
    "github.com/geoffjay/templ-charts/charts/network"
    "github.com/geoffjay/templ-charts/charts/render"
)

svg, _ := render.String(network.Network(network.NetworkProps{
    Width: 440, Height: 440,
    Nodes: []network.NetworkInputNode{
        {ID: "hub-0", Size: 20, Color: "#e8c1a0"},
        {ID: "leaf-0", Size: 10, Color: "#f47560"},
    },
    Links: []network.NetworkInputLink{
        {Source: "hub-0", Target: "leaf-0"},
    },
    LinkDistance: 90,
    Repulsivity:  120,
}))`,
		Render: func(theme *theming.Theme, palette colors.PaletteID) (string, error) {
			p := network.NetworkProps{
				Width: 440, Height: 440, Responsive: true,
				Nodes: []network.NetworkInputNode{
					{ID: "hub-0", Size: 20, Color: "#e8c1a0"},
					{ID: "leaf-0", Size: 10, Color: "#f47560"},
					{ID: "leaf-1", Size: 10, Color: "#f1e15b"},
				},
				Links: []network.NetworkInputLink{
					{Source: "hub-0", Target: "leaf-0"},
					{Source: "hub-0", Target: "leaf-1"},
				},
				LinkDistance: 90,
				Repulsivity:  120,
				Theme:        theme,
			}
			return render.String(network.Network(p))
		},
	})
}
