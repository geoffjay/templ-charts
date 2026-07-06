// Package force is a port of the subset of d3-force that templ-charts needs to
// lay out the network and swarmplot charts server-side.
//
// It implements d3's velocity-Verlet integrator and the five forces those two
// charts use — ForceLink, ForceManyBody, ForceCenter, ForceX/ForceY and
// ForceCollide — driven the same way nivo drives them on the server: build the
// simulation, Stop() it, then Tick(n) a fixed number of iterations (nivo's
// default is 120). There is no wall-clock timer and no alpha-based auto-stop,
// so a run is a pure function of its inputs.
//
// Determinism. Node positions are seeded with d3's phyllotaxis spiral (not
// Math.random), and the microscopic "jiggle" that d3 applies to coincident
// nodes is driven by d3's built-in linear congruential generator (constants
// a=1664525, c=1013904223, m=2^32), matching internal/d3/hierarchy's LCG. So
// every run of the same simulation produces byte-identical output.
//
// Relationship to d3. ForceLink, ForceCenter, ForceX/ForceY and ForceCollide
// reproduce d3's math exactly. ForceManyBody uses an exact all-pairs n-body
// sum (equivalent to a Barnes–Hut θ of 0) rather than d3's quadtree
// approximation; for chart-scale node counts this is both cheaper to reason
// about and more accurate, at the cost of not being bit-identical to d3's
// approximated charge.
package force

import "math"

// Node is a simulation body. X/Y are its position and Vx/Vy its velocity;
// Index is assigned by the simulation. A node with Fx/Fy set is pinned to that
// coordinate (its velocity is zeroed each tick). Data carries the caller's
// datum so force accessors can read it. Build nodes with NewNode so unset
// positions are seeded via phyllotaxis, matching d3.
type Node struct {
	Index  int
	X, Y   float64
	Vx, Vy float64
	Fx, Fy *float64
	Data   any
}

// Link connects two nodes. Data carries the caller's datum for the
// distance/strength accessors. Index is assigned by ForceLink.
type Link struct {
	Source *Node
	Target *Node
	Index  int
	Data   any
}

// NewNode returns a node whose position and velocity are marked "unset" (NaN)
// so the simulation seeds them deterministically (phyllotaxis for position, 0
// for velocity), mirroring how a fresh object behaves in d3's JS.
func NewNode(data any) *Node {
	nan := math.NaN()
	return &Node{X: nan, Y: nan, Vx: nan, Vy: nan, Data: data}
}

// NewLink connects source to target, carrying data for force accessors.
func NewLink(source, target *Node, data any) *Link {
	return &Link{Source: source, Target: target, Data: data}
}

// Force is a physics force attached to a simulation. Ports mirror d3's
// force.initialize(nodes)/force(alpha) pair.
type Force interface {
	initialize(nodes []*Node, random *lcg)
	apply(alpha float64)
}

type namedForce struct {
	name  string
	force Force
}

// phyllotaxis seeding constants, matching d3-force.
const initialRadius = 10.0

var initialAngle = math.Pi * (3 - math.Sqrt(5))

// Simulation is a velocity-Verlet force simulation. Unlike d3 it has no timer;
// call Tick to advance it a fixed number of iterations.
type Simulation struct {
	nodes         []*Node
	alpha         float64
	alphaMin      float64
	alphaDecay    float64
	alphaTarget   float64
	velocityDecay float64
	forces        []namedForce
	random        *lcg
}

// NewSimulation creates a simulation over nodes, seeding node indices and any
// unset positions via phyllotaxis. Defaults match d3: alpha 1, alphaMin 0.001,
// alphaDecay 1-alphaMin^(1/300), alphaTarget 0, velocityDecay 0.6.
func NewSimulation(nodes []*Node) *Simulation {
	s := &Simulation{
		nodes:         nodes,
		alpha:         1,
		alphaMin:      0.001,
		alphaTarget:   0,
		velocityDecay: 0.6,
		random:        newLCG(),
	}
	s.alphaDecay = 1 - math.Pow(s.alphaMin, 1.0/300.0)
	s.initializeNodes()
	return s
}

func (s *Simulation) initializeNodes() {
	for i, node := range s.nodes {
		node.Index = i
		if node.Fx != nil {
			node.X = *node.Fx
		}
		if node.Fy != nil {
			node.Y = *node.Fy
		}
		if math.IsNaN(node.X) || math.IsNaN(node.Y) {
			radius := initialRadius * math.Sqrt(0.5+float64(i))
			angle := float64(i) * initialAngle
			node.X = radius * math.Cos(angle)
			node.Y = radius * math.Sin(angle)
		}
		if math.IsNaN(node.Vx) || math.IsNaN(node.Vy) {
			node.Vx, node.Vy = 0, 0
		}
	}
}

// Force attaches (or, with a nil force, removes) a force under name and
// initializes it against the current nodes, preserving insertion order.
// Returns the simulation for chaining.
func (s *Simulation) Force(name string, f Force) *Simulation {
	for i := range s.forces {
		if s.forces[i].name == name {
			if f == nil {
				s.forces = append(s.forces[:i], s.forces[i+1:]...)
			} else {
				f.initialize(s.nodes, s.random)
				s.forces[i].force = f
			}
			return s
		}
	}
	if f != nil {
		f.initialize(s.nodes, s.random)
		s.forces = append(s.forces, namedForce{name: name, force: f})
	}
	return s
}

// Tick advances the simulation by iterations steps (default 1), applying each
// force then integrating positions — d3's simulation.tick. Returns the
// simulation for chaining.
func (s *Simulation) Tick(iterations int) *Simulation {
	if iterations <= 0 {
		iterations = 1
	}
	for k := 0; k < iterations; k++ {
		s.alpha += (s.alphaTarget - s.alpha) * s.alphaDecay
		for _, nf := range s.forces {
			nf.force.apply(s.alpha)
		}
		for _, node := range s.nodes {
			if node.Fx == nil {
				node.Vx *= s.velocityDecay
				node.X += node.Vx
			} else {
				node.X = *node.Fx
				node.Vx = 0
			}
			if node.Fy == nil {
				node.Vy *= s.velocityDecay
				node.Y += node.Vy
			} else {
				node.Y = *node.Fy
				node.Vy = 0
			}
		}
	}
	return s
}

// Stop is a no-op kept for API parity with d3 (there is no internal timer to
// halt server-side). Returns the simulation for chaining.
func (s *Simulation) Stop() *Simulation { return s }

// Nodes returns the simulation's nodes (the same slice passed to
// NewSimulation, mutated in place).
func (s *Simulation) Nodes() []*Node { return s.nodes }

// jiggle returns a tiny random displacement, matching d3's jiggle: used when
// two bodies occupy the same coordinate so the force has a direction.
func jiggle(random *lcg) float64 { return (random.next() - 0.5) * 1e-6 }
