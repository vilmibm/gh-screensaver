package savers

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

type FireworksSaver struct {
	screen tcell.Screen
	style  tcell.Style
	color  bool
	inputs map[string]string
}

func NewFireworksSaver(opts shared.ScreensaverOpts) (shared.Screensaver, error) {
	fs := &FireworksSaver{}
	if err := fs.Initialize(opts); err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *FireworksSaver) Initialize(opts shared.ScreensaverOpts) error {
	fs.screen = opts.Screen
	fs.style = opts.Style

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func (fs *FireworksSaver) Inputs() map[string]shared.SaverInput {
	return map[string]shared.SaverInput{
		"color": {
			Default:     "full",
			Description: "whether to use full color or monochrome. Values: full, off",
		},
	}
}

func (fs *FireworksSaver) SetInputs(inputs map[string]string) error {
	fs.inputs = inputs
	if inputs["color"] == "full" {
		fs.color = true
	}
	return nil
}

func (fs *FireworksSaver) Update() error {
	// TODO
	return nil
}
