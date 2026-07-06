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
		Render: func(theme *theming.Theme, palette colors.PaletteID, animate bool) (string, error) {
			// Seven service clusters: a hub node plus a few leaves each, with
			// mostly intra-cluster links and a few cross-cluster ties. Each
			// cluster gets one color; the selected palette drives them when set.
			groups := []string{"api", "web", "db", "auth", "cache", "queue", "cdn"}
			leafCounts := []int{4, 3, 3, 2, 3, 2, 3}
			clusterColors := []string{"#e8c1a0", "#f47560", "#f1e15b", "#e8a838", "#61cdbb", "#97e3d5", "#8da0cb"}
			if pal, ok := colors.LookupPalette(palette); ok {
				if sw := pal.Swatch(len(groups)); len(sw) == len(groups) {
					clusterColors = sw
				}
			}
			var nodes []network.NetworkInputNode
			var links []network.NetworkInputLink
			for gi, g := range groups {
				c := clusterColors[gi]
				nodes = append(nodes, network.NetworkInputNode{ID: g, Size: 18, Color: c})
				for i := 0; i < leafCounts[gi]; i++ {
					id := g + "-" + string(rune('0'+i))
					nodes = append(nodes, network.NetworkInputNode{ID: id, Size: 9 + float64((gi+i)%3)*2, Color: c})
					links = append(links, network.NetworkInputLink{Source: g, Target: id})
				}
			}
			// Ring of hub-to-hub links keeps clusters connected but distinct.
			for gi := range groups {
				links = append(links, network.NetworkInputLink{Source: groups[gi], Target: groups[(gi+1)%len(groups)], Distance: 130})
			}
			links = append(links,
				network.NetworkInputLink{Source: "api", Target: "db", Distance: 150},
				network.NetworkInputLink{Source: "web", Target: "cache", Distance: 150},
				network.NetworkInputLink{Source: "api", Target: "queue", Distance: 150},
			)
			p := network.NetworkProps{
				Width: 440, Height: 440, Responsive: true,
				Nodes:        nodes,
				Links:        links,
				LinkDistance: 90,
				Repulsivity:  120,
				Theme:        theme,
			}
			p.Animate = animate
			return render.String(network.Network(p))
		},
	})
}
