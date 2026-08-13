package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameOverState struct {
	IsVictory  bool
	FinalScore int
}

func NewGameOverState(isVictory bool, score int) *GameOverState {
	return &GameOverState{
		IsVictory:  isVictory,
		FinalScore: score,
	}
}

func (s *GameOverState) Enter() {}

func (s *GameOverState) Update(g *Game) StateType {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return StateTitle
	}
	return g.currentState
}

func (s *GameOverState) Draw(screen *ebiten.Image) {
	// Dark backdrop
	vector.DrawFilledRect(screen, 0, 0, 640, 360, color.RGBA{10, 15, 30, 240}, false)

	// Result Box
	vector.DrawFilledRect(screen, 120, 60, 400, 240, color.RGBA{25, 35, 60, 240}, false)

	if s.IsVictory {
		vector.StrokeRect(screen, 120, 60, 400, 240, 3, color.RGBA{80, 220, 100, 255}, false)
		ebitenutil.DebugPrintAt(screen, "🎉 STAGE CLEAR! LONDON IS CLEANER! 🎉", 170, 90)
	} else {
		vector.StrokeRect(screen, 120, 60, 400, 240, 3, color.RGBA{220, 60, 60, 255}, false)
		ebitenutil.DebugPrintAt(screen, "⏰ TIME EXPIRED! GAME OVER ⏰", 210, 90)
	}

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FINAL ECO SCORE: %06d", s.FinalScore), 230, 150)
	ebitenutil.DebugPrintAt(screen, "Thank you for recycling plastic bottles!", 180, 190)

	ebitenutil.DebugPrintAt(screen, "PRESS SPACE OR ENTER TO RESTART", 195, 240)
}

func (s *GameOverState) Exit() {}
