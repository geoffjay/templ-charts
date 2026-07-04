package samples

import "github.com/geoffjay/templ-charts/charts/funnel"

// Funnel returns a five-step conversion funnel, ready to assign to
// FunnelProps.Data.
func Funnel() []funnel.FunnelDatum {
	return []funnel.FunnelDatum{
		{ID: "sent", Label: "Sent", Value: 60000},
		{ID: "viewed", Label: "Viewed", Value: 38000},
		{ID: "clicked", Label: "Clicked", Value: 22000},
		{ID: "cart", Label: "Add to Cart", Value: 12000},
		{ID: "purchased", Label: "Purchased", Value: 7000},
	}
}
