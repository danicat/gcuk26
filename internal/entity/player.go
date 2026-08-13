package entity

import (
	"image/color"
	"london-eco-rider/assets"
	"london-eco-rider/internal/physics"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayerMode int

const (
	PlayerModeOnFoot PlayerMode = iota
	PlayerModeOnBike
)

type BikeColor int

const (
	BikeColorRed BikeColor = iota
	BikeColorBlue
	BikeColorGold
)

// Player represents Leo walking on foot or riding his bicycle.
type Player struct {
	Mode            PlayerMode
	BikeColorChoice BikeColor
	Physics        *physics.BikePhysics
	IsTurbo        bool
	TurboTimer     float64
	Invincible     bool
	InvincTimer    float64
	FlashTimer     float64
	Particles      *ParticlePool
	FacingRight    bool
	AnimTimer      float64
	StepCycle      float64
}

// NewPlayer creates Leo on foot initially.
func NewPlayer(startX, startY float64) *Player {
	p := &Player{
		Mode:            PlayerModeOnFoot,
		BikeColorChoice: BikeColorRed,
		Physics:         physics.NewBikePhysics(startX, startY),
		Particles:       NewParticlePool(100),
		FacingRight:     true,
	}
	p.Physics.Width = 32.0
	p.Physics.Height = 38.0
	p.Physics.MaxSpeed = 2.8
	p.Physics.JumpForce = -7.5
	return p
}

// SwitchToBike unlocks the bicycle with chosen color frame.
func (p *Player) SwitchToBike(chosenColor BikeColor) {
	p.Mode = PlayerModeOnBike
	p.BikeColorChoice = chosenColor
	p.Physics.Width = 48.0
	p.Physics.Height = 36.0
	p.Physics.MaxSpeed = 5.5
	p.Physics.TurboMaxSpeed = 9.0
	p.Physics.JumpForce = -9.8
}

// Update handles input, physics deltas, and motion timing.
func (p *Player) Update(pedalForward, pedalBack, isJumping, triggerBoost bool) {
	if p.Mode == PlayerModeOnBike && triggerBoost && p.TurboTimer > 0 {
		p.IsTurbo = true
	} else if p.TurboTimer <= 0 {
		p.IsTurbo = false
	}

	if p.IsTurbo {
		p.TurboTimer -= 0.016
		p.Particles.Emit(
			p.Physics.X,
			p.Physics.Y+p.Physics.Height*0.7,
			-p.Physics.VX*0.5,
			(randFloat()*2.0-1.0)*1.0,
			3.0,
			0.3,
			color.RGBA{255, 200, 50, 255},
		)
	}

	if p.Invincible {
		p.InvincTimer -= 0.016
		p.FlashTimer += 0.016
		if p.InvincTimer <= 0 {
			p.Invincible = false
		}
	}

	p.Physics.Update(pedalForward, pedalBack, isJumping, p.IsTurbo)

	if pedalForward {
		p.FacingRight = true
	} else if pedalBack {
		p.FacingRight = false
	}

	if p.Physics.Grounded && math.Abs(p.Physics.VX) > 0.1 {
		p.AnimTimer += 0.016 * math.Abs(p.Physics.VX) * 8.0
		p.StepCycle = math.Sin(p.AnimTimer)
	} else {
		p.StepCycle = 0
		p.AnimTimer = 0
	}

	p.Particles.Update()
}

// TriggerTurbo activates turbo boost.
func (p *Player) TriggerTurbo(seconds float64) {
	if p.Mode == PlayerModeOnBike {
		p.TurboTimer = seconds
		p.IsTurbo = true
	}
}

// TriggerInvincibility makes player invincible temporarily.
func (p *Player) TriggerInvincibility(seconds float64) {
	p.Invincible = true
	p.InvincTimer = seconds
}

// Draw renders consistent character sprite with realistic cycling pedaling leg animation.
func (p *Player) Draw(screen *ebiten.Image, cameraX float64) {
	p.Particles.Draw(screen, cameraX)

	if p.Invincible && int(p.FlashTimer*20)%2 == 0 {
		return
	}

	screenX := p.Physics.X - cameraX
	screenY := p.Physics.Y

	// Ground Shadow under feet
	vector.DrawFilledCircle(screen, float32(screenX+p.Physics.Width*0.5), float32(screenY+p.Physics.Height-2), float32(p.Physics.Width*0.35), color.RGBA{20, 20, 30, 120}, false)

	bobY := 0.0
	tiltAngle := p.Physics.Angle

	if p.Physics.Grounded && math.Abs(p.Physics.VX) > 0.1 {
		bobY = math.Abs(p.StepCycle) * -2.0
		tiltAngle += p.StepCycle * 0.04
	}

	var sprite *ebiten.Image
	if p.Mode == PlayerModeOnFoot {
		sprite = assets.BoyWalkImage
	} else {
		sprite = assets.BoyBikeImage
	}

	if sprite != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := sprite.Bounds()

		// Color Tint for Bike Choice
		if p.Mode == PlayerModeOnBike {
			switch p.BikeColorChoice {
			case BikeColorRed:
				op.ColorScale.Scale(1.2, 0.8, 0.8, 1.0)
			case BikeColorBlue:
				op.ColorScale.Scale(0.8, 0.9, 1.3, 1.0)
			case BikeColorGold:
				op.ColorScale.Scale(1.3, 1.1, 0.6, 1.0)
			}
		}

		targetH := p.Physics.Height
		scale := targetH / float64(bounds.Dy())
		op.GeoM.Scale(scale, scale)

		drawnW := float64(bounds.Dx()) * scale
		drawnH := float64(bounds.Dy()) * scale

		op.GeoM.Translate(-drawnW/2.0, -drawnH/2.0)

		if !p.FacingRight {
			op.GeoM.Scale(-1, 1)
		}

		op.GeoM.Rotate(tiltAngle)
		op.GeoM.Translate(screenX+p.Physics.Width/2.0, screenY+p.Physics.Height/2.0+bobY)

		screen.DrawImage(sprite, op)

		// Overlay Realistic Cycling Leg Pedaling Animation when on Bicycle
		if p.Mode == PlayerModeOnBike && math.Abs(p.Physics.VX) > 0.1 {
			drawCyclingLegOverlay(screen, float32(screenX), float32(screenY+bobY), p.AnimTimer, p.FacingRight, p.BikeColorChoice)
		}
	} else {
		if p.Mode == PlayerModeOnFoot {
			drawProceduralBoy(screen, float32(screenX), float32(screenY+bobY), p.StepCycle)
		} else {
			drawProceduralBike(screen, float32(screenX), float32(screenY+bobY), p.Physics.Angle, p.IsTurbo, p.AnimTimer, p.BikeColorChoice)
		}
	}
}

