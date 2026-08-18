package game

import (
	"image/color"
	"london-eco-rider/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type TitleState struct{}

func NewTitleState() *TitleState {
	return &TitleState{}
}

func (s *TitleState) Enter() {}

func (s *TitleState) Update(g *Game) StateType {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return StatePlay
	}
	return StateTitle
}

func (s *TitleState) Draw(screen *ebiten.Image) {
	// Background Panorama Banner
	if assets.LondonPanoramaImg != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := assets.LondonPanoramaImg.Bounds()
		scaleX := 640.0 / float64(bounds.Dx())
		scaleY := 360.0 / float64(bounds.Dy())
		op.GeoM.Scale(scaleX, scaleY)
		screen.DrawImage(assets.LondonPanoramaImg, op)
	} else {
		vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{100, 180, 240, 255}, false)
	}

	// Darkened Translucent Overlay
	vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{10, 15, 30, 210}, false)

	// Main Title Card
	vector.DrawFilledRect(screen, 60, 30, 520, 300, color.RGBA{20, 30, 50, 235}, false)
	vector.StrokeRect(screen, 60, 30, 520, 300, 2, color.RGBA{0, 200, 255, 255}, false)

	// Header Title
	ebitenutil.DebugPrintAt(screen, "🚴 LONDON ECO-RIDER 🍦", 235, 45)
	vector.StrokeLine(screen, 80, 62, 560, 62, 1, color.RGBA{0, 200, 255, 180}, false)

	// Short Goal Summary
	ebitenutil.DebugPrintAt(screen, "GOAL: Collect bottles -> Unlock Bike at Bike Shop -> Patrol London!", 90, 75)

	// Item Showcase Card
	vector.DrawFilledRect(screen, 75, 100, 490, 140, color.RGBA{30, 45, 75, 240}, false)
	vector.StrokeRect(screen, 75, 100, 490, 140, 1, color.RGBA{80, 180, 255, 200}, false)

	// 4 Concise Showcase Items
	drawShowcaseItemShort(screen, assets.PlasticBottleImg, 90, 115, 36, "BOTTLE", "+150 Pts")
	drawShowcaseItemShort(screen, assets.IceCreamImage, 205, 115, 36, "ICE CREAM", "Turbo Speed")
	drawShowcaseItemShort(screen, assets.BoyBikeImage, 320, 115, 36, "BICYCLE", "Pick Color")
	drawShowcaseItemShort(screen, assets.RedBusImage, 435, 115, 36, "RED BUS", "Jump Over!")

	// Solid Steady Non-Flashing Start Button
	vector.DrawFilledRect(screen, 180, 255, 280, 28, color.RGBA{0, 180, 240, 240}, false)
	vector.StrokeRect(screen, 180, 255, 280, 28, 1, color.RGBA{255, 255, 255, 255}, false)
	ebitenutil.DebugPrintAt(screen, "PRESS SPACE TO START", 230, 262)

	// Controls Footer
	ebitenutil.DebugPrintAt(screen, "Controls: Move (Arrow Keys / A-D) | Jump (Space) | Turbo (Shift)", 110, 300)
}

func drawShowcaseItemShort(screen *ebiten.Image, img *ebiten.Image, x, y float64, size float64, title string, desc string) {
	vector.DrawFilledRect(screen, float32(x), float32(y), 100, 110, color.RGBA{15, 23, 42, 210}, false)
	vector.StrokeRect(screen, float32(x), float32(y), 100, 110, 1, color.RGBA{100, 150, 200, 150}, false)

	if img != nil {
		op := &ebiten.DrawImageOptions{}
		bounds := img.Bounds()
		scale := size / float64(bounds.Dy())
		op.GeoM.Scale(scale, scale)

		drawnW := float64(bounds.Dx()) * scale
		offsetX := (100.0 - drawnW) / 2.0

		op.GeoM.Translate(x+offsetX, y+10.0)
		screen.DrawImage(img, op)
	}

	ebitenutil.DebugPrintAt(screen, title, int(x)+12, int(y)+60)
	ebitenutil.DebugPrintAt(screen, desc, int(x)+8, int(y)+80)
}

func (s *TitleState) Exit() {}
