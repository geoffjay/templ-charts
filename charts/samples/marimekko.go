package samples

import "github.com/geoffjay/templ-charts/charts/marimekko"

// Marimekko returns ready-to-render marimekko data: four countries sized by
// sample count, each split across four survey-response dimensions. The returned
// dimensions slice maps each Marimekko datum's Dimensions keys to display ids.
func Marimekko() ([]marimekko.MarimekkoDatum, []marimekko.MarimekkoDimension) {
	dims := []marimekko.MarimekkoDimension{
		{ID: "agree strongly", Key: "agreeStrongly"},
		{ID: "agree", Key: "agree"},
		{ID: "disagree", Key: "disagree"},
		{ID: "disagree strongly", Key: "disagreeStrongly"},
	}
	data := []marimekko.MarimekkoDatum{
		{ID: "France", Value: 42, Dimensions: map[string]float64{"agreeStrongly": 18, "agree": 12, "disagree": 8, "disagreeStrongly": 4}},
		{ID: "Japan", Value: 28, Dimensions: map[string]float64{"agreeStrongly": 6, "agree": 10, "disagree": 8, "disagreeStrongly": 4}},
		{ID: "USA", Value: 63, Dimensions: map[string]float64{"agreeStrongly": 30, "agree": 18, "disagree": 10, "disagreeStrongly": 5}},
		{ID: "Germany", Value: 35, Dimensions: map[string]float64{"agreeStrongly": 12, "agree": 14, "disagree": 6, "disagreeStrongly": 3}},
	}
	return data, dims
}
