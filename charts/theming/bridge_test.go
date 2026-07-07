package theming

import "testing"

func TestConvertStyleAttribute(t *testing.T) {
	cases := []struct {
		engine Engine
		attr   string
		value  string
		want   string
	}{
		{EngineSVG, "textAlign", "start", "start"},
		{EngineSVG, "textAlign", "center", "middle"},
		{EngineSVG, "textAlign", "end", "end"},
		{EngineSVG, "textBaseline", "top", "text-before-edge"},
		{EngineSVG, "textBaseline", "center", "middle"},
		{EngineSVG, "textBaseline", "bottom", "text-after-edge"},
		{EngineCSS, "textAlign", "start", "left"},
		{EngineCSS, "textAlign", "end", "right"},
		{EngineCSS, "textBaseline", "center", "middle"},
		{EngineCanvas, "textAlign", "center", "center"},
		{EngineCanvas, "textBaseline", "top", "top"},
		// Unknown engine passes the value through.
		{Engine("webgl"), "textAlign", "center", "center"},
		// Unknown attribute passes through.
		{EngineSVG, "fontWeight", "bold", "bold"},
		// Unknown value passes through.
		{EngineSVG, "textAlign", "justify", "justify"},
	}
	for _, c := range cases {
		if got := ConvertStyleAttribute(c.engine, c.attr, c.value); got != c.want {
			t.Errorf("ConvertStyleAttribute(%s, %s, %s) = %q, want %q", c.engine, c.attr, c.value, got, c.want)
		}
	}
}

func TestSanitizeSvgTextStyleFull(t *testing.T) {
	fs := 12
	style := TextStyle{
		FontFamily:     "sans-serif",
		FontSize:       fs,
		Fill:           "#333333",
		OutlineWidth:   2,
		OutlineColor:   "#ffffff",
		OutlineOpacity: 1,
		Extra:          map[string]any{"fontWeight": "bold"},
	}
	got := SanitizeSvgTextStyle(style)
	if got["fontFamily"] != "sans-serif" || got["fontSize"] != fs || got["fill"] != "#333333" {
		t.Errorf("sanitized svg style = %v", got)
	}
	if got["fontWeight"] != "bold" {
		t.Errorf("extra props should be preserved: %v", got)
	}
	for _, k := range []string{"outlineWidth", "outlineColor", "outlineOpacity"} {
		if _, ok := got[k]; ok {
			t.Errorf("outline field %q should be stripped", k)
		}
	}
	// Zero-value style produces an empty map.
	if got := SanitizeSvgTextStyle(TextStyle{}); len(got) != 0 {
		t.Errorf("empty style = %v, want empty map", got)
	}
}

func TestSanitizeHtmlTextStyle(t *testing.T) {
	style := TextStyle{
		FontFamily: "serif",
		FontSize:   "11px",
		Fill:       "#222222",
		Extra:      map[string]any{"fontStyle": "italic"},
	}
	got := SanitizeHtmlTextStyle(style)
	if got["fontFamily"] != "serif" || got["fontSize"] != "11px" {
		t.Errorf("sanitized html style = %v", got)
	}
	// fill is rewritten as color for HTML.
	if got["color"] != "#222222" {
		t.Errorf("color = %v, want #222222", got["color"])
	}
	if _, ok := got["fill"]; ok {
		t.Error("fill key should not be present in html style")
	}
	if got["fontStyle"] != "italic" {
		t.Errorf("extra props should be preserved: %v", got)
	}
	if got := SanitizeHtmlTextStyle(TextStyle{}); len(got) != 0 {
		t.Errorf("empty style = %v, want empty map", got)
	}
}
