package savers

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

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

func (fs *FireworksSaver) Clear() {
	fs.screen.Clear()
}

func (fs *FireworksSaver) Initialize(opts shared.ScreensaverOpts) error {
	fs.screen = opts.Screen
	fs.style = opts.Style

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func (fs *FireworksSaver) Inputs() map[string]shared.SaverInput {
	// TODO eventually support truecolor
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
		f.Draw(fs.color)
	}
	fs.fireworks = next

	// TODO tweak as needed
	if rand.Intn(10) < 1 {
		fs.fireworks = append(fs.fireworks, newFirework(fs.screen, fs.style))
	}

	return nil
}

type sprite struct {
	frames []string
	Frame  int
	loop   bool
}

func (s *sprite) Advance() {
	s.Frame++
	if s.loop {
		if s.Frame == len(s.frames) {
			s.Frame = 0
		}
	}
}

func (s *sprite) CurrentFrame() string {
	return s.frames[s.Frame]
}

func (s *sprite) Done() bool {
	return s.Frame == len(s.frames)
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
	exploding     bool
	done          bool
}

func parensTrail() *sprite {
	return &sprite{
		loop: true,
		frames: []string{
			"(", "|", ")",
		},
	}
}

func sparkyTrail() *sprite {
	return &sprite{
		loop: true,
		frames: []string{
			"*", "x", ".",
		},
	}
}

func bangTrail() *sprite {
	return &sprite{
		loop: true,
		frames: []string{
			"i", "!", "|",
		},
	}
}

var trails = []func() *sprite{
	sparkyTrail,
	parensTrail,
	bangTrail,
}

func tinyBoomer() *sprite {
	return &sprite{
		frames: []string{
			`


      .


			`,
			`

      *

      
			`,
			`

     * *
    * * *
     * *
      
			`,
			`
    *   *
        
   *     *
        
    *   *
			`,
			`
         
            
            
           
           
			`,
		},
	}
}

func basicExplode() *sprite {
	return &sprite{
		frames: []string{
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
  /     \ `,
			`
  \     /
   *   *
          			 
   *   *
  /     \ `,
			`
  \     /
           
				  
         
  /     \ `,
			`
            
            
           
          
          `,
		},
	}
}

func sparkly() *sprite {
	return &sprite{
		frames: []string{
			`

     * *
      *  *
       *

			`,
			`
        *
       *
    * *  *
       *  *
    *
			`,
			`
  *    *   *
     *   *
    *  * *
   *   *    *
    *    *
			`,
			`
  * *      *
     *   *
       *    *
   *   *    *
 *  *    *   *
			`,
			`
    *    *  *
     *     
      *   * *
   *   * *   
 * *     *   *
			`,
			`
    *    *   
              
      *     *
               
 *       *   *
			`,
			`
                 
              
              
               
              
			`,
		},
	}
}

var explosions = []func() *sprite{
	tinyBoomer,
	basicExplode,
	sparkly,
}

var colors = []tcell.Color{
	tcell.ColorBlue,
	tcell.ColorCoral,
	tcell.ColorGoldenrod,
	tcell.ColorGray,
	tcell.ColorGreen,
	tcell.ColorPink,
	tcell.ColorSalmon,
	tcell.ColorSeaGreen,
	tcell.ColorDeepSkyBlue,
	tcell.ColorSlateGray,
	tcell.ColorSteelBlue,
	tcell.ColorYellow,
}

var lightColors = []tcell.Color{
	tcell.ColorLightBlue,
	tcell.ColorLightCoral,
	tcell.ColorLightGoldenrodYellow,
	tcell.ColorLightGray,
	tcell.ColorLightGreen,
	tcell.ColorLightPink,
	tcell.ColorLightSalmon,
	tcell.ColorLightSeaGreen,
	tcell.ColorLightSkyBlue,
	tcell.ColorLightSlateGray,
	tcell.ColorLightSteelBlue,
	tcell.ColorLightYellow,
}

func newFirework(screen tcell.Screen, style tcell.Style) *firework {
	width, height := screen.Size()
	colorIx := rand.Intn(len(colors))
	trailIx := rand.Intn(len(trails))
	explosionIx := rand.Intn(len(explosions))
	f := &firework{
		screen:        screen,
		style:         style,
		x:             rand.Intn(width-5) + 5,
		y:             height,
		height:        rand.Intn(height - 8),
		TrailSprite:   trails[trailIx](),
		ExplodeSprite: explosions[explosionIx](),
		Color1:        colors[colorIx],
		Color2:        lightColors[colorIx],
	}
	return f
}

func (f *firework) Update() {
	if f.y == f.height {
		f.exploding = true
	} else {
		f.y--
	}
}

// TODO affordance for setting animation interval from within screensaver

func (f *firework) Done() bool {
	return f.done
}

func (f *firework) Draw(useColor bool) {
	if f.exploding {
		if f.ExplodeSprite.Done() {
			f.done = true
			return
		}

		color := f.Color1
		colorChoice := f.ExplodeSprite.Frame % 2
		if colorChoice == 1 {
			color = f.Color2
		}

		lines := strings.Split(f.ExplodeSprite.CurrentFrame(), "\n")
		for ix, line := range lines {
			s := f.style
			if useColor {
				s = f.style.Foreground(color)
			}
			drawStr(f.screen, f.x-2, f.y+ix-2, s, line)
		}

		f.ExplodeSprite.Advance()

		return
	}

	color := f.Color1
	colorChoice := f.TrailSprite.Frame % 2
	if colorChoice == 1 {
		color = f.Color2
	}
	s := f.style
	if useColor {
		s = f.style.Foreground(color)
	}
	drawStr(f.screen, f.x, f.y, s, f.TrailSprite.CurrentFrame())
	f.TrailSprite.Advance()
}
