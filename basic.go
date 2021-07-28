package main

import "github.com/gdamore/tcell/v2"

type BasicSaver struct {
	screen tcell.Screen
	style  tcell.Style
	frame  int
}

func NewBasicSaver(opts screensaverOpts) (Screensaver, error) {
	bs := &BasicSaver{}
	if err := bs.Initialize(opts); err != nil {
		return nil, err
	}
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
	bs.frame++
	x := 0
	y := 0
	drawStr(bs.screen, x, y, bs.style, "HELLO WORLD")
	return nil
}