// drawCyclingLegOverlay renders realistic circular 360-degree pedaling leg motion
func drawCyclingLegOverlay(screen *ebiten.Image, x, y float32, animTimer float64, facingRight bool, bikeCol BikeColor) {
	pedalAngle := animTimer * 2.5

	// Crank Center
	crankX := x + 24.0
	crankY := y + 26.0

	// Hip joint on seat
	hipX := x + 20.0
	hipY := y + 16.0

	// Pedal 1 (Right Leg)
	p1X := crankX + float32(math.Cos(pedalAngle)*5.0)
	p1Y := crankY + float32(math.Sin(pedalAngle)*5.0)

	// Knee 1 bending forward and up
	k1X := hipX + (p1X-hipX)*0.5 + 4.0
	k1Y := hipY + (p1Y-hipY)*0.5 - 3.0

	// Pedal 2 (Left Leg, 180 degrees opposite)
	p2X := crankX - float32(math.Cos(pedalAngle)*5.0)
	p2Y := crankY - float32(math.Sin(pedalAngle)*5.0)

	// Knee 2
	k2X := hipX + (p2X-hipX)*0.5 + 4.0
	k2Y := hipY + (p2Y-hipY)*0.5 - 3.0

	legCol1 := color.RGBA{40, 40, 50, 255}
	legCol2 := color.RGBA{60, 60, 80, 255}
	shoeCol := color.RGBA{220, 50, 50, 255}

	// Draw Back Leg
	vector.StrokeLine(screen, hipX, hipY, k2X, k2Y, 2.5, legCol2, false)
	vector.StrokeLine(screen, k2X, k2Y, p2X, p2Y, 2.5, legCol2, false)
	vector.DrawFilledCircle(screen, p2X, p2Y, 2.5, shoeCol, false)

	// Draw Front Leg
	vector.StrokeLine(screen, hipX, hipY, k1X, k1Y, 3.0, legCol1, false)
	vector.StrokeLine(screen, k1X, k1Y, p1X, p1Y, 3.0, legCol1, false)
	vector.DrawFilledCircle(screen, p1X, p1Y, 3.0, shoeCol, false)
}

