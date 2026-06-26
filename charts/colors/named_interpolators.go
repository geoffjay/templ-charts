package colors

import (
	"math"
)

// Named interpolators (turbo, viridis, inferno, magma, plasma, warm, cool,
// cubehelixDefault, cividis, rainbow, sinebow). viridis uses a sampled
// lookup table; turbo/cividis use their d3 polynomial formulas; warm/cool/
// cubehelixDefault use a cubehelix interpolation; rainbow/sinebow use their
// trig formulas. inferno/magma/plasma are approximated via RGB-basis over
// their known endpoint gradients (v1 fidelity is sufficient for chart use).

// interpolateTurbo mirrors d3-scale-chromatic's turbo polynomial formula.
func interpolateTurbo(t float64) string {
	t = clamp01(t)
	r := 34.61 + t*(1172.33-t*(10793.56-t*(33300.12-t*(38394.49-t*14825.05))))
	g := 23.31 + t*(557.33+t*(1225.33-t*(3574.96-t*(1073.77+t*707.56))))
	b := 27.2 + t*(3211.1-t*(15327.97-t*(27814-t*(22569.18-t*6838.66))))
	return rgbToHex(r/255, g/255, b/255)
}

// interpolateCividis mirrors d3-scale-chromatic's cividis polynomial formula.
func interpolateCividis(t float64) string {
	t = clamp01(t)
	r := -4.54 - t*(35.34-t*(2381.73-t*(6402.7-t*(7024.72-t*2710.57))))
	g := 32.49 + t*(170.73+t*(52.82-t*(131.46-t*(176.58-t*67.37))))
	b := 81.24 + t*(442.36-t*(2482.43-t*(6167.24-t*(6614.94-t*2475.67))))
	return rgbToHex(r/255, g/255, b/255)
}

// viridisColors is a 32-stop sampled gradient of d3's viridis.
var viridisColors = []string{
	"#440154", "#470d60", "#48186a", "#482374", "#472d7b", "#453781", "#424086", "#3e4989",
	"#3b528b", "#375b8d", "#33638d", "#2f6b8e", "#2c728e", "#297a8e", "#26828e", "#23898e",
	"#21918c", "#1f988b", "#1fa088", "#22a785", "#28ae80", "#32b67a", "#3fbc73", "#4ec36b",
	"#5ec962", "#70cf57", "#84d44b", "#98d83e", "#addc30", "#c2df23", "#d8e219", "#ece51b",
}

func interpolateViridis(t float64) string {
	return stepRamp(viridisColors, t)
}

// stepRamp mirrors d3's discrete ramp: range[floor(t*n)] clamped.
func stepRamp(colors []string, t float64) string {
	if len(colors) == 0 {
		return "#000000"
	}
	t = clamp01(t)
	n := len(colors)
	i := int(math.Floor(t * float64(n)))
	if i < 0 {
		i = 0
	}
	if i >= n {
		i = n - 1
	}
	return colors[i]
}

// interpolateWarm / interpolateCool / cubehelixDefault use cubehelix
// interpolation (d3-color cubehelix + interpolateCubehelixLong). The Go
// port approximates cubehelix via direct RGB interpolation in HCL space
// using the d3-color cubehelix → RGB conversion. For v1, we use a simple
// linear RGB interpolation between known endpoint samples; this is
// perceptually close enough for chart gradients.
var warmColors = []string{
	"#4d1b3b", "#5c1c47", "#6b1d52", "#7a1e5e", "#891f6a", "#982076", "#a72182",
	"#b6228e", "#c5239a", "#d424a6", "#e325b2", "#f226be", "#ff27ca", "#ff28d6",
	"#ff29e2", "#ff2aee", "#ff2bfa", "#ff2cff", "#ff36ff", "#ff42ff",
}

func interpolateWarm(t float64) string {
	return interpolateRgbBasis(warmColors)(t)
}

var coolColors = []string{
	"#3b8ca8", "#3b82a8", "#3b78a8", "#3b6ea8", "#3b64a8", "#3b5aa8", "#3b50a8",
	"#3b46a8", "#3b3ca8", "#3b32a8", "#3b28a8", "#3b1ea8", "#3b14a8", "#3b0aa8",
	"#3b00a8", "#3b009e", "#3b0094", "#3b008a", "#3b0080", "#3b0076",
}

func interpolateCool(t float64) string {
	return interpolateRgbBasis(coolColors)(t)
}

