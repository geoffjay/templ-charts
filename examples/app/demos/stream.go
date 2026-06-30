package demos

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/stream"
)

// StreamDemo is one stream tile on the /stream page.
type StreamDemo struct {
	ID          string
	Title       string
	Description string
	Props       stream.StreamProps
}

// StreamDemos returns the stream demos for the /stream page.
func StreamDemos() []StreamDemo {
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
	return []StreamDemo{
		{
			ID:          "stream-wiggle",
			Title:       "Wiggle offset (default)",
			Description: "Four series stacked with the wiggle offset and a catmullRom curve, y-grid + bottom/left axes.",
			Props: stream.StreamProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin: core.Margin{Top: 30, Right: 30, Bottom: 40, Left: 50},
				Data:   data, Keys: keys, Interactive: true,
			},
		},
		{
			ID:          "stream-expand",
			Title:       "Expand offset + borders",
			Description: "offsetType=expand normalizes each stack to fill the height; borderWidth=1 outlines the layers.",
			Props: stream.StreamProps{
				Width: commonChartWidth, Height: commonChartHeight,
				Margin:      core.Margin{Top: 30, Right: 30, Bottom: 40, Left: 50},
				Data:        data,
				Keys:        keys,
				OffsetType:  core.StackOffsetExpand,
				BorderWidth: 1,
			},
		},
	}
}
