package savers

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

type PollockSaver struct {
	screen tcell.Screen
	style  tcell.Style

	width  int
	height int

	splats    []*splat
	maxSplats int
}

func NewPollockSaver(opts shared.ScreensaverOpts) (shared.Screensaver, error) {
	p := &PollockSaver{}

	if err := p.Initialize(opts); err != nil {
		return nil, err
	}

	return p, nil
}

// pick random spot
// pick random color
// create a bounding box
// for each square in the box, splat with a chance

func (p *PollockSaver) Initialize(opts shared.ScreensaverOpts) error {
	p.screen = opts.Screen
	p.style = opts.Style
	p.width, p.height = p.screen.Size()

	p.maxSplats = 1000

	rand.Seed(time.Now().UTC().UnixNano())
	return nil
}

func (p *PollockSaver) Inputs() map[string]shared.SaverInput {
	return map[string]shared.SaverInput{}
}

func (p *PollockSaver) SetInputs(inputs map[string]string) error {
	return nil
}

type splat struct {
	x      int
	y      int
	width  int
	height int
	color  tcell.Color
	cells  [][]string
}

func newSplat(width, height int) *splat {
	s := &splat{
		x:      rand.Intn(width),
		y:      rand.Intn(height),
		width:  rand.Intn(width / 2),
		height: rand.Intn(width / 2),
		color:  randColor(),
		cells:  [][]string{},
	}
	for ix := 0; ix < s.width; ix++ {
		s.cells = append(s.cells, []string{})
		for iy := 0; iy < s.height; iy++ {
			splat := rand.Intn(10)
			if splat < 2 {
				c := "#"
				cChance := rand.Intn(10)
				if cChance > 8 {
					c = "+"
				} else if cChance > 4 {
					c = "*"
				}

				s.cells[ix] = append(s.cells[ix], c)
			} else {
				s.cells[ix] = append(s.cells[ix], " ")
			}
		}
	}

	return s
}

func (p *PollockSaver) Update() error {
	if rand.Intn(10) > 5 {
		return nil
	}

	p.splats = append(p.splats, newSplat(p.width, p.height))

	if len(p.splats) > p.maxSplats {
		p.splats = p.splats[1:]
	}

	for _, splat := range p.splats {
		ix := 0
		for _, splatRow := range splat.cells {
			iy := 0
			for _, splatCell := range splatRow {
				if splatCell != " " {
					drawStr(p.screen, ix+splat.x, iy+splat.y, p.style.Foreground(splat.color), splatCell)
				}
				iy++
			}
			ix++
		}
	}

	return nil
}
