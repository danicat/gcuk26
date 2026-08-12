package entity

import (
	"image/color"
	"london-eco-rider/assets"
	"london-eco-rider/internal/physics"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ObstacleType int

const (
	ObstacleRedBus ObstacleType = iota
	ObstaclePhoneBox
	ObstacleRamp
	ObstacleRecyclingBin
	ObstacleIceCreamVan
	ObstacleBikeShop
	ObstacleTree
)

// Obstacle represents moving red buses, phone boxes, ramps, recycling bins, ice cream vans, bike shops, or trees.
type Obstacle struct {
	Type    ObstacleType
	X, Y    float64
	VX      float64
	Width   float64
	Height  float64
	IsRamp  bool
	Stopped bool
}

// NewObstacle creates a London obstacle or interactive station.
func NewObstacle(obsType ObstacleType, x, y float64) *Obstacle {
	w, h := 40.0, 40.0
	isRamp := false
	vx := 0.0

	switch obsType {
	case ObstacleRedBus:
		w, h = 100.0, 52.0
		vx = -1.2 // Dynamic moving bus
	case ObstaclePhoneBox:
		w, h = 26.0, 52.0
	case ObstacleRamp:
		w, h = 42.0, 22.0
		isRamp = true
	case ObstacleRecyclingBin:
		w, h = 32.0, 42.0
	case ObstacleIceCreamVan:
		w, h = 88.0, 52.0
	case ObstacleBikeShop:
		w, h = 95.0, 60.0
	case ObstacleTree:
		w, h = 55.0, 75.0 // Lush London oak tree
	}

	return &Obstacle{
		Type:    obsType,
		X:       x,
		Y:       y - h,
		VX:      vx,
		Width:   w,
		Height:  h,
		IsRamp:  isRamp,
		Stopped: false,
	}
}

// Update advances moving obstacles (e.g. buses).
func (o *Obstacle) Update() {
	if !o.Stopped && o.VX != 0 {
		o.X += o.VX
	}
}

// Stop stops obstacle movement (e.g. upon bus collision).
func (o *Obstacle) Stop() {
	o.Stopped = true
	o.VX = 0
}

// GetAABB returns bounding box for collision detection.
func (o *Obstacle) GetAABB() physics.AABB {
	return physics.AABB{
		X:      o.X,
		Y:      o.Y,
		Width:  o.Width,
		Height: o.Height,
	}
}

// Draw renders obstacle sprite with drop shadow, without outline boxes.
func (o *Obstacle) Draw(screen *ebiten.Image, cameraX float64) {
	screenX := o.X - cameraX
	screenY := o.Y

	// Ground Drop Shadow on Asphalt Road
	vector.DrawFilledCircle(screen, float32(screenX+o.Width*0.5), float32(screenY+o.Height-2), float32(o.Width*0.45), color.RGBA{15, 15, 25, 140}, false)

	var img *ebiten.Image
	switch o.Type {
	case ObstacleRedBus:
		img = assets.RedBusImage
	case ObstaclePhoneBox:
		img = assets.RedPhoneBoxImg
	case ObstacleRecyclingBin:
		img = assets.RecyclingBinImg
	case ObstacleIceCreamVan:
		img = assets.IceCreamVanImg
	case ObstacleBikeShop:
		img = assets.BikeShopImg
	case ObstacleTree:
		img = assets.TreeImage
	}

	if img != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := img.Bounds()

		scale := o.Height / float64(bounds.Dy())
		op.GeoM.Scale(scale, scale)

		drawnW := float64(bounds.Dx()) * scale
		offsetX := (o.Width - drawnW) / 2.0

		op.GeoM.Translate(screenX+offsetX, screenY)
		screen.DrawImage(img, op)
	} else {
		// Fallbacks
		switch o.Type {
		case ObstacleRedBus:
			vector.DrawFilledRect(screen, float32(screenX), float32(screenY), float32(o.Width), float32(o.Height), color.RGBA{220, 30, 30, 255}, false)
		case ObstaclePhoneBox:
			vector.DrawFilledRect(screen, float32(screenX), float32(screenY), float32(o.Width), float32(o.Height), color.RGBA{200, 20, 20, 255}, false)
		case ObstacleRecyclingBin:
			vector.DrawFilledRect(screen, float32(screenX), float32(screenY), float32(o.Width), float32(o.Height), color.RGBA{40, 180, 80, 255}, false)
		case ObstacleIceCreamVan:
			vector.DrawFilledRect(screen, float32(screenX), float32(screenY), float32(o.Width), float32(o.Height), color.RGBA{255, 230, 180, 255}, false)
		case ObstacleBikeShop:
			vector.DrawFilledRect(screen, float32(screenX), float32(screenY), float32(o.Width), float32(o.Height), color.RGBA{50, 150, 220, 255}, false)
		case ObstacleTree:
			vector.DrawFilledCircle(screen, float32(screenX+o.Width*0.5), float32(screenY+o.Height*0.4), float32(o.Width*0.4), color.RGBA{40, 160, 60, 255}, false)
			vector.DrawFilledRect(screen, float32(screenX+o.Width*0.4), float32(screenY+o.Height*0.5), float32(o.Width*0.2), float32(o.Height*0.5), color.RGBA{120, 80, 30, 255}, false)
		case ObstacleRamp:
			path := vector.Path{}
			path.MoveTo(float32(screenX), float32(screenY+o.Height))
			path.LineTo(float32(screenX+o.Width), float32(screenY))
			path.LineTo(float32(screenX+o.Width), float32(screenY+o.Height))
			path.Close()
			vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
			for i := range vs {
				vs[i].ColorR = 0.9
				vs[i].ColorG = 0.6
				vs[i].ColorB = 0.2
				vs[i].ColorA = 1.0
			}
			op := &ebiten.DrawTrianglesOptions{}
			screen.DrawTriangles(vs, is, whiteSubImage, op)
		}
	}
}

var whiteSubImage = func() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}()
