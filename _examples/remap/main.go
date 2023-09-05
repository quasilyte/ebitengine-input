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
	ActionPing
	ActionRemap
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	k        input.Key
	prevK    input.Key
	scanning bool

	keyScanner input.KeyScanner

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
	if g.scanning {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("keybind: %s\n<scanning the new keybing>", g.k))
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("keybind: %s\npress ctrl+enter to remap", g.k))
	}
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()

	if !g.started {
		g.Init()
		g.started = true
	}

	g.handleRemap()

	if !g.scanning && g.inputHandler.ActionIsJustPressed(ActionPing) {
		fmt.Printf("ping! (activated with %s keybind)\n", g.k)
	}

	return nil
}

func (g *exampleGame) handleRemap() {
	if !g.scanning {
		if g.inputHandler.ActionIsJustPressed(ActionRemap) {
			g.prevK = g.k // Save it for an easier fallback
			g.scanning = true
		}
		return
	}

	k, status := g.keyScanner.Scan()
	if status != input.KeyScanUnchanged {
		g.k = k
	}
	if status == input.KeyScanCompleted {
		g.scanning = false
		// Check for the new key to be available.
		// Resolve the conflicts here; I'll just reject
		// the combination that is already in use.
		if g.k == input.KeyWithModifier(input.KeyEnter, input.ModControl) {
			g.k = g.prevK
		} else {
			g.inputHandler.Remap(g.makeKeymap())
		}
	}
}

func (g *exampleGame) makeKeymap() input.Keymap {
	return input.Keymap{
		ActionPing:  {g.k},
		ActionRemap: {input.KeyWithModifier(input.KeyEnter, input.ModControl)},
	}
}

func (g *exampleGame) Init() {
	g.k = input.KeyQ
	g.inputHandler = g.inputSystem.NewHandler(0, g.makeKeymap())
}
