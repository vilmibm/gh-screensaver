package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/lukesampson/figlet/figletlib"
	"github.com/mattn/go-runewidth"
)

type MarqueeSaver struct {
	screen tcell.Screen
	style  tcell.Style
	x      int
	banner string
	inputs map[string]string
}

func NewMarqueeSaver(opts screensaverOpts) (Screensaver, error) {
	bs := &MarqueeSaver{}
	if err := bs.Initialize(opts); err != nil {
		return nil, err
	}
	width, _ := bs.screen.Size()
	// TODO vendor fonts
	f, err := figletlib.GetFontByName("/usr/share/figlet", "slant")
	if err != nil {
		return nil, err
	}
	// TODO parameterize message
	bs.banner = figletlib.SprintMsg("RUIN YOURSELF BEFORE THEY DO", f, width, f.Settings(), "left")
	bs.x = width
	return bs, nil
}

func (bs *MarqueeSaver) Inputs() map[string]SaverInput {
	return map[string]SaverInput{
		"font": {
			Default:     "slant",
			Description: "Font file to use",
		},
	}
}

func (bs *MarqueeSaver) SetInputs(inputs map[string]string) {
	bs.inputs = inputs
}

func (bs *MarqueeSaver) Initialize(opts screensaverOpts) error {
	bs.screen = opts.Screen
	bs.style = opts.Style

	return nil
}

func (bs *MarqueeSaver) Update() error {
	width, height := bs.screen.Size()
	y := height / 2
	bs.x--

	lines := strings.Split(bs.banner, "\n")
	if len(lines) == 0 {
		return nil
	}

	maxWidth := 0
	for _, line := range lines {
		if runewidth.StringWidth(line) > maxWidth {
			maxWidth = runewidth.StringWidth(line)
		}
	}

	// TODO fix this to work with banners
	// - ensure enough room for max width line
	// - draw str each line
	if bs.x+maxWidth < 0 {
		bs.x = width
	}

	for ix, line := range lines {
		drawStr(bs.screen, bs.x, y+ix, bs.style, line)
	}

	return nil
}
