package sankey

import "testing"

// TestLayoutCyclicLinksNoPanic guards the v1.0 contract that user-supplied
// cyclic link data degrades to a best-effort layout instead of panicking.
// Sankey diagrams are acyclic, but the input is not validated upstream, and a
// panic in the layout would crash the consumer's render path.
func TestLayoutCyclicLinksNoPanic(t *testing.T) {
	cases := map[string]*Graph{
		"two-node cycle": {
			Nodes: []*Node{{ID: "a"}, {ID: "b"}},
			Links: []*Link{
				{SourceID: "a", TargetID: "b", Value: 1},
				{SourceID: "b", TargetID: "a", Value: 1},
			},
		},
		"self loop": {
			Nodes: []*Node{{ID: "a"}},
			Links: []*Link{{SourceID: "a", TargetID: "a", Value: 1}},
		},
		"three-node cycle": {
			Nodes: []*Node{{ID: "a"}, {ID: "b"}, {ID: "c"}},
			Links: []*Link{
				{SourceID: "a", TargetID: "b", Value: 1},
				{SourceID: "b", TargetID: "c", Value: 1},
				{SourceID: "c", TargetID: "a", Value: 1},
			},
		},
	}
	for name, g := range cases {
		t.Run(name, func(t *testing.T) {
			// Must not panic and must terminate.
			New().NodeAlign(SankeyJustify).NodeWidth(12).NodePadding(4).
				Size(600, 400).Layout(g)
		})
	}
}
