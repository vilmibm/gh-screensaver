package savers

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

// Want to launch fireworks from bottom of screen at a random x coord. Interval for shooting the firework with chance to fire per x coord chosen.
// Each firework should be its own updating entity that advances up a random number of y coords bounded by screen height, removing itself onces it's exploded.
// Fireworks should have a color pallete and then a particular explody animation.
// Not for MVP but would be cool to have "rare" fireworks that occur way less frequently.

type FireworksSaver struct {
	screen    tcell.Screen
	style     tcell.Style
	color     bool
	inputs    map[string]string
	fireworks []*firework
}

func NewFireworksSaver(opts shared.ScreensaverOpts) (shared.Screensaver, error) {
	fs := &FireworksSaver{}
	if err := fs.Initialize(opts); err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *FireworksSaver) Initialize(opts shared.ScreensaverOpts) error {
	fs.screen = opts.Screen
	fs.style = opts.Style

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func (fs *FireworksSaver) Inputs() map[string]shared.SaverInput {
	return map[string]shared.SaverInput{
		"color": {
			Default:     "full",
			Description: "whether to use full color or monochrome. Values: full, off",
		},
	}
}

func (fs *FireworksSaver) SetInputs(inputs map[string]string) error {
	fs.inputs = inputs
	if inputs["color"] == "full" {
		fs.color = true
	}
	return nil
}

func (fs *FireworksSaver) Update() error {
	next := []*firework{}
	for _, f := range fs.fireworks {
		f.Update()
		if !f.Done() {
			next = append(next, f)
		} else {
			continue
		}
		f.Draw()
	}
	fs.fireworks = next

	// TODO tweak as needed
	if rand.Intn(10) < 3 {
		fs.fireworks = append(fs.fireworks, newFirework(fs.screen, fs.style))
	}

	return nil
}

type spriteFrame string

type sprite struct {
	frames []spriteFrame
	frame  int
}

func (s *sprite) Advance() {
	s.frame++
	if s.frame == len(s.frames) {
		s.frame = 0
	}
}

type firework struct {
	Color1        tcell.Color
	Color2        tcell.Color
	TrailSprite   *sprite
	ExplodeSprite *sprite
	x             int
	y             int
	height        int
	screen        tcell.Screen
	style         tcell.Style
}

func parensTrail() *sprite {
	return &sprite{
		frames: []spriteFrame{"(", "|", ")"},
	}
}

func basicExplode() *sprite {
	return &sprite{
		frames: []spriteFrame{
			`
      *
`,
			`
     ( )
`,
			`
     ^
		( )
		 v
`,
			`
   * ^ *
	(     )
	 * v *
`,
			`
  \     /
   *   *
	(     )
	 *   *
	/     \
`,
			`
  \     /
   *   *
	       
	 *   *
	/     \
`,
			`
  \     /
        
	       
	      
	/     \
`,
			``,
		},
	}
}

func newFirework(screen tcell.Screen, style tcell.Style) *firework {
	width, height := screen.Size()
	f := &firework{
		screen: screen,
		style:  style,
		x:      rand.Intn(width),
		height: rand.Intn(height),
		// TODO randomize
		TrailSprite: parensTrail(),
		// TODO randomize
		ExplodeSprite: basicExplode(),
		// TODO randomize
		Color1: tcell.ColorBlue,
		// TODO randomize
		Color2: tcell.ColorLightBlue,
	}
	return f
}

func (f *firework) Update() {
	// TODO animate shootsprite
	// TODO increase height
	// TODO check height, prep to explode
	// TODO exlode
	// TODO remove
}

func (f *firework) Done() bool {
	return false
}

func (f *firework) Draw() {
}
