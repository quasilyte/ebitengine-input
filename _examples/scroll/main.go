//go:build example

package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionUnknown input.Action = iota
	ActionScrollVertical
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	pos input.Vec

	inputHandler *input.Handler
	inputSystem  input.System
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

func (g *exampleGame) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "scroll up/down")
	ebitenutil.DebugPrintAt(screen, "*", int(g.pos.X), int(g.pos.Y))
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()

	if !g.started {
		g.Init()
		g.started = true
	}

	if info, ok := g.inputHandler.JustPressedActionInfo(ActionScrollVertical); ok {
		g.pos.Y += info.Pos.Y
	}

	return nil
}

func (g *exampleGame) Init() {
	g.pos = input.Vec{X: 200, Y: 200}

	keymap := input.Keymap{
		ActionScrollVertical: {input.KeyWheelVertical},
	}

	g.inputHandler = g.inputSystem.NewHandler(0, keymap)
}