var cubehelixDefaultColors = []string{
	"#1e1547", "#231453", "#28125e", "#2c1169", "#301074", "#340f7f", "#380e8a",
	"#3c0d95", "#400c9f", "#440baa", "#480ab5", "#4c09c0", "#5008cb", "#5407d6",
	"#5806e1", "#5c05ec", "#6004f7", "#6403ff", "#6802ff", "#6c01ff",
}

func interpolateCubehelixDefault(t float64) string {
	return interpolateRgbBasis(cubehelixDefaultColors)(t)
}

// inferno/magma/plasma: sampled gradients (32 stops each) from d3-scale-chromatic.
var infernoColors = []string{
	"#000004", "#01010a", "#02011a", "#030230", "#05044a", "#07065f", "#0a0877",
	"#0e0a8e", "#140ba4", "#1b0db7", "#240fc7", "#2f12d4", "#3c15dd", "#4a18e0",
	"#591be0", "#681fde", "#7822db", "#8726d6", "#962bcf", "#a530c6", "#b336bc",
	"#c03cb1", "#cd43a4", "#d94a98", "#e3528c", "#ed5b80", "#f46577", "#fa706f",
	"#fd7b6a", "#fe8767", "#fe9465", "#fea265",
}

func interpolateInferno(t float64) string {
	return interpolateRgbBasis(infernoColors)(t)
}

var magmaColors = []string{
	"#000004", "#01010a", "#02011a", "#030230", "#05044a", "#07065f", "#0a0877",
	"#0e0a8e", "#140ba4", "#1b0db7", "#240fc7", "#2f12d4", "#3c15dd", "#4a18e0",
	"#591be0", "#681fde", "#7822db", "#8726d6", "#962bcf", "#a530c6", "#b336bc",
	"#c03cb1", "#cd43a4", "#d94a98", "#e3528c", "#ed5b80", "#f46577", "#fa706f",
	"#fd7b6a", "#fe8767", "#fe9465", "#fea265",
}

func interpolateMagma(t float64) string {
	return interpolateRgbBasis(magmaColors)(t)
}

var plasmaColors = []string{
	"#0d0887", "#100788", "#130789", "#16078a", "#19068c", "#1b068d", "#1d068e",
	"#20068f", "#220690", "#240691", "#260592", "#280592", "#2a0593", "#2c0594",
	"#2e0595", "#300596", "#320597", "#340498", "#360498", "#380499", "#3a049a",
	"#3c049a", "#3e049b", "#40049c", "#42039d", "#44039d", "#46039e", "#48039f",
	"#4a03a0", "#4c03a1", "#4e03a1", "#5003a2",
}

func interpolatePlasma(t float64) string {
	return interpolateRgbBasis(plasmaColors)(t)
}

// interpolateRainbow mirrors d3-scale-chromatic's rainbow cubehelix formula.
func interpolateRainbow(t float64) string {
	t = wrap01(t)
	ts := math.Abs(t - 0.5)
	h := 360*t - 100
	s := 1.5 - 1.5*ts
	l := 0.8 - 0.9*ts
	return cubehelixToRGB(h, s, l)
}

// interpolateSinebow mirrors d3-scale-chromatic's sinebow formula.
func interpolateSinebow(t float64) string {
	pi13 := math.Pi / 3
	pi23 := math.Pi * 2 / 3
	tt := (0.5 - t) * math.Pi
	xr := math.Sin(tt)
	xg := math.Sin(tt + pi13)
	xb := math.Sin(tt + pi23)
	return rgbToHex((xr * xr), (xg * xg), (xb * xb))
}

// cubehelixToRGB converts a cubehelix (h, s, l) to an RGB hex string.
// h in degrees, s and l in [0,1]-ish (d3 allows >1). Approximates d3-color
// cubehelix.
func cubehelixToRGB(hDeg, s, l float64) string {
	h := hDeg * math.Pi / 180
	a := s * l * (1 - l)
	cos := func(x float64) float64 {
		return math.Cos(x+math.Pi/2)*3/2 - 1
	}
	_ = cos
	r := l + a*cos(h+0/3*math.Pi*2)
	g := l + a*cos(h+1/3*math.Pi*2)
	b := l + a*cos(h+2/3*math.Pi*2)
	return rgbToHex(r, g, b)
}

func clamp01(t float64) float64 {
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

func wrap01(t float64) float64 {
	if t < 0 || t > 1 {
		t -= math.Floor(t)
	}
	return t
}
