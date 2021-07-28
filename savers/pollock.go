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

func (p *PollockSaver) Clear() {}

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

type paintCell struct {
	x    int
	y    int
	char string
}

type splat struct {
	color tcell.Color
	cells []paintCell
}

func newSplat(width, height int) *splat {
	s := &splat{
		color: randColor(),
	}

	c := "#"
	cChance := rand.Intn(10)
	if cChance > 8 {
		c = "+"
	} else if cChance > 4 {
		c = "*"
	}

	cell := paintCell{
		x:    rand.Intn(width),
		y:    rand.Intn(height),
		char: c,
	}

	s.cells = []paintCell{cell}

	return s
}

func (s *splat) Spread(width, height int) {
	last := s.cells[len(s.cells)-1]

	all := []paintCell{}

	if last.y > 0 {
		all = append(all, paintCell{x: last.x, y: last.y - 1})
	}
	if last.x < width {
		all = append(all, paintCell{x: last.x + 1, y: last.y})
	}
	if last.y < height {
		all = append(all, paintCell{x: last.x, y: last.y + 1})
	}
	if last.x > 0 {
		all = append(all, paintCell{x: last.x - 1, y: last.y})
	}

	if len(all) == 0 {
		return
	}

	next := all[rand.Intn(len(all))]
	next.char = "#"
	cChance := rand.Intn(10)
	if cChance > 8 {
		next.char = "+"
	} else if cChance > 4 {
		next.char = "*"
	}
	s.cells = append(s.cells, next)
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
		if rand.Intn(10) < 5 {
			splat.Spread(p.width, p.height)
		}

		for _, cell := range splat.cells {
			drawStr(p.screen, cell.x, cell.y, p.style.Foreground(splat.color), cell.char)
		}
	}

	return nil
}
