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
	ActionDrag
	ActionClick
	ActionLongClick
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	pos         input.Vec
	fallbackPos input.Vec

	numDrags    int
	numTaps     int
	numLongTaps int

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
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("use drag gesture to move the star\nnum drags: %d\nnum taps: %d\nnum long taps: %d", g.numDrags, g.numTaps, g.numLongTaps))
	ebitenutil.DebugPrintAt(screen, "*", int(g.pos.X), int(g.pos.Y))
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()

	if !g.started {
		g.Init()
		g.started = true
	}

	if g.inputHandler.ActionIsJustPressed(ActionClick) {
		g.numTaps++
	}
	if g.inputHandler.ActionIsJustPressed(ActionLongClick) {
		g.numLongTaps++
	}

	if info, ok := g.inputHandler.JustPressedActionInfo(ActionDrag); ok {
		// Start dragging.
		g.numDrags++
		g.fallbackPos = info.StartPos
	} else if info, ok := g.inputHandler.PressedActionInfo(ActionDrag); ok {
		// Continue dragging.
		g.pos = info.Pos
	} else {
		// Not being dragged.
		g.pos = g.fallbackPos
	}

	return nil
}

func (g *exampleGame) Init() {
	g.pos = input.Vec{X: 200, Y: 200}
	g.fallbackPos = g.pos

	keymap := input.Keymap{
		ActionDrag:      {input.KeyTouchDrag},
		ActionClick:     {input.KeyTouchTap},
		ActionLongClick: {input.KeyTouchLongTap},
	}

	g.inputHandler = g.inputSystem.NewHandler(0, keymap)
}
