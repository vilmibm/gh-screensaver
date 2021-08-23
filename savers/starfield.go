package savers

import (
	"math"
	"math/rand"
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

	maxStars int
	interval time.Duration

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

	interval, _ := time.ParseDuration("1s")
	s.interval = interval
	s.maxStars = 10

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

func (s *star) Project(m [16]float64) []float64 {
	original := []float64{}
	for _, f := range s.vec {
		original = append(original, f)
	}
	s.vec[0] = s.vec[0]*m[0] + s.vec[1]*m[4] + s.vec[2]*m[8] + s.vec[3]*m[12]
	s.vec[1] = s.vec[0]*m[1] + s.vec[1]*m[5] + s.vec[2]*m[9] + s.vec[3]*m[13]
	s.vec[2] = s.vec[0]*m[2] + s.vec[1]*m[6] + s.vec[2]*m[10] + s.vec[3]*m[14]
	s.vec[3] = s.vec[0]*m[3] + s.vec[1]*m[7] + s.vec[2]*m[11] + s.vec[3]*m[15]

	// "dehomogenize"
	s.vec[0] = s.vec[0] / s.vec[3]
	s.vec[1] = s.vec[1] / s.vec[3]
	s.vec[2] = s.vec[2] / s.vec[3]

	return original
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
	for len(s.stars) < s.maxStars {
		s.stars = append(s.stars, newStar(s.projAspect, s.f))
	}

	starsNext := []*star{}

	for _, st := range s.stars {
		v := st.Project(s.projMatrix)
		distance := v[0]*v[0] + v[1]*v[1] + v[2]*v[2]
		c := "#"
		if distance > 50 {
			c = "."
		} else if distance > 20 {
			c = "*"
		}
		x := int((st.vec[0] + 1) * 0.5 * float64(s.width))
		y := int((st.vec[1] + 1) * 0.5 * float64(s.height))
		if x > 0 && x < s.width && y > 0 && y < s.height {
			// TODO set light or dark color based on distance
			drawStr(s.screen, x, y, s.style, c)
			starsNext = append(starsNext, st)
		}
		time.Sleep(1)
	}

	s.stars = starsNext

	// TODO wtf with the depth calculation

	return nil
}
