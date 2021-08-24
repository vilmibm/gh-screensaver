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

func (ps *PipesSaver) Clear() {
	ps.screen.Clear()
}

func (ps *PipesSaver) Initialize(opts shared.ScreensaverOpts) error {
	ps.screen = opts.Screen
	ps.style = opts.Style

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func (ps *PipesSaver) Inputs() map[string]shared.SaverInput {
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

type dir int

const (
	up    dir = 0
	right dir = 1
	down  dir = 2
	left  dir = 3
	cw    dir = 4
	ccw   dir = 5
)

type relDir int

type pipe struct {
	dir    dir
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

func turn(pd dir, d dir) dir {
	if d == cw {
		switch pd {
		case up:
			return right
		case right:
			return down
		case down:
			return left
		case left:
			return up
		}
	} else if d == ccw {
		switch pd {
		case up:
			return left
		case right:
			return up
		case down:
			return right
		case left:
			return down
		}
	} else {
		panic("bad dir")
	}

	return pd
}

func (p *pipe) Next() {
	last := p.coords[len(p.coords)-1]
	// 80% continue in current dir
	// 10% turn right
	// 10% turn left
	score := rand.Intn(10)

	// continue
	if score == 8 {
		p.dir = turn(p.dir, cw)
	} else if score == 9 {
		p.dir = turn(p.dir, ccw)
	}

	var next coord
	switch p.dir {
	case up:
		next = coord{last.x, last.y - 1}
	case right:
		next = coord{last.x + 1, last.y}
	case down:
		next = coord{last.x, last.y + 1}
	case left:
		next = coord{last.x - 1, last.y}
	}

	if next.x < 0 || next.y > p.height || next.y < 0 || next.y > p.width {
		// just hang out until we get a legal move
		return
	}

	p.coords = append(p.coords, next)
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
		var pipe *pipe
		for pipe == nil {
			pipe = newPipe(width, height)
			for _, p := range ps.pipes {
				if p.coords[0].x == pipe.coords[0].x && p.coords[0].y == pipe.coords[0].y {
					pipe = nil
				} else {
					break
				}
			}
		}

		ps.pipes = append(ps.pipes, pipe)
	}

	for _, p := range ps.pipes {
		p.Next()
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
