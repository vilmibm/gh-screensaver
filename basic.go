package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type BasicSaver struct {
	screen tcell.Screen
	style  tcell.Style
	x      int
	msg    string
}

func NewBasicSaver(opts screensaverOpts) (Screensaver, error) {
	bs := &BasicSaver{}
	if err := bs.Initialize(opts); err != nil {
		return nil, err
	}
	bs.msg = "HELLO WORLD"
	width, _ := bs.screen.Size()
	bs.x = width
	return bs, nil
}

func (bs *BasicSaver) Initialize(opts screensaverOpts) error {
	bs.screen = opts.Screen
	bs.style = opts.Style

	return nil
}

// TODO remember how to get terminal dimensions
// TODO track x pos in struct

func (bs *BasicSaver) Update() error {
	width, height := bs.screen.Size()
	y := height / 2
	bs.x--
	if bs.x+runewidth.StringWidth(bs.msg) < 0 {
		bs.x = width
	}

	drawStr(bs.screen, bs.x, y, bs.style, bs.msg)
	return nil
}
