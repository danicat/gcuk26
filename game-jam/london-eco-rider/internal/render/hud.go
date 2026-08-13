package render

import (
	"fmt"
	"image/color"
	"london-eco-rider/internal/entity"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// HUD renders score, mode status, recycling meter, and prompts.
type HUD struct {
	Width, Height float64
}

// NewHUD creates HUD manager.
func NewHUD(w, h float64) *HUD {
	return &HUD{Width: w, Height: h}
}

// Draw renders HUD elements.
func (h *HUD) Draw(screen *ebiten.Image, score int, bottles int, maxBottles int, turboSec float64, timeSec float64, isFever bool, mode entity.PlayerMode, promptMsg string) {
	// Top HUD Panel Bar
	vector.DrawFilledRect(screen, 10, 10, 260, 60, color.RGBA{20, 20, 30, 220}, false)
	vector.StrokeRect(screen, 10, 10, 260, 60, 2, color.RGBA{0, 200, 255, 255}, false)

	// Score & Time Text
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SCORE: %06d", score), 20, 16)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TIME:  %02d:%02d", int(timeSec)/60, int(timeSec)%60), 160, 16)

	// Player Mode Badge
	modeStr := "🚶 ON FOOT"
	if mode == entity.PlayerModeOnBike {
		modeStr = "🚴 BICYCLE"
	}
	ebitenutil.DebugPrintAt(screen, modeStr, 20, 36)

	// Recycling Bottle Bar
	recycleRatio := float32(bottles) / float32(maxBottles)
	if recycleRatio > 1.0 {
		recycleRatio = 1.0
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("BOTTLES: %d/%d", bottles, maxBottles), 110, 36)
	vector.DrawFilledRect(screen, 195, 38, 65, 10, color.RGBA{50, 50, 50, 255}, false)

	barCol := color.RGBA{80, 220, 100, 255}
	if isFever {
		barCol = color.RGBA{255, 215, 0, 255}
	}
	vector.DrawFilledRect(screen, 195, 38, 65*recycleRatio, 10, barCol, false)
	vector.StrokeRect(screen, 195, 38, 65, 10, 1, color.RGBA{255, 255, 255, 255}, false)

	// Interactive Story Prompt Banner
	if promptMsg != "" {
		vector.DrawFilledRect(screen, float32(h.Width/2)-180, 280, 360, 32, color.RGBA{15, 23, 42, 230}, false)
		vector.StrokeRect(screen, float32(h.Width/2)-180, 280, 360, 32, 2, color.RGBA{255, 200, 50, 255}, false)
		ebitenutil.DebugPrintAt(screen, promptMsg, int(h.Width/2)-160, 288)
	}

	// Turbo Boost Meter
	if turboSec > 0 {
		vector.DrawFilledRect(screen, float32(h.Width)-180, 10, 170, 30, color.RGBA{255, 100, 20, 220}, false)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("🍦 TURBO BOOST! %.1fs", turboSec), int(h.Width)-170, 18)
	}

	// Eco Fever Banner
	if isFever {
		vector.DrawFilledRect(screen, float32(h.Width/2)-100, 10, 200, 28, color.RGBA{0, 200, 100, 230}, false)
		ebitenutil.DebugPrintAt(screen, "🌟 ECO FEVER 2X MULTIPLIER! 🌟", int(h.Width/2)-92, 18)
	}
}
