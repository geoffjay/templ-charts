package colors

// ColorSchemeId is the union of all color scheme identifiers nivo supports.
type ColorSchemeId string

// ColorSchemes merges categorical + diverging + sequential scheme maps.
// Categorical entries are []string; diverging/sequential are map[int][]string.
// Use CategoricalColorSchemes / DivergingColorSchemes / SequentialColorSchemes
// for typed access.
var ColorSchemes = map[string]any{}

func init() {
	for k, v := range CategoricalColorSchemes {
		ColorSchemes[k] = v
	}
	for k, v := range DivergingColorSchemes {
		ColorSchemes[k] = v
	}
	for k, v := range SequentialColorSchemes {
		ColorSchemes[k] = v
	}
}

// ColorSchemeIds is the full ordered list of scheme ids.
var ColorSchemeIds = func() []string {
	ids := make([]string, 0, len(CategoricalColorSchemes)+len(DivergingColorSchemes)+len(SequentialColorSchemes))
	ids = append(ids, categoricalColorSchemeIds...)
	ids = append(ids, divergingColorSchemeIds...)
	ids = append(ids, sequentialColorSchemeIds...)
	return ids
}()

// IsCategoricalColorScheme reports whether scheme is a categorical id.
func IsCategoricalColorScheme(scheme string) bool {
	if _, ok := CategoricalColorSchemes[scheme]; ok {
		return true
	}
	return false
}

// IsDivergingColorScheme reports whether scheme is a diverging id.
func IsDivergingColorScheme(scheme string) bool {
	if _, ok := DivergingColorSchemes[scheme]; ok {
		return true
	}
	return false
}

// IsSequentialColorScheme reports whether scheme is a sequential id.
func IsSequentialColorScheme(scheme string) bool {
	if _, ok := SequentialColorSchemes[scheme]; ok {
		return true
	}
	return false
}
