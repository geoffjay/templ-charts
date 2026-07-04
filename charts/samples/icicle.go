package samples

import "github.com/geoffjay/templ-charts/charts/icicle"

// Icicle returns a small partition hierarchy, ready to assign to
// IcicleProps.Data.
func Icicle() icicle.IcicleNode {
	return icicle.IcicleNode{ID: "root", Children: []icicle.IcicleNode{
		{ID: "analytics", Children: []icicle.IcicleNode{
			{ID: "charts", Value: 40}, {ID: "reports", Value: 22},
		}},
		{ID: "billing", Children: []icicle.IcicleNode{
			{ID: "invoices", Value: 18}, {ID: "refunds", Value: 9}, {ID: "plans", Value: 15},
		}},
		{ID: "auth", Value: 28},
	}}
}
