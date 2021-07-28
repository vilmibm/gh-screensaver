package savers

import (
	"embed"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/lukesampson/figlet/figletlib"
	"github.com/mattn/go-runewidth"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

//go:embed fonts/*
var fonts embed.FS

// TODO will likely share this
func drawStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

type MarqueeSaver struct {
	screen tcell.Screen
	style  tcell.Style
	x      int
	y      int
	banner string
	inputs map[string]string
}

func NewMarqueeSaver(opts shared.ScreensaverOpts) (shared.Screensaver, error) {
	bs := &MarqueeSaver{}
	if err := bs.Initialize(opts); err != nil {
		return nil, err
	}
	width, _ := bs.screen.Size()
	bs.x = width
	return bs, nil
}

func (bs *MarqueeSaver) Clear() {
	bs.screen.Clear()
}

func (bs *MarqueeSaver) Inputs() map[string]shared.SaverInput {
	// TODO list fonts in documentation
	return map[string]shared.SaverInput{
		"font": {
			Default:     "slant",
			Description: "Font file to use",
		},
		"message": {
			Default:     "text is cool",
			Description: "Message to display",
		},
	}
}

func (bs *MarqueeSaver) SetInputs(inputs map[string]string) error {
	bs.inputs = inputs
	data, err := fonts.ReadFile("fonts/" + bs.inputs["font"] + ".flf")
	if err != nil {
		return fmt.Errorf("no such font: %s: %w", bs.inputs["font"], err)
	}

	f, err := figletlib.ReadFontFromBytes(data)
	if err != nil {
		return err
	}
	width, _ := bs.screen.Size()
	bs.banner = figletlib.SprintMsg(bs.inputs["message"], f, width, f.Settings(), "left")
	return nil
}

func (bs *MarqueeSaver) Initialize(opts shared.ScreensaverOpts) error {
	bs.screen = opts.Screen
	bs.style = opts.Style

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func (bs *MarqueeSaver) Update() error {
	width, height := bs.screen.Size()
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

	if bs.x+maxWidth < 0 {
		bs.x = width
		bs.y = rand.Intn(height - len(lines))
	}

	for ix, line := range lines {
		drawStr(bs.screen, bs.x, bs.y+ix, bs.style, line)
	}

	return nil
}
