//go:build example

package main

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

// For the basics, see "_examples/basic";
// this example omits some explanations for brevity.

const (
	ActionUnknown input.Action = iota
	ActionMoveLeft
	ActionMoveRight
	ActionDebug
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	state gameState

	currentScene scene

	inputSystem input.System
}

type gameState struct {
	inputHandlers []*input.Handler
}

type scene interface {
	Update() scene
	Draw(screen *ebiten.Image)
}

func newExampleGame() *exampleGame {
	g := &exampleGame{}

	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyInput,
	})

	return g
}

func (g *exampleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()
	if !g.started {
		g.init()
		g.started = true
	}
	g.currentScene = g.currentScene.Update()
	return nil
}

func (g *exampleGame) Draw(screen *ebiten.Image) {
	g.currentScene.Draw(screen)
}

func (g *exampleGame) init() {
	keymap := input.Keymap{
		// ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyGamepadLStickLeft},
		// ActionMoveRight: {input.KeyGamepadRight, input.KeyGamepadLStickRight},
		ActionDebug: {input.KeyGamepadLeft},
	}

	g.state.inputHandlers = make([]*input.Handler, 4)
	for i := range g.state.inputHandlers {
		g.state.inputHandlers[i] = g.inputSystem.NewHandler(i, keymap)
	}

	g.currentScene = &lobbyScene{
		state:   &g.state,
		timeout: time.Now().Add(9 * time.Second),
	}
}

type lobbyScene struct {
	state       *gameState
	timeout     time.Time
	secondsLeft float64
	gamepads    int
}

func (s *lobbyScene) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "connect gamepads, press buttons", 200, 160)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("the scene changes in %.1f seconds", s.secondsLeft), 200, 200)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("gamepads detected: %d", s.gamepads), 240, 240)

}

func (s *lobbyScene) Update() scene {
	s.secondsLeft = time.Until(s.timeout).Seconds()

	numGamepads := 0
	for _, h := range s.state.inputHandlers {
		if h.GamepadConnected() {
			numGamepads++
		}
	}
	s.gamepads = numGamepads

	if s.secondsLeft <= 0 {
		return newMainScene(s.state, s.gamepads)
	}
	return s
}

type mainScene struct {
	state   *gameState
	players []*player
}

func newMainScene(state *gameState, gamepads int) *mainScene {
	s := &mainScene{state: state}

	// Depending on the actual number of players, create
	// player objects and give them associated input handlers.
	s.players = make([]*player, gamepads)
	pos := image.Point{X: 256, Y: 128}
	for i := range s.players {
		s.players[i] = &player{
			input: state.inputHandlers[i],
			pos:   pos,
			label: fmt.Sprintf("[player%d]", i+1),
		}
		pos.Y += 64
	}

	return s
}

func (s *mainScene) Draw(screen *ebiten.Image) {
	for _, p := range s.players {
		p.Draw(screen)
	}
}

func (s *mainScene) Update() scene {
	for _, p := range s.players {
		p.Update()
	}
	return s
}

type player struct {
	label string
	input *input.Handler
	pos   image.Point
}

func (p *player) Update() {
	if p.input.ActionIsPressed(ActionMoveLeft) {
		p.pos.X -= 4
	}
	if p.input.ActionIsPressed(ActionMoveRight) {
		p.pos.X += 4
	}
	if p.input.ActionIsJustPressed(ActionDebug) {
		fmt.Printf("%s: debug action is pressed\n", p.label)
	}
}

func (p *player) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, p.label, p.pos.X, p.pos.Y)
}
