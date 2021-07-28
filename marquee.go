package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type MarqueeSaver struct {
	screen tcell.Screen
	style  tcell.Style
	x      int
	msg    string
	inputs map[string]string
}

func NewMarqueeSaver(opts screensaverOpts) (Screensaver, error) {
	bs := &MarqueeSaver{}
	if err := bs.Initialize(opts); err != nil {
		return nil, err
	}
	bs.msg = "HELLO WORLD"
	width, _ := bs.screen.Size()
	bs.x = width
	return bs, nil
}

func (bs *MarqueeSaver) Inputs() map[string]SaverInput {
	return map[string]SaverInput{
		"script": {
			Default:     "script",
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
	if bs.x+runewidth.StringWidth(bs.msg) < 0 {
		bs.x = width
	}

	drawStr(bs.screen, bs.x, y, bs.style, bs.msg)
	return nil
}
