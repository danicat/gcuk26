package physics

import "math"

// AABB represents an Axis-Aligned Bounding Box.
type AABB struct {
	X, Y, Width, Height float64
}

// Intersects checks if two AABB bounding boxes overlap.
func (a AABB) Intersects(b AABB) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

// BikePhysics handles momentum, acceleration, friction, gravity, and jumping.
type BikePhysics struct {
	X, Y           float64
	VX, VY         float64
	Ax, Ay         float64
	Width, Height  float64
	Grounded       bool
	Angle          float64
	MaxSpeed       float64
	TurboMaxSpeed  float64
	Acceleration   float64
	Friction       float64
	Gravity        float64
	JumpForce      float64
	GroundY        float64
}

// NewBikePhysics initializes bike physics parameters.
func NewBikePhysics(startX, startY float64) *BikePhysics {
	return &BikePhysics{
		X:             startX,
		Y:             startY,
		Width:         32.0,
		Height:        38.0,
		Grounded:      true,
		MaxSpeed:      3.2,
		TurboMaxSpeed: 9.0,
		Acceleration:  0.25,
		Friction:      0.96,
		Gravity:       0.45,
		JumpForce:     -7.5,
		GroundY:       startY,
	}
}

// Update updates physics deltas with symmetrical left/right max speeds.
func (p *BikePhysics) Update(pedalForward, pedalBack, isJumping, isTurbo bool) {
	topSpeed := p.MaxSpeed
	if isTurbo {
		topSpeed = p.TurboMaxSpeed
	}

	// Symmetrical Left / Right Movement Speed
	if pedalForward {
		p.VX += p.Acceleration
		if p.VX > topSpeed {
			p.VX = topSpeed
		}
	} else if pedalBack {
		p.VX -= p.Acceleration
		if p.VX < -topSpeed { // Symmetrical leftward speed!
			p.VX = -topSpeed
		}
	} else {
		p.VX *= p.Friction
		if math.Abs(p.VX) < 0.01 {
			p.VX = 0
		}
	}

	// Jumping
	if isJumping && p.Grounded {
		p.VY = p.JumpForce
		p.Grounded = false
	}

	// Gravity
	if !p.Grounded {
		p.VY += p.Gravity
		p.Angle = p.VY * 0.025
	} else {
		p.Angle *= 0.8
	}

	// Position step
	p.X += p.VX
	p.Y += p.VY

	// Street Ground Floor Collision
	if p.Y >= p.GroundY {
		p.Y = p.GroundY
		p.VY = 0
		p.Grounded = true
	}
}

// GetAABB returns bounding box.
func (p *BikePhysics) GetAABB() AABB {
	return AABB{
		X:      p.X,
		Y:      p.Y,
		Width:  p.Width,
		Height: p.Height,
	}
}
