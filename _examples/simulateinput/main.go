//go:build example

package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionUnknown input.Action = iota
	ActionSpace
	ActionEnter
	ActionClick
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	frame         int
	pressingEnter bool

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
	ebitenutil.DebugPrint(screen, "check out the stdout logs\ntry clicking, pressing space/enter")
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()

	if !g.started {
		g.Init()
		g.started = true
	}

	g.frame++
	// Every 90 frames, emit some events and toggle the enter pressing mode.
	// Note: the simulated input events won't be detected until the next frame.
	if g.frame%90 == 0 {
		g.inputHandler.EmitEvent(input.SimulatedEvent{
			Key: input.KeySpace,
		})
		g.inputHandler.EmitEvent(input.SimulatedEvent{
			Key: input.KeyMouseLeft,
			Pos: input.Vec{X: 100, Y: 100},
		})
		g.pressingEnter = !g.pressingEnter
		fmt.Printf(">> frame %d: switch 'pressing enter' (now %v)\n", g.frame, g.pressingEnter)
		fmt.Printf(">> frame %d: simulate space key press\n", g.frame)
		fmt.Printf(">> frame %d: simulate a mouse click at (100,100)\n", g.frame)
	}
	if g.pressingEnter {
		g.inputHandler.EmitEvent(input.SimulatedEvent{
			Key: input.KeyEnter,
		})
	}

	if info, ok := g.inputHandler.JustPressedActionInfo(ActionClick); ok {
		fmt.Printf("<< frame %d: click action is pressed (pos: %f,%f)\n", g.frame, info.Pos.X, info.Pos.Y)
	}
	if g.inputHandler.ActionIsPressed(ActionSpace) {
		fmt.Printf("<< frame %d: space action is pressed\n", g.frame)
	}
	if g.inputHandler.ActionIsJustPressed(ActionEnter) {
		fmt.Printf("<< frame %d: enter action is just pressed\n", g.frame)
	}

	return nil
}

func (g *exampleGame) Init() {
	keymap := input.Keymap{
		ActionClick: {input.KeyMouseLeft},
		ActionSpace: {input.KeySpace},
		ActionEnter: {input.KeyEnter},
	}

	g.pressingEnter = true
	g.inputHandler = g.inputSystem.NewHandler(0, keymap)
}
