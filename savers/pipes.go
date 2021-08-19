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

func (p *pipe) ValidMoves() []coord {
	last := p.coords[len(p.coords)-1]
	penul := coord{-1, -1}
	if len(p.coords) > 1 {
		penul = p.coords[len(p.coords)-2]
	}
	out := []coord{}
	if last.y > 0 {
		y := last.y - 1
		x := last.x
		if x != penul.x && y != penul.y {
			out = append(out, coord{last.x, last.y - 1})
		}
	}
	if last.x < p.width {
		x := last.x + 1
		y := last.y
		if x != penul.x && y != penul.y {
			out = append(out, coord{last.x + 1, last.y})
		}
	}
	if last.y < p.height {
		x := last.x
		y := last.y + 1
		if x != penul.x && y != penul.y {
			out = append(out, coord{last.x, last.y + 1})
		}
	}
	if last.x > 0 {
		x := last.x - 1
		y := last.y
		if x != penul.x && y != penul.y {
			out = append(out, coord{last.x - 1, last.y})
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
			// TODO determine if it's a change in direction and use "*"
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
