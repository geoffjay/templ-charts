package samples

// Radar returns ready-to-render radar data: three wine varieties scored across
// five taste axes. The returned keys slice is the series (varieties); the
// index-by property for the axes is "taste".
func Radar() ([]map[string]any, []string) {
	data := []map[string]any{
		{"taste": "fruity", "chardonnay": 93.0, "carmenere": 61.0, "syrah": 114.0},
		{"taste": "bitter", "chardonnay": 91.0, "carmenere": 37.0, "syrah": 72.0},
		{"taste": "heavy", "chardonnay": 56.0, "carmenere": 95.0, "syrah": 99.0},
		{"taste": "strong", "chardonnay": 64.0, "carmenere": 90.0, "syrah": 30.0},
		{"taste": "sunny", "chardonnay": 119.0, "carmenere": 94.0, "syrah": 103.0},
	}
	keys := []string{"chardonnay", "carmenere", "syrah"}
	return data, keys
}
