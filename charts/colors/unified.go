package colors

// This file adds one ergonomic, additive entry point for setting colors that
// resolves to whichever of the five existing config shapes a chart's field
// expects — OrdinalColorScaleConfig (categorical `Colors`), InheritedColorConfig
// (`BorderColor`/`LabelTextColor`/…), SequentialColorScaleConfig,
// DivergingColorScaleConfig, and plain color strings. It renames/removes
// nothing: every existing type, field, and constructor is untouched, so no
// consumer breaks and no golden moves. Callers who prefer the explicit configs
// keep using them; this is a convenience layer on top.

// colorSource discriminates what a ColorSetting was built from.
type colorSource int

const (
	srcStatic     colorSource = iota // a single CSS color string
	srcScheme                        // a named palette/scheme id
	srcColors                        // an explicit color list
	srcFunc                          // a datum→color function
	srcDatum                         // a datum property path (ordinal) / from-path (inherited)
	srcTheme                         // a theme property path (inherited)
	srcOrdinal                       // a pre-built OrdinalColorScaleConfig
	srcInherited                     // a pre-built InheritedColorConfig
	srcSequential                    // a pre-built SequentialColorScaleConfig
	srcDiverging                     // a pre-built DivergingColorScaleConfig
)

// ColorSetting is a unified, additive color intent. Build it with Set (or the
// typed Set* constructors), optionally refine it with WithModifiers/InSpace,
// then resolve it to the config a chart field wants via Ordinal/Inherited/
// Sequential/Diverging (or Static for a raw color string).
type ColorSetting struct {
	source colorSource

	static string
	scheme string
	colors []string
	fn     func(any) string
	path   string // datum or theme path
	mods   []ColorModifier
	space  Space

	ordinal    OrdinalColorScaleConfig
	inherited  InheritedColorConfig
	sequential SequentialColorScaleConfig
	diverging  DivergingColorScaleConfig
}

// Set is the ergonomic entry point. It infers the intent from v's dynamic type:
//   - nil                          → empty static color
//   - PaletteID / Palette          → that palette's scheme
//   - string                       → a scheme id if recognized, else a static color
//   - []string                     → an explicit color list
//   - func(any) string             → a datum→color function
//   - a pre-built *ColorScaleConfig → wrapped verbatim (Ordinal/Inherited/Sequential/Diverging)
//
// Anything unrecognized resolves to an empty static color.
func Set(v any) ColorSetting {
	switch x := v.(type) {
	case nil:
		return ColorSetting{source: srcStatic, static: ""}
	case ColorSetting:
		return x
	case PaletteID:
		return ColorSetting{source: srcScheme, scheme: string(x)}
	case Palette:
		return ColorSetting{source: srcScheme, scheme: string(x.ID)}
	case string:
		if IsCategoricalColorScheme(x) || IsDivergingColorScheme(x) || IsSequentialColorScheme(x) {
			return ColorSetting{source: srcScheme, scheme: x}
		}
		return ColorSetting{source: srcStatic, static: x}
	case []string:
		return ColorSetting{source: srcColors, colors: x}
	case func(any) string:
		return ColorSetting{source: srcFunc, fn: x}
	case OrdinalColorScaleConfig:
		return ColorSetting{source: srcOrdinal, ordinal: x}
	case InheritedColorConfig:
		return ColorSetting{source: srcInherited, inherited: x}
	case SequentialColorScaleConfig:
		return ColorSetting{source: srcSequential, sequential: x}
	case DivergingColorScaleConfig:
		return ColorSetting{source: srcDiverging, diverging: x}
	}
	return ColorSetting{source: srcStatic, static: ""}
}

// SetScheme selects a named palette/scheme (categorical, sequential, or
// diverging depending on the resolver used).
func SetScheme(id PaletteID) ColorSetting {
	return ColorSetting{source: srcScheme, scheme: string(id)}
}

// SetColorList uses an explicit ordered color list.
func SetColorList(cols ...string) ColorSetting {
	return ColorSetting{source: srcColors, colors: cols}
}

// SetStatic uses a single static CSS color.
func SetStatic(color string) ColorSetting {
	return ColorSetting{source: srcStatic, static: color}
}

// SetFunc uses a datum→color function.
func SetFunc(fn func(any) string) ColorSetting {
	return ColorSetting{source: srcFunc, fn: fn}
}

// SetFromDatum draws the color from a datum property path (an ordinal `datum`
// config, or the `from` path of an inherited config when resolved via
// Inherited).
func SetFromDatum(path string) ColorSetting {
	return ColorSetting{source: srcDatum, path: path}
}

// SetFromTheme draws the color from a theme property path (inherited).
func SetFromTheme(path string) ColorSetting {
	return ColorSetting{source: srcTheme, path: path}
}

// WithModifiers attaches brighter/darker/opacity modifiers (used when resolved
// to an InheritedColorConfig).
func (c ColorSetting) WithModifiers(mods ...ColorModifier) ColorSetting {
	c.mods = mods
	return c
}