func drawProceduralBoy(screen *ebiten.Image, x, y float32, stepCycle float64) {
	boyCol := color.RGBA{240, 180, 50, 255}
	shirtCol := color.RGBA{50, 120, 220, 255}

	vector.DrawFilledCircle(screen, x+16, y+8, 7, boyCol, false)
	vector.DrawFilledRect(screen, x+10, y+15, 12, 12, shirtCol, false)

	legOffset := float32(stepCycle * 4.0)
	vector.StrokeLine(screen, x+12, y+27, x+12-legOffset, y+36, 2, color.RGBA{30, 30, 30, 255}, false)
	vector.StrokeLine(screen, x+20, y+27, x+20+legOffset, y+36, 2, color.RGBA{30, 30, 30, 255}, false)
}

func drawProceduralBike(screen *ebiten.Image, x, y float32, angle float64, isTurbo bool, animTimer float64, bikeColChoice BikeColor) {
	wheelCol := color.RGBA{40, 40, 40, 255}
	frameCol := color.RGBA{0, 150, 220, 255}
	switch bikeColChoice {
	case BikeColorRed:
		frameCol = color.RGBA{230, 40, 40, 255}
	case BikeColorGold:
		frameCol = color.RGBA{255, 215, 0, 255}
	}

	boyCol := color.RGBA{240, 180, 50, 255}

	vector.DrawFilledCircle(screen, x+10, y+28, 8, wheelCol, false)
	vector.DrawFilledCircle(screen, x+38, y+28, 8, wheelCol, false)

	spokeAngle := float32(animTimer * 2.5)
	vector.StrokeLine(screen, x+10, y+28, x+10+float32(math.Cos(float64(spokeAngle)))*6, y+28+float32(math.Sin(float64(spokeAngle)))*6, 1, color.RGBA{180, 180, 180, 255}, false)
	vector.StrokeLine(screen, x+38, y+28, x+38+float32(math.Cos(float64(spokeAngle)))*6, y+28+float32(math.Sin(float64(spokeAngle)))*6, 1, color.RGBA{180, 180, 180, 255}, false)

	vector.StrokeLine(screen, x+10, y+28, x+24, y+28, 2, frameCol, false)
	vector.StrokeLine(screen, x+24, y+28, x+34, y+16, 2, frameCol, false)
	vector.StrokeLine(screen, x+10, y+28, x+20, y+16, 2, frameCol, false)
	vector.StrokeLine(screen, x+20, y+16, x+34, y+16, 2, frameCol, false)
	vector.StrokeLine(screen, x+34, y+16, x+36, y+8, 2, color.RGBA{180, 180, 180, 255}, false)

	vector.DrawFilledCircle(screen, x+22, y+6, 6, boyCol, false)
	vector.StrokeLine(screen, x+22, y+12, x+20, y+22, 3, color.RGBA{50, 120, 200, 255}, false)

	drawCyclingLegOverlay(screen, x, y, animTimer, true, bikeColChoice)
}

func randFloat() float64 {
	return 0.5
}
