package entity

import (
	"image/color"
	"london-eco-rider/assets"
	"london-eco-rider/internal/physics"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ItemType int

const (
	ItemIceCream ItemType = iota
	ItemPlasticBottle
)

// Item represents a collectible in the level.
type Item struct {
	Type      ItemType
	X, Y      float64
	Width     float64
	Height    float64
	Collected bool
	BobTimer  float64
	BaseY     float64
}

// NewItem spawns an ice cream cone or plastic bottle.
func NewItem(itemType ItemType, x, y float64) *Item {
	return &Item{
		Type:   itemType,
		X:      x,
		Y:      y,
		BaseY:  y,
		Width:  18.0,
		Height: 28.0,
	}
}

// Update handles floating bob animation.
func (it *Item) Update() {
	if it.Collected {
		return
	}
	it.BobTimer += 0.06
	it.Y = it.BaseY + math.Sin(it.BobTimer)*3.5
}

// GetAABB returns bounding box for collision detection.
func (it *Item) GetAABB() physics.AABB {
	return physics.AABB{
		X:      it.X,
		Y:      it.Y,
		Width:  it.Width,
		Height: it.Height,
	}
}

// Draw renders item sprite preserving natural aspect ratio.
func (it *Item) Draw(screen *ebiten.Image, cameraX float64) {
	if it.Collected {
		return
	}
	screenX := it.X - cameraX
	screenY := it.Y

	var img *ebiten.Image
	if it.Type == ItemIceCream {
		img = assets.IceCreamImage
	} else {
		img = assets.PlasticBottleImg
	}

	if img != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := img.Bounds()

		// Preserve natural aspect ratio
		scale := 26.0 / float64(bounds.Dy())
		op.GeoM.Scale(scale, scale)

		// Center horizontally over entity bounding box
		drawnW := float64(bounds.Dx()) * scale
		offsetX := (it.Width - drawnW) / 2.0

		op.GeoM.Translate(screenX+offsetX, screenY)
		screen.DrawImage(img, op)
	} else {
		// Fallback vector drawing
		if it.Type == ItemIceCream {
			vector.DrawFilledCircle(screen, float32(screenX+9), float32(screenY+8), 6, color.RGBA{255, 180, 200, 255}, false)
			vector.StrokeLine(screen, float32(screenX+9), float32(screenY+14), float32(screenX+4), float32(screenY+24), 2, color.RGBA{200, 140, 60, 255}, false)
			vector.StrokeLine(screen, float32(screenX+9), float32(screenY+14), float32(screenX+14), float32(screenY+24), 2, color.RGBA{200, 140, 60, 255}, false)
		} else {
			vector.DrawFilledRect(screen, float32(screenX+5), float32(screenY+6), 8, 18, color.RGBA{80, 220, 160, 220}, false)
			vector.DrawFilledRect(screen, float32(screenX+7), float32(screenY+2), 4, 4, color.RGBA{40, 100, 220, 255}, false)
		}
	}
}
