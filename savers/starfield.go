package savers

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

// TODO try and adapt asciifield.c

const deg2rad = math.Pi / 180.0

type StarfieldSaver struct {
	screen tcell.Screen
	style  tcell.Style

	width  int
	height int

	n          float64
	f          float64
	fontAspect float64
	projAspect float64
	theta      float64

	projMatrix [16]float64

	speed    float64
	maxStars int

	stars []*star
}

func NewStarfieldSaver(opts shared.ScreensaverOpts) (shared.Screensaver, error) {
	s := &StarfieldSaver{}

	if err := s.Initialize(opts); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *StarfieldSaver) Initialize(opts shared.ScreensaverOpts) error {
	s.screen = opts.Screen
	s.style = opts.Style

	s.width, s.height = s.screen.Size()

	s.n = 0.1
	s.f = 10.0
	s.fontAspect = 0.5
	s.projAspect = float64(s.width) / float64(s.height) * s.fontAspect
	s.theta = 45 * deg2rad

	s.speed = 4.0
	s.maxStars = 300

	s.projMatrix = [16]float64{
		1.0 / math.Tan(s.theta*0.5) / s.projAspect,
		0, 0, 0,

		0,
		1.0 / math.Tan(s.theta*0.5),
		0, 0,

		0, 0,
		(s.f + s.n) / (s.f - s.n),
		-1,

		0, 0,
		(2 * s.n * s.f) / (s.f - s.n),
		0,
	}

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func (s *StarfieldSaver) Inputs() map[string]shared.SaverInput {
	// TODO support rainbow mode
	return map[string]shared.SaverInput{}
}

func (s *StarfieldSaver) SetInputs(inputs map[string]string) error {
	return nil
}

type star struct {
	vec []float64
}

func (s *star) Step(stepsize float64) {
	s.vec[2] += stepsize
}

func (s *star) Project(m [16]float64) [4]float64 {
	out := [4]float64{}
	out[0] = s.vec[0]*m[0] + s.vec[1]*m[4] + s.vec[2]*m[8] + s.vec[3]*m[12]
	out[1] = s.vec[0]*m[1] + s.vec[1]*m[5] + s.vec[2]*m[9] + s.vec[3]*m[13]
	out[2] = s.vec[0]*m[2] + s.vec[1]*m[6] + s.vec[2]*m[10] + s.vec[3]*m[14]
	out[3] = s.vec[0]*m[3] + s.vec[1]*m[7] + s.vec[2]*m[11] + s.vec[3]*m[15]

	// "dehomogenize"
	out[0] /= out[3]
	out[1] /= out[3]
	out[2] /= out[3]

	return out
}

func newStar(projAspect float64, f float64) *star {
	return &star{
		vec: []float64{
			(rand.Float64()*2 - 1) * 4 * projAspect,
			(rand.Float64()*2 - 1) * 4,
			float64(-f),
			1.0,
		},
	}
}

func (s *StarfieldSaver) Update() error {
	/*
		newStars := []*star{}
		for _, st := range s.stars {
			if !(st.vec[2] > -s.n || st.vec[2] < -s.f) {
				newStars = append(newStars, st)
			}
		}
		s.stars = newStars
	*/
	for len(s.stars) < s.maxStars {
		fmt.Fprintln(os.Stderr, len(s.stars))
		s.stars = append(s.stars, newStar(s.projAspect, s.f))
	}

	next := []*star{}

	for _, st := range s.stars {
		projected := st.Project(s.projMatrix)
		distance := st.vec[0]*st.vec[0] + st.vec[1]*st.vec[1] + st.vec[2]*st.vec[2]
		c := "#"
		if distance > 50 {
			c = "."
		} else if distance > 20 {
			c = "*"
		}
		x := int((projected[0] + 1) * 0.5 * float64(s.width))
		y := int((-projected[1] + 1) * 0.5 * float64(s.height))

		if x > 0 && x < s.width && y > 0 && y < s.width {
			// TODO set light or dark color based on distance
			drawStr(s.screen, x, y, s.style, c)
			// TODO WAG
			stepsize := s.speed * .04
			st.Step(stepsize)
			next = append(next, st)
		}
	}

	s.stars = next

	// TODO wtf with the depth calculation

	return nil
}
