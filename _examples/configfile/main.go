//go:build example

package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

// Here is a basic algorithm to load a keymap from a file:
//
// 1. Store an action=>[]keyname mapping somewhere in a file;
//
// 2. When initializing an input.Keymap, you need to associate
//    an action string key with an actual input.Action constant;
//
// 3. Map keyname to input.Key, this can be done by using input.ParseKey function.
//
// Only the 2nd step requires some extra efforts.
// Since input.Action is an external type, you can't use a stringer tool to
// generate the string mappings. You can try using some other tool to do that.
// Or you can write the mapping manually (see actionString).
//
// When you have an action=>actionname mapping, it's easy to build a
// reverse index for the second step. A special sentinel value like actionLast
// can be useful (see the code below).

const (
	ActionUnknown input.Action = iota
	ActionLeft
	ActionRight
	ActionPause
	ActionRestart
	ActionSecret

	actionLast
)

func actionString(a input.Action) string {
	// This is the only function we have to implement.
	// Write it manually or use the tools to generate it (not stringer though).
	switch a {
	case ActionLeft:
		return "Left"
	case ActionRight:
		return "Right"
	case ActionPause:
		return "Pause"
	case ActionRestart:
		return "Restart"
	case ActionSecret:
		return "Secret"
	default:
		return "?"
	}
}

//go:embed keymap.json
var keymapConfigData []byte

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	lastActionPressed string

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
	ebitenutil.DebugPrint(screen, "last pressed action: "+g.lastActionPressed)
}

func (g *exampleGame) Update() error {
	g.inputSystem.Update()

	if !g.started {
		g.Init()
		g.started = true
	}

	actions := [...]input.Action{
		ActionLeft,
		ActionRight,
		ActionPause,
		ActionRestart,
		ActionSecret,
	}
	for _, a := range actions {
		if g.inputHandler.ActionIsJustPressed(a) {
			g.lastActionPressed = actionString(a)
			break
		}
	}

	return nil
}

func (g *exampleGame) Init() {
	var keymapConfig map[string][]string
	if err := json.Unmarshal(keymapConfigData, &keymapConfig); err != nil {
		panic(err)
	}

	// Build a reverse index to get an action ID by its name.
	actionNameToID := map[string]input.Action{}
	for a := ActionUnknown; a < actionLast; a++ {
		actionNameToID[actionString(a)] = a
	}

	// Parse our config file into a keymap object.
	keymap := input.Keymap{}
	for actionName, keyNames := range keymapConfig {
		a, ok := actionNameToID[actionName]
		if !ok {
			panic(fmt.Sprintf("unexpected action name: %s", actionName))
		}
		keys := make([]input.Key, len(keyNames))
		for i, keyString := range keyNames {
			k, err := input.ParseKey(keyString)
			if err != nil {
				panic(err)
			}
			keys[i] = k
		}
		keymap[a] = keys
	}

	g.inputHandler = g.inputSystem.NewHandler(0, keymap)
}
