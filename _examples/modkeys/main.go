//go:build example

package main

import (
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	input "github.com/quasilyte/ebitengine-input"
)

const (
	ActionUnknown input.Action = iota
	ActionAddGreenLine
	ActionAddBlueLine
	ActionResetGreenLine
	ActionResetBlueLine
)

func main() {
	ebiten.SetWindowSize(640, 480)

	if err := ebiten.RunGame(newExampleGame()); err != nil {
		log.Fatal(err)
	}
}

type exampleGame struct {
	started bool

	points []drawPoint

	inputHandler *input.Handler
	inputSystem  input.System
}

type drawPoint struct {
	pos   input.Vec
	color color.RGBA
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
	ebitenutil.DebugPrint(screen, "press lmb with ctrl/shift (any combo)")

	pos := input.Vec{X: 320, Y: 240}
	for _, pt := range g.points {
		drawLine(screen, pos, pt.pos, pt.color)
		pos = pt.pos
	}
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
	blue := color.RGBA{B: 0xf0, A: 0xff}
	green := color.RGBA{G: 0xf0, A: 0xff}
	if info, ok := g.inputHandler.JustPressedActionInfo(ActionAddGreenLine); ok {
		g.points = append(g.points, drawPoint{color: green, pos: info.Pos})
		return
	}
	if info, ok := g.inputHandler.JustPressedActionInfo(ActionResetGreenLine); ok {
		g.points = append(g.points[:0], drawPoint{color: green, pos: info.Pos})
		return
	}
	if info, ok := g.inputHandler.JustPressedActionInfo(ActionAddBlueLine); ok {
		g.points = append(g.points, drawPoint{color: blue, pos: info.Pos})
		return
	}
	if info, ok := g.inputHandler.JustPressedActionInfo(ActionResetBlueLine); ok {
		g.points = append(g.points[:0], drawPoint{color: blue, pos: info.Pos})
		return
	}
}

func (g *exampleGame) Init() {
	keymap := input.Keymap{
		ActionResetBlueLine: {input.KeyMouseLeft},
		ActionAddBlueLine:   {input.KeyWithModifier(input.KeyMouseLeft, input.ModControl)},

		ActionResetGreenLine: {input.KeyWithModifier(input.KeyMouseLeft, input.ModShift)},
		ActionAddGreenLine:   {input.KeyWithModifier(input.KeyMouseLeft, input.ModControlShift)},
	}

	g.inputHandler = g.inputSystem.NewHandler(0, keymap)
}

func drawLine(dst *ebiten.Image, pos1, pos2 input.Vec, c color.RGBA) {
	x1 := pos1.X
	y1 := pos1.Y
	x2 := pos2.X
	y2 := pos2.Y

	length := math.Hypot(x2-x1, y2-y1)

	var drawOptions ebiten.DrawImageOptions
	drawOptions.GeoM.Scale(length, 2)
	drawOptions.GeoM.Rotate(math.Atan2(y2-y1, x2-x1))
	drawOptions.GeoM.Translate(x1, y1)
	drawOptions.ColorScale.Scale(float32(c.R), float32(c.G), float32(c.B), float32(c.A))

	dst.DrawImage(whitePixel, &drawOptions)
}

var whitePixel *ebiten.Image

func init() {
	emptyImage := ebiten.NewImage(3, 3)
	emptyImage.Fill(color.White)
	whitePixel = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
}
