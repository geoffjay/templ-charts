// position.go — ports of d3-force's forceX and forceY: per-axis springs
// pulling each node toward a target coordinate. Strength defaults to 0.1.
package force

// XForce is d3's forceX.
type XForce struct {
	nodes      []*Node
	xFn        func(*Node) float64
	strengthFn func(*Node) float64
	xs         []float64
	strengths  []float64
}

// ForceX creates an x-positioning force with a per-node target accessor.
func ForceX(fn func(*Node) float64) *XForce {
	return &XForce{xFn: fn, strengthFn: func(*Node) float64 { return 0.1 }}
}

// Strength sets a constant per-tick pull toward the target x (d3 default 0.1).
func (f *XForce) Strength(s float64) *XForce {
	f.strengthFn = func(*Node) float64 { return s }
	return f
}

func (f *XForce) initialize(nodes []*Node, _ *lcg) {
	f.nodes = nodes
	f.xs = make([]float64, len(nodes))
	f.strengths = make([]float64, len(nodes))
	for i, node := range nodes {
		f.strengths[i] = f.strengthFn(node)
		f.xs[i] = f.xFn(node)
	}
}

func (f *XForce) apply(alpha float64) {
	for i, node := range f.nodes {
		node.Vx += (f.xs[i] - node.X) * f.strengths[i] * alpha
	}
}

// YForce is d3's forceY.
type YForce struct {
	nodes      []*Node
	yFn        func(*Node) float64
	strengthFn func(*Node) float64
	ys         []float64
	strengths  []float64
}

// ForceY creates a y-positioning force with a per-node target accessor.
func ForceY(fn func(*Node) float64) *YForce {
	return &YForce{yFn: fn, strengthFn: func(*Node) float64 { return 0.1 }}
}

// Strength sets a constant per-tick pull toward the target y (d3 default 0.1).
func (f *YForce) Strength(s float64) *YForce {
	f.strengthFn = func(*Node) float64 { return s }
	return f
}

func (f *YForce) initialize(nodes []*Node, _ *lcg) {
	f.nodes = nodes
	f.ys = make([]float64, len(nodes))
	f.strengths = make([]float64, len(nodes))
	for i, node := range nodes {
		f.strengths[i] = f.strengthFn(node)
		f.ys[i] = f.yFn(node)
	}
}

func (f *YForce) apply(alpha float64) {
	for i, node := range f.nodes {
		node.Vy += (f.ys[i] - node.Y) * f.strengths[i] * alpha
	}
}
