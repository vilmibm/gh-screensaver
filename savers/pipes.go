package savers

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

type PipesSaver struct {
	screen tcell.Screen
	style  tcell.Style
	color  bool
	inputs map[string]string
}

func NewPipesSaver(opts shared.ScreensaverOpts) (shared.Screensaver, error) {
	ps := &PipesSaver{}
	if err := ps.Initialize(opts); err != nil {
		return nil, err
	}

	return ps, nil
}

func (ps *PipesSaver) Initialize(opts shared.ScreensaverOpts) error {
	ps.screen = opts.Screen
	ps.style = opts.Style

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func (fs *PipesSaver) Inputs() map[string]shared.SaverInput {
	// TODO eventually support truecolor
	return map[string]shared.SaverInput{
		"color": {
			Default:     "full",
			Description: "whether to use full color or monochrome. Values: full, off",
		},
	}
}

func (ps *PipesSaver) SetInputs(inputs map[string]string) error {
	ps.inputs = inputs
	if inputs["color"] == "full" {
		ps.color = true
	}
	return nil
}

type coord struct {
	x int
	y int
}

type pipe struct {
	coords []coord
	color  tcell.Color
}

// TODO validMoves

func randColor() tcell.Color {
	r := rand.Intn(255)
	g := rand.Intn(255)
	b := rand.Intn(255)

	return tcell.NewRGBColor(r, g, b)
}

func (ps *PipesSaver) Update() error {
	// TODO check new pipe chance
	//       - random color
	//       - random wall
	//       - add to list
	// TODO for each pipe,
	//      - decide on next spot options based on proximity to wall
	//      - pick random next spot

	// TODO the OG pipes clears itself at some interval. I think it will take far
	// more time for us to fill up a screen, so initially I think i'll just let
	// it fill up.
	return nil
}
