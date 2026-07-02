// center.go — port of d3-force's forceCenter: translates every node each tick
// so the mean position sits at (x, y). It adjusts positions directly (not
// velocities) and ignores alpha, matching d3.
package force

// CenterForce is d3's forceCenter. Strength defaults to 1.
type CenterForce struct {
	nodes    []*Node
	x, y     float64
	strength float64
}

// ForceCenter creates a centering force pulling the mean position to (x, y).
func ForceCenter(x, y float64) *CenterForce {
	return &CenterForce{x: x, y: y, strength: 1}
}

// Strength sets how strongly each tick corrects toward the center (d3 default
// 1). Returns f for chaining.
func (f *CenterForce) Strength(s float64) *CenterForce {
	f.strength = s
	return f
}

func (f *CenterForce) initialize(nodes []*Node, _ *lcg) { f.nodes = nodes }

func (f *CenterForce) apply(float64) {
	n := len(f.nodes)
	if n == 0 {
		return
	}
	var sx, sy float64
	for _, node := range f.nodes {
		sx += node.X
		sy += node.Y
	}
	sx = (sx/float64(n) - f.x) * f.strength
	sy = (sy/float64(n) - f.y) * f.strength
	for _, node := range f.nodes {
		node.X -= sx
		node.Y -= sy
	}
}
