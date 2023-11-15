//go:build example

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionUnknown input.Action = iota
	ActionAdd
	ActionSub
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	score int

	inputHandler *input.Handler
	inputSystem  input.System
}

func newExampleGame() *exampleGame {
	g := &exampleGame{}

	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyDevice,
	})

	return g
}

func (g *exampleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *exampleGame) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, strings.Join([]string{
		"mouse controls:",
		"  released [left mouse button]: increase score",
		"  released [ctrl]+[left mouse button]: decrease score",
		"keyboard controls:",
		"  released [enter]: increase score",
		"  released [ctrl]+[enter]: decrease score",
		"gamepad controls:",
		"  released [R1]: increase score",
		"  released [L1]: increase score",
		"",
		fmt.Sprintf("score: %d", g.score),
	}, "\n"))
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()

	if !g.started {
		g.Init()
		g.started = true
	}

	g.handleInput()

	return nil
}

func (g *exampleGame) handleInput() {
	// Due to the fact that ActionAdd requires just an LMB while
	// ActionSub is ctrl+LMB, the ActionAdd can be confused with
	// ActionSub if it's checked first.
	// This is due to the fact that ebitengine-input does no
	// conflict resolution of any kind.
	// It may start preferring the "longest" key combination
	// in the future, but for now, you might want to order
	// the key checks carefully if your keymap might have these conflicts.
	// See #36 for more details.
	if g.inputHandler.ActionIsJustReleased(ActionSub) {
		g.score--
	} else if g.inputHandler.ActionIsJustReleased(ActionAdd) {
		g.score++
	}
}

func (g *exampleGame) Init() {
	keymap := input.Keymap{
		ActionAdd: {input.KeyMouseLeft, input.KeyEnter, input.KeyGamepadR1},
		ActionSub: {
			input.KeyWithModifier(input.KeyMouseLeft, input.ModControl),
			input.KeyWithModifier(input.KeyEnter, input.ModControl),
			input.KeyGamepadL1,
		},
	}

	g.inputHandler = g.inputSystem.NewHandler(0, keymap)
}
