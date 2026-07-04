package samples

// Pie returns ready-to-render pie data: nine programming-language slices as
// {id, value} maps (the []any shape PieProps.Data accepts).
func Pie() []any {
	labels := []string{"Go", "Rust", "TypeScript", "Python", "Ruby", "Elixir", "C", "Zig", "Kotlin"}
	data := make([]any, len(labels))
	for i, l := range labels {
		data[i] = map[string]any{"id": l, "value": float64((i*17+3)%80 + 10)}
	}
	return data
}
