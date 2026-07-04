package samples

import "github.com/geoffjay/templ-charts/charts/stream"

// Stream returns seven time steps of four stacked series (Raoul, Josiane,
// Marcel, René). The returned data and keys populate a stream.StreamProps'
// Data and Keys fields.
func Stream() ([]stream.StreamDatum, []string) {
	keys := []string{"Raoul", "Josiane", "Marcel", "René"}
	data := []stream.StreamDatum{
		{"Raoul": 10, "Josiane": 20, "Marcel": 30, "René": 14},
		{"Raoul": 15, "Josiane": 18, "Marcel": 25, "René": 20},
		{"Raoul": 12, "Josiane": 22, "Marcel": 28, "René": 11},
		{"Raoul": 18, "Josiane": 16, "Marcel": 35, "René": 24},
		{"Raoul": 14, "Josiane": 24, "Marcel": 20, "René": 18},
		{"Raoul": 20, "Josiane": 12, "Marcel": 30, "René": 16},
		{"Raoul": 16, "Josiane": 28, "Marcel": 22, "René": 21},
	}
	return data, keys
}
