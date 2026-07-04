package samples

import "github.com/geoffjay/templ-charts/charts/sankey"

// Sankey returns a small energy-flow graph: coal, gas, and solar feed a shared
// grid (solar also feeds homes directly) which splits into homes and industry.
// The two returned slices are the node and link inputs a sankey.SankeyProps
// expects (Nodes, Links).
func Sankey() ([]sankey.SankeyInputNode, []sankey.SankeyInputLink) {
	nodes := []sankey.SankeyInputNode{
		{ID: "Coal"},
		{ID: "Gas"},
		{ID: "Grid"},
		{ID: "Solar"},
		{ID: "Homes"},
		{ID: "Industry"},
	}
	links := []sankey.SankeyInputLink{
		{Source: "Coal", Target: "Grid", Value: 12},
		{Source: "Gas", Target: "Grid", Value: 8},
		{Source: "Solar", Target: "Grid", Value: 5},
		{Source: "Solar", Target: "Homes", Value: 3},
		{Source: "Grid", Target: "Homes", Value: 15},
		{Source: "Grid", Target: "Industry", Value: 10},
	}
	return nodes, links
}
