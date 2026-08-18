package game

import "github.com/hajimehoshi/ebiten/v2"

type StateType int

const (
	StateTitle StateType = iota
	StatePlay
	StateGameOver
	StateGameClear
)

// State interface for FSM scenes.
type State interface {
	Enter()
	Update(game *Game) StateType
	Draw(screen *ebiten.Image)
	Exit()
}
