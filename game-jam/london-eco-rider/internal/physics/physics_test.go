package physics

import "testing"

func TestBikePhysics_Acceleration(t *testing.T) {
	bike := NewBikePhysics(0, 0)
	initialVX := bike.VX

	// Pedal forward
	bike.Update(true, false, false, false)

	if bike.VX <= initialVX {
		t.Errorf("Expected bike velocity to increase when pedaling forward, got %f", bike.VX)
	}
}

func TestBikePhysics_Jump(t *testing.T) {
	bike := NewBikePhysics(0, 100)
	bike.GroundY = 100

	// Jump
	bike.Update(false, false, true, false)

	if bike.Grounded {
		t.Errorf("Expected bike to be airborne after jump")
	}
	if bike.VY >= 0 {
		t.Errorf("Expected negative upward jump velocity, got %f", bike.VY)
	}
}

func TestAABB_Intersects(t *testing.T) {
	boxA := AABB{X: 0, Y: 0, Width: 10, Height: 10}
	boxB := AABB{X: 5, Y: 5, Width: 10, Height: 10}
	boxC := AABB{X: 20, Y: 20, Width: 10, Height: 10}

	if !boxA.Intersects(boxB) {
		t.Errorf("Expected boxA and boxB to intersect")
	}
	if boxA.Intersects(boxC) {
		t.Errorf("Expected boxA and boxC not to intersect")
	}
}
