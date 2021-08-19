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
	pipes  []*pipe
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
	width  int
	height int
}

func newPipe(width, height int) *pipe {
	p := &pipe{
		color:  randColor(),
		width:  width,
		height: height,
	}
	x := 0
	y := 0
	switch rand.Intn(4) {
	case 0: // top
		x = rand.Intn(width)
	case 1: // right
		x = width
		y = rand.Intn(height)
	case 2: // bottom
		y = height
		x = rand.Intn(width)
	case 3: // left
		y = rand.Intn(height)
	}

	p.coords = []coord{{x, y}}

	return p
}

// TODO don't spawn multiple pipes at the same point

// problem: i want to have more negative space; i'd like the pipes to not grow so close to themselves.

// what if instead of storing absolute points I stored a series of vectors? the valid moves shouldn't be one of three surrounding points, but a choice between:

// - continue
// - turn left
// - turn right

// while forbidding two turns in a row

type move int

const (
	straight  move = 0
	turnLeft  move = 1
	turnRight move = 2
)

func (p *pipe) ValidMoves() []coord {
	// current
	// ###
	//   #*
	//   ##*
	//    *

	// two straights
	// ###
	//   #
	//   ##*
	//

	// TODO want to do two straighs before a turn; how to compute that?
	// - either x or y unchanged for last and penul
	last := p.coords[len(p.coords)-1]
	penul := coord{-1, -1}
	if len(p.coords) > 1 {
		penul = p.coords[len(p.coords)-2]
	}
	all := []coord{}
	if last.y > 0 {
		all = append(all, coord{last.x, last.y - 1})
	}
	if last.x < p.width {
		all = append(all, coord{last.x + 1, last.y})
	}
	if last.y < p.height {
		all = append(all, coord{last.x, last.y + 1})
	}
	if last.x > 0 {
		all = append(all, coord{last.x - 1, last.y})
	}
	out := []coord{}
	for _, c := range all {
		if c.x != penul.x && c.y != penul.y {
			// Is it not the penultimate point? Want to prevent backtracking
			out = append(out, c)
		}
	}
	return out
}

func (p *pipe) AddCoord(c coord) {
	p.coords = append(p.coords, c)
}

func randColor() tcell.Color {
	r := rand.Int31n(255)
	g := rand.Int31n(255)
	b := rand.Int31n(255)

	return tcell.NewRGBColor(r, g, b)
}

func (ps *PipesSaver) Update() error {
	width, height := ps.screen.Size()
	if rand.Intn(10) < 1 {
		pipe := newPipe(width, height)
		ps.pipes = append(ps.pipes, pipe)
	}

	for _, p := range ps.pipes {
		moves := p.ValidMoves()
		ix := rand.Intn(len(moves))
		p.AddCoord(moves[ix])
	}

	for _, p := range ps.pipes {
		for _, c := range p.coords {
			s := ps.style
			if ps.color {
				s = s.Foreground(p.color)
			}
			drawStr(ps.screen, c.x, c.y, s, "#")
		}
	}

	// TODO the OG pipes clears itself at some interval. I think it will take far
	// more time for us to fill up a screen, so initially I think i'll just let
	// it fill up.
	return nil
}
