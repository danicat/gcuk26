package render

import (
	"image/color"
	"london-eco-rider/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Background handles flat 2D London skyline scrolling.
type Background struct {
	Width, Height float64
}

// NewBackground creates background renderer.
func NewBackground(w, h float64) *Background {
	return &Background{Width: w, Height: h}
}

// Draw renders sky gradient, flat London panorama, and road pavement.
func (bg *Background) Draw(screen *ebiten.Image, cameraX float64) {
	roadY := float32(bg.Height * 0.72)

	// Layer 1: Sky Gradient (Soft London Sky Blue)
	vector.DrawFilledRect(screen, 0, 0, float32(bg.Width), roadY, color.RGBA{135, 200, 245, 255}, false)

	// Layer 2: Flat London Panorama Banner (Parallax 0.25x)
	bgImg := assets.LondonPanoramaImg

	if bgImg != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := bgImg.Bounds()

		// Scale banner to fill from y = 30 to roadY (top of pavement)
		bannerHeight := float64(roadY) - 30.0
		scaleY := bannerHeight / float64(bounds.Dy())
		scaleX := scaleY
		op.GeoM.Scale(scaleX, scaleY)

		scaledWidth := float64(bounds.Dx()) * scaleX
		parallaxX := -cameraX * 0.25
		startTile := int(parallaxX / scaledWidth)

		for tile := startTile - 2; tile <= startTile+3; tile++ {
			tileOp := *op
			tileOp.GeoM.Translate(parallaxX+float64(tile)*scaledWidth, 30.0)
			screen.DrawImage(bgImg, &tileOp)
		}
	}

	// Layer 3: Pavement / Street (Parallax 1.0x)
	roadHeight := float32(bg.Height - float64(roadY))

	// Dark asphalt road
	vector.DrawFilledRect(screen, 0, roadY, float32(bg.Width), roadHeight, color.RGBA{45, 50, 60, 255}, false)

	// Curb
	vector.DrawFilledRect(screen, 0, roadY-8, float32(bg.Width), 8, color.RGBA{180, 185, 195, 255}, false)

	// Dashed white road markings
	dashWidth := float32(32.0)
	dashGap := float32(22.0)
	offset := float32(int(-cameraX) % int(dashWidth+dashGap))

	for x := offset - dashWidth; x < float32(bg.Width)+dashWidth; x += dashWidth + dashGap {
		vector.DrawFilledRect(screen, x, roadY+24, dashWidth, 4, color.RGBA{255, 255, 255, 230}, false)
	}
}
