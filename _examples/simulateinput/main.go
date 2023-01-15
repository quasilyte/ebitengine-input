//go:build example

package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

// The virtual input API are most useful when you want to:
//
// * Simulate the user input (specific device input event)
// * Test the game by triggering bound actions programmatically
// * Implement a remote controller object
//
// There are two main ways to emit such an event:
//
// 1. Handler.EmitKeyEvent(...)
// 2. Handler.EmitEvent(...)
//
// The first option requires a Key object to be specified.
// It makes the system believe that this key was actually in its pressed
// state during the frame. All device-related info is preserved.
// This is really good for the user input emulation or
// for the games that want to know which input device was used to
// emit an event. All device-related behavior will be preserved too.
// So, a gamepad button press will be controller-local, but
// keyboard events like KeyEnter will be visible to all handlers.
//
// The second option just triggers an Action directly.
// There will be no input device associated with that event.
// This means that the info object methods like IsGamepadEvent() and alike
// will always return false. It's possible to trigger an action that has
// none keys associated with it.
// All artificial actions triggered this way are only visible to
// the handlers of the same player ID. So they're always handler-local
// (like the gamepad button press would be).

const (
	ActionUnknown input.Action = iota
	ActionSpace
	ActionEnter
	ActionClick
	ActionGamepadButton
	ActionUnbound
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

	firstHandler  *input.Handler
	secondHandler *input.Handler
	inputSystem   input.System
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
		g.firstHandler.EmitKeyEvent(input.SimulatedKeyEvent{Key: input.KeySpace})
		g.firstHandler.EmitKeyEvent(input.SimulatedKeyEvent{Key: input.KeyGamepadA})
		g.firstHandler.EmitKeyEvent(input.SimulatedKeyEvent{
			Key: input.KeyMouseLeft,
			Pos: input.Vec{X: 100, Y: 100},
		})
		g.pressingEnter = !g.pressingEnter
		fmt.Printf(">> frame %d: switch 'pressing enter' (now %v)\n", g.frame, g.pressingEnter)
		fmt.Printf(">> frame %d: simulate space key press\n", g.frame)
		fmt.Printf(">> frame %d: simulate gamepad A key press\n", g.frame)
		fmt.Printf(">> frame %d: simulate a mouse click at (100,100)\n", g.frame)
	}
	if g.pressingEnter {
		// It's possible to trigger an action directly.
		// This would create a special
		g.firstHandler.EmitEvent(input.SimulatedAction{Action: ActionEnter})
	}
	if g.frame%100 == 0 {
		g.firstHandler.EmitEvent(input.SimulatedAction{
			Action: ActionUnbound,
			Pos:    input.Vec{X: 1, Y: 2},
		})
		fmt.Printf(">> frame %d: simulate an unbound action\n", g.frame)
	}

	if info, ok := g.firstHandler.JustPressedActionInfo(ActionClick); ok {
		fmt.Printf("<< frame %d: click action is just pressed (pos: %f,%f)\n", g.frame, info.Pos.X, info.Pos.Y)
	}
	if g.firstHandler.ActionIsPressed(ActionSpace) {
		fmt.Printf("<< frame %d: space action is pressed\n", g.frame)
	}
	if g.firstHandler.ActionIsJustPressed(ActionEnter) {
		fmt.Printf("<< frame %d: enter action is just pressed\n", g.frame)
	}
	if g.firstHandler.ActionIsJustPressed(ActionUnbound) {
		info, _ := g.firstHandler.JustPressedActionInfo(ActionUnbound)
		fmt.Printf("<< frame %d: unbound action is just pressed (pos: %f,%f)\n",
			g.frame, info.Pos.X, info.Pos.Y)
	}
	if g.firstHandler.ActionIsPressed(ActionGamepadButton) {
		fmt.Printf("<< frame %d: gamepad button action is pressed\n", g.frame)
	}

	// Gamepad actions are bound to the device ID (same as player ID).
	// Therefore, the second handler doesn't get a gamepad button event.
	// But it's possible to trigger this action from the second gamepad!
	if g.secondHandler.ActionIsJustPressed(ActionGamepadButton) {
		fmt.Printf("<<<< frame %d: gamepad button action is pressed\n", g.frame)
	}
	// Artificial input events are id-bound too.
	if g.secondHandler.ActionIsJustPressed(ActionUnbound) {
		panic("should never happen")
	}
	// None of the simulated actions are broadcasted.
	// Even things that are bound to keys like "enter".
	// But it's possible to trigger this action by pressing the enter manually!
	if g.secondHandler.ActionIsJustPressed(ActionEnter) {
		fmt.Printf("<<<< frame %d: enter button action is just pressed\n", g.frame)
	}

	return nil
}

func (g *exampleGame) Init() {
	keymap := input.Keymap{
		ActionClick:         {input.KeyMouseLeft},
		ActionSpace:         {input.KeySpace},
		ActionEnter:         {input.KeyEnter},
		ActionGamepadButton: {input.KeyGamepadA},
		ActionUnbound:       {}, // Empty by choice
	}

	g.pressingEnter = true
	g.firstHandler = g.inputSystem.NewHandler(0, keymap)
	g.secondHandler = g.inputSystem.NewHandler(1, keymap)
}
