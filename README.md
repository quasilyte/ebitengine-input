## Ebitengine input library

### Overview

A [Godot](https://godotengine.org/)-inspired action input handling system for [Ebitengine](https://github.com/hajimehoshi/ebiten).

**Key features:**

* [Actions](https://docs.godotengine.org/en/stable/tutorials/inputs/inputevent.html#actions) paradigm instead of the raw input events
* Configurable keymaps
* Bind more than one key to a single action
* Bind keys with modifiers to a single action (like `ctrl+c`)
* Simplified multi-input handling (like multiple gamepads)
* No extra dependencies (apart from the [Ebitengine](https://github.com/hajimehoshi/ebiten) of course)
* Solves some issues related to gamepads in browsers
* Can be used without extra deps or with [gmath](https://github.com/quasilyte/gmath) integration

### Installation

```bash
go get github.com/quasilyte/ebitengine-input
```

A runnable [example](_examples/basic/main.go) is available:

```bash
git clone https://github.com/quasilyte/ebitengine-input.git
cd ebitengine-input
go run ./_examples/basic/main.go
```

### Quick Start

```go
package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionMoveLeft input.Action = iota
	ActionMoveRight
)

func main() {
	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	p           *player
	inputSystem input.System
}

func newExampleGame() *exampleGame {
	g := &exampleGame{}
	g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyInput,
	})
	keymap := input.Keymap{
		ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
		ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
	}
	g.p = &player{
		input: g.inputSystem.NewHandler(0, keymap),
		pos:   image.Point{X: 96, Y: 96},
	}
	return g
}

func (g *exampleGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *exampleGame) Draw(screen *ebiten.Image) {
	g.p.Draw(screen)
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()
	g.p.Update()
	return nil
}

type player struct {
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
}

func (p *player) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "player", p.pos.X, p.pos.Y)
}
```

### Introduction

Let's assume that we have a simple game where you can move a character left or right.

You might end up checking the specific key events in your code like this:

```go
if ebiten.IsKeyPressed(ebiten.KeyLeft) {
    // Move left
}
```

But there are a few issues here:

1. This approach doesn't allow a key rebinding for the user
2. There is no clean way to add a gamepad support without making things messy
3. And even if you add a gamepad support, how would you handle multiple gamepads?

All of these issues can be solved by our little library. First, we need to declare our abstract actions as enum-like constants:

```go
const (
	ActionUnknown input.Action = iota
	ActionMoveLeft
	ActionMoveRight
)
```

Then we change the keypress handling code to this:

```go
if h.ActionIsPressed(ActionMoveLeft) {
    // Move left
}
```

Now, what is `h`? It's an [`input.Handler`](https://pkg.go.dev/github.com/quasilyte/ebitengine-input#Handler).

The input handler is bound to some keymap and device ID (only useful for the multi-devices setup with multiple gamepads being connected to the computer).

Having a keymap solves the first issue. The keymap associates an [`input.Action`](https://pkg.go.dev/github.com/quasilyte/ebitengine-input#Action) with a list of [`input.Key`](https://pkg.go.dev/github.com/quasilyte/ebitengine-input#Key). This means that the second issue is resolved too. The third issue is covered by the bound device ID.

So how do we create an input handler? We use a constructor provided by the [`input.System`](https://pkg.go.dev/github.com/quasilyte/ebitengine-input#System).

```go
// The ID argument is important for devices like gamepads.
// The input handlers can have the same keymaps.
player1input := inputSystem.NewHandler(0, keymap)
player2input := inputSystem.NewHandler(1, keymap)
```

The input system is an object that you integrate into your game `Update()` loop.

```go
func (g *myGame) Update() {
    g.inputSystem.Update() // called every Update()

    // ...rest of the function
}
```

You usually put this object into the game state. It could be either a global state (which I don't recommend) or a part of the state-like object that you pass through your game explicitely.

```go
type myGame struct {
    inputSystem input.System

    // ...rest of the fields
}
```

You'll need to call the `input.System.Init()` once before calling its `Update()` method. This `Init()` can be called **before** Ebitengine game is executed.

```go
func newMyGame() *myGame {
    g := &myGame{}
    g.inputSystem.Init(input.SystemConfig{
		DevicesEnabled: input.AnyInput,
	})
    // ... rest of the game object initialization
    return g
}
```

The keymaps are quite straightforward. We're hardcoding the keymap here, but it could be read from the config file.

```go
keymap := input.Keymap{
    ActionMoveLeft:  {input.KeyGamepadLeft, input.KeyLeft, input.KeyA},
    ActionMoveRight: {input.KeyGamepadRight, input.KeyRight, input.KeyD},
}
```

With the keymap above, when we check for the `ActionMoveLeft`, it doesn't matter if it was activated by a gamepad left button on a D-pad or by a keyboard left/A key.

Another benefit of this system is that we can get a list of relevant key events that can activate a given action. This is useful when you want to prompt player to press some button.

```go
// If gamepad is connected, show only gamepad-related keys.
// Otherwise show only keyboard-related keys.
inputDeviceMask := input.KeyboardInput
if h.GamepadConnected() {
    inputDeviceMask = GamepadInput
}
keyNames := h.ActionKeyNames(ActionMoveLeft, inputDeviceMask)
```

Since the pattern above is quite common, there is a shorthand for that:

```go
keyNames := h.ActionKeyNames(ActionMoveLeft, h.DefaultInputMask())
```

If the gamepad is connected, the `keyNames` will be `["gamepad_left"]`. Otherwise it will contain two entries for our example: `["left", "a"]`.

To build a combined key like `ctrl+c`, use `KeyWithModifier` function:

```go
// trigger an action when c is pressed while ctrl is down
input.KeyWithModifier(input.KeyC, input.ModControl)
```

See an [example](_examples/basic/main.go) for a complete source code.