// InSpace sets the interpolation/modifier color space, applied to the
// Sequential/Diverging scale Space and the Inherited modifier space on resolve.
func (c ColorSetting) InSpace(space Space) ColorSetting {
	c.space = space
	return c
}

// Ordinal resolves to an OrdinalColorScaleConfig (the categorical `Colors`
// shape). Pre-built ordinal configs pass through; scheme-bearing continuous
// configs reuse their scheme id.
func (c ColorSetting) Ordinal() OrdinalColorScaleConfig {
	switch c.source {
	case srcScheme:
		return OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: c.scheme}
	case srcColors:
		return OrdinalColorScaleConfig{Type: OrdinalTypeColors, Colors: c.colors}
	case srcStatic:
		return OrdinalColorScaleConfig{Type: OrdinalTypeStatic, Static: c.static}
	case srcFunc:
		return OrdinalColorScaleConfig{Type: OrdinalTypeFunc, Func: c.fn}
	case srcDatum:
		return OrdinalColorScaleConfig{Type: OrdinalTypeDatum, DatumPath: c.path}
	case srcOrdinal:
		return c.ordinal
	case srcSequential:
		if c.sequential.Scheme != "" {
			return OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: c.sequential.Scheme}
		}
	case srcDiverging:
		if c.diverging.Scheme != "" {
			return OrdinalColorScaleConfig{Type: OrdinalTypeScheme, Scheme: c.diverging.Scheme}
		}
	case srcInherited:
		if c.inherited.Type == InheritedColorTypeStatic {
			return OrdinalColorScaleConfig{Type: OrdinalTypeStatic, Static: c.inherited.Static}
		}
	}
	return OrdinalColorScaleConfig{Type: OrdinalTypeStatic, Static: c.static}
}

// Inherited resolves to an InheritedColorConfig (the `BorderColor`/
// `LabelTextColor` shape). Static/func/theme/from intents map directly; scheme
// and color-list intents (which have no single inherited meaning) fall back to
// a static color.
func (c ColorSetting) Inherited() InheritedColorConfig {
	switch c.source {
	case srcStatic:
		return NewStaticColor(c.static)
	case srcFunc:
		return NewFuncColor(c.fn)
	case srcTheme:
		return NewThemeColor(c.path)
	case srcDatum:
		return NewFromContextColorInSpace(c.path, c.mods, c.space)
	case srcInherited:
		cfg := c.inherited
		if len(c.mods) > 0 {
			cfg.Modifiers = c.mods
		}
		if c.space != SpaceRGB {
			cfg.ModifierSpace = c.space
		}
		return cfg
	case srcColors:
		if len(c.colors) > 0 {
			return NewStaticColor(c.colors[0])
		}
	}
	return NewStaticColor(c.static)
}

// Sequential resolves to a SequentialColorScaleConfig. Scheme and two-color
// intents map directly; the InSpace selector flows into the scale's Space.
func (c ColorSetting) Sequential() SequentialColorScaleConfig {
	switch c.source {
	case srcSequential:
		cfg := c.sequential
		cfg.Type = "sequential"
		if c.space != SpaceRGB {
			cfg.Space = c.space
		}
		return cfg
	case srcScheme:
		return SequentialColorScaleConfig{Type: "sequential", Scheme: c.scheme, Space: c.space}
	case srcColors:
		var cols [2]string
		if len(c.colors) > 0 {
			cols[0] = c.colors[0]
			cols[1] = c.colors[len(c.colors)-1]
		}
		return SequentialColorScaleConfig{Type: "sequential", Colors: cols, Space: c.space}
	}
	return SequentialColorScaleConfig{Type: "sequential", Space: c.space}
}

// Diverging resolves to a DivergingColorScaleConfig. Scheme and three-color
// intents map directly; the InSpace selector flows into the scale's Space.
func (c ColorSetting) Diverging() DivergingColorScaleConfig {
	switch c.source {
	case srcDiverging:
		cfg := c.diverging
		cfg.Type = "diverging"
		if c.space != SpaceRGB {
			cfg.Space = c.space
		}
		return cfg
	case srcScheme:
		return DivergingColorScaleConfig{Type: "diverging", Scheme: c.scheme, Space: c.space}
	case srcColors:
		var cols [3]string
		for i := 0; i < len(c.colors) && i < 3; i++ {
			cols[i] = c.colors[i]
		}
		return DivergingColorScaleConfig{Type: "diverging", Colors: cols, Space: c.space}
	}
	return DivergingColorScaleConfig{Type: "diverging", Space: c.space}
}

// Static resolves to a single color string (for raw string color fields like
// NodeColor/LinkColor). Non-static intents fall back to the first list color or
// an empty string.
func (c ColorSetting) Static() string {
	switch c.source {
	case srcStatic:
		return c.static
	case srcColors:
		if len(c.colors) > 0 {
			return c.colors[0]
		}
	case srcInherited:
		if c.inherited.Type == InheritedColorTypeStatic {
			return c.inherited.Static
		}
	}
	return c.static
}
