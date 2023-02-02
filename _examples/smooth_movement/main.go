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
	ActionMove
	ActionAlternativeMove
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	status   unitStatus
	pos      input.Vec
	startPos input.Vec

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
	ebitenutil.DebugPrint(screen, "move by using the gamepad left stick")

	// We'll use some ASCII art instead of the real graphics.
	var sprite string
	offsetX := -4
	offsetY := -6
	switch g.status {
	case statusIdle:
		sprite = "o"
		offsetX = 0
		offsetY = 0
	case statusMoving, statusStartMovement:
		sprite = "@@\n@@"
	}
	if g.startPos != (input.Vec{}) {
		ebitenutil.DebugPrintAt(screen, "O", int(g.startPos.X), int(g.startPos.Y))
	}
	ebitenutil.DebugPrintAt(screen, sprite, int(g.pos.X)+offsetX, int(g.pos.Y)+offsetY)
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()

	if !g.started {
		g.Init()
		g.started = true
	}

	// We simulate the left stick events using our right controller stick.
	// It's better just to bind both L+R sticks to the action, but this example
	// demonstrates the power of simulated events.
	if info, ok := g.inputHandler.PressedActionInfo(ActionAlternativeMove); ok {
		g.inputHandler.EmitKeyEvent(input.SimulatedKeyEvent{
			Key: input.KeyGamepadLStickMotion,
			Pos: info.Pos,
		})
	}

	// You can control all movement phases:
	// - its first frame (started to move)
	// - its end (just stopped to move)
	// - its active phase (on the move)
	// - the movement absence (an idle state)
	if g.inputHandler.ActionIsJustPressed(ActionMove) {
		// The movement is just started.
		g.startPos = g.pos
		g.status = statusStartMovement
	} else if info, ok := g.inputHandler.PressedActionInfo(ActionMove); ok {
		// We're in the middle of the movement.
		// The info.Pos value is like a direction vector of the stick.
		g.pos.X += info.Pos.X * 2
		g.pos.Y += info.Pos.Y * 2
		g.status = statusMoving
	} else if g.status != statusIdle {
		// The movement has just finished.
		g.startPos = input.Vec{}
		g.status = statusIdle
	} else {
		// No movement is being executed.
		g.status = statusIdle
	}

	return nil
}

func (g *exampleGame) Init() {
	g.pos = input.Vec{X: 200, Y: 200}

	keymap := input.Keymap{
		ActionMove:            {input.KeyGamepadLStickMotion},
		ActionAlternativeMove: {input.KeyGamepadRStickMotion},
	}

	g.inputHandler = g.inputSystem.NewHandler(0, keymap)
}

type unitStatus int

const (
	statusIdle unitStatus = iota
	statusStartMovement
	statusMoving
)
