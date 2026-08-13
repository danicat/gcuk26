package game

import (
	"london-eco-rider/assets"
	"london-eco-rider/internal/audio"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game implements the ebiten.Game interface.
type Game struct {
	currentState StateType
	states       map[StateType]State
	Audio        *audio.AudioSystem
	FinalScore   int
}

// NewGame initializes assets, audio system, and states.
func NewGame() *Game {
	assets.LoadAssets()
	audioSys := audio.NewAudioSystem()

	g := &Game{
		currentState: StateTitle,
		states:       make(map[StateType]State),
		Audio:        audioSys,
	}

	g.states[StateTitle] = NewTitleState()
	g.states[StatePlay] = NewPlayState()
	g.states[StateGameOver] = NewGameOverState(false, 0)
	g.states[StateGameClear] = NewGameOverState(true, 0)

	g.states[g.currentState].Enter()
	g.Audio.StartBGM()

	return g
}

// Update handles state transitions and main loop updates.
func (g *Game) Update() error {
	// Fullscreen toggle (F11 or Alt+Enter)
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	nextState := g.states[g.currentState].Update(g)

	if nextState != g.currentState {
		g.states[g.currentState].Exit()

		if nextState == StateGameOver {
			g.states[StateGameOver] = NewGameOverState(false, g.FinalScore)
		} else if nextState == StateGameClear {
			g.states[StateGameClear] = NewGameOverState(true, g.FinalScore)
		}

		g.currentState = nextState
		g.states[g.currentState].Enter()
	}

	return nil
}

// Draw renders current state.
func (g *Game) Draw(screen *ebiten.Image) {
	g.states[g.currentState].Draw(screen)
}

// Layout returns the 16:9 virtual canvas resolution (640x360).
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
