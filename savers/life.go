package savers

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

var seeds = []string{"dragon", "gun", "noise", "r", "pulsar", "glider"}

var aliveColorsBlue = []tcell.Color{
	tcell.ColorDeepSkyBlue,
	tcell.ColorBlue,
}

var aliveColorsGreen = []tcell.Color{
	tcell.ColorGreen,
	tcell.ColorGoldenrod,
}

var aliveColorsRed = []tcell.Color{
	tcell.ColorYellow,
	tcell.ColorRed,
}

var aliveColorsWhite = []tcell.Color{
	tcell.ColorWhite,
	tcell.ColorGold,
}

type LifeSaver struct {
	screen tcell.Screen
	style  tcell.Style

	width  int
	height int

	useColor   bool
	colors     []tcell.Color
	aliveCells [][]int
}

func NewLifeSaver(opts shared.ScreensaverOpts) (shared.Screensaver, error) {
	lf := &LifeSaver{}
	if err := lf.Initialize(opts); err != nil {
		return nil, err
	}

	lf.aliveCells = make([][]int, lf.width)
	for i := range lf.aliveCells {
		lf.aliveCells[i] = make([]int, lf.height)
	}

	return lf, nil
}

func (lf *LifeSaver) Clear() {
	lf.screen.Clear()
}

func (lf *LifeSaver) Initialize(opts shared.ScreensaverOpts) error {
	lf.screen = opts.Screen
	lf.style = opts.Style
	lf.width, lf.height = lf.screen.Size()

	rand.Seed(time.Now().UTC().UnixNano())

	return nil
}

func (lf *LifeSaver) Inputs() map[string]shared.SaverInput {
	return map[string]shared.SaverInput{
		"seed": {
			Default:     "rand",
			Description: "seed state. Values: dragon, gun, R, pulsar, noise",
		},
		"color": {
			Default:     "full",
			Description: "whether to use full color or monochrome. Values: full, off",
		},
	}
}

func (lf *LifeSaver) SetInputs(inputs map[string]string) error {
	lf.useColor = inputs["color"] == "full"
	seed := strings.ToLower(inputs["seed"])
	if seed == "rand" {
		idx := rand.Intn(len(seeds))
		seed = seeds[idx]
	}

	// default to noise if terminal is too small
	defaultSeed := "noise"
	switch seed {
	case "gun":
		if lf.width < 50 || lf.height < 42 {
			seed = defaultSeed
		}
	case "pulsar":
		if lf.width < 40 || lf.height < 40 {
			seed = defaultSeed
		}
	case "dragon":
		if lf.width < 26 || lf.height < 30 {
			seed = defaultSeed
		}
	case "r":
		if lf.width < 15 || lf.height < 15 {
			seed = defaultSeed
		}
	case "glider":
		if lf.width < 10 || lf.height < 10 {
			seed = defaultSeed
		}
	default:
		seed = defaultSeed
	}
	err := lf.initState(seed)
	if err != nil {
		return err
	}
	return nil
}

func (lf *LifeSaver) initState(seed string) error {
	tX := lf.width / 3
	tY := lf.height / 3
	hX := lf.width / 2
	hY := lf.height / 2
	switch seed {
	case "pulsar":
		// some oscillators
		for i := 0; i < 37; i++ {
			for j := 0; j < 37; j++ {
				lf.aliveCells[j+hX-15][i+hY-17] = pulsar[i][j]
			}
		}

		if hX/2 > 10 {
			for i := 0; i < 3; i++ {
				for j := 0; j < 8; j++ {
					lf.aliveCells[i+hX/2][j+hY] = pentadec[i][j]
					lf.aliveCells[i+3*hX/2][j+hY] = pentadec[i][j]
				}
			}
		}
		lf.colors = aliveColorsRed
	case "gun":
		// colliding glider guns!!!
		for i := 0; i < 36; i++ {
			for j := 0; j < 9; j++ {
				lf.aliveCells[5+j][5+i] = glidergun[i][j]
				lf.aliveCells[lf.width-10+j][5+i] = glidergun[i][8-j]
			}
		}
		lf.colors = aliveColorsBlue
	case "dragon":
		// there be dragons
		for n := 5; n+20 < lf.width; n += 27 {
			for i := 0; i < 18; i++ {
				for j := 0; j < 29; j++ {
					if n%2 == 0 {
						lf.aliveCells[i+n][(j+2*n*tY)%lf.height] = dragon[i][j]
					} else {
						lf.aliveCells[i+n][(j+n*tY)%lf.height] = dragon[i][28-j]
					}
				}
			}
		}

		lf.colors = aliveColorsGreen
	case "r":
		//R-pentominos - chaotic
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				lf.aliveCells[tX+i][hY+j] = rPentomino[j][i]
				lf.aliveCells[2*tX+i][hY+j] = rPentomino[i][j]
			}
		}
		lf.colors = aliveColorsRed
	case "noise":
		//random noise seed
		for i := 0; i < lf.width; i++ {
			for j := 0; j < lf.height; j++ {
				if rand.Intn(10) < 2 {
					lf.aliveCells[i][j] = 1
				} else {
					lf.aliveCells[i][j] = 0
				}
			}
		}
		lf.colors = aliveColorsBlue
	case "glider":
		//glider fleet
		for k := 2; k+3 < lf.width; k += 15 {
			h := rand.Intn(lf.height)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					lf.aliveCells[k+i][(j+h+lf.height)%lf.height] = glider[i][j]
				}
			}
		}
		lf.colors = aliveColorsWhite
	default:
		return errors.New("Error initiliazing seed")
	}

	return nil
}

func (lf *LifeSaver) isOnGrid(x, y int) bool {
	return x >= 0 && y >= 0 && x < lf.width && y < lf.height
}

func (lf *LifeSaver) Update() error {
	for i := 0; i < lf.width; i++ {
		for j := 0; j < lf.height; j++ {
			if lf.aliveCells[i][j] > 0 {
				// cell is alive
				for k := 0; k < 8; k++ {
					ni := (i + nbrX[k] + lf.width) % lf.width
					nj := (j + nbrY[k] + lf.height) % lf.height
					if lf.aliveCells[ni][nj] > 0 {
						lf.aliveCells[i][j]++
					}
				}
			} else {
				// cell is dead
				for k := 0; k < 8; k++ {
					ni := (i + nbrX[k] + lf.width) % lf.width
					nj := (j + nbrY[k] + lf.height) % lf.height
					if lf.aliveCells[ni][nj] > 0 {
						lf.aliveCells[i][j]--
					}
				}
			}
		}
	}

	// next generation
	for i := 0; i < lf.width; i++ {
		for j := 0; j < lf.height; j++ {
			n := lf.aliveCells[i][j]
			switch {
			case n == -3 || n == 3:
				if lf.useColor {
					drawStr(lf.screen, i, j, lf.style.Foreground(lf.colors[0]), "*")
				} else {
					drawStr(lf.screen, i, j, lf.style, "*")
				}
				lf.aliveCells[i][j] = 1
			case n == 4:
				if lf.useColor {
					drawStr(lf.screen, i, j, lf.style.Foreground(lf.colors[1]), "#")
				} else {
					drawStr(lf.screen, i, j, lf.style, "*")
				}
				lf.aliveCells[i][j] = 1
			default:
				lf.aliveCells[i][j] = 0
				drawStr(lf.screen, i, j, lf.style, " ")
			}
		}
	}

	return nil
}

var nbrX = []int{1, -1, 0, 1, -1, 0, 1, -1}
var nbrY = []int{0, 0, -1, -1, -1, 1, 1, 1}

var rPentomino = [3][3]int{
	{0, 1, 0},
	{1, 1, 0},
	{0, 1, 1}}

var glider = [3][3]int{
	{1, 0, 0},
	{1, 0, 1},
	{1, 1, 0}}

var pentadec = [3][8]int{
	{1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 1, 1, 1, 1, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1}}

var glidergun = [36][9]int{
	{0, 0, 0, 0, 1, 1, 0, 0, 0},
	{0, 0, 0, 0, 1, 1, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 1, 0},
	{0, 0, 1, 0, 0, 0, 0, 0, 1},
	{0, 0, 1, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 1, 0},
	{0, 0, 0, 0, 1, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 1, 0, 0, 0, 0},
	{0, 0, 1, 1, 1, 0, 0, 0, 0},
	{0, 1, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{1, 1, 0, 0, 0, 1, 1, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 0, 0, 0, 0, 0}}

var dragon = [18][29]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	{1, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1},
	{1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0},
	{1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0},
	{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0},
	{1, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0},
	{1, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1},
	{0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}

var pulsar = [37][37]int{
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
	{0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
	{1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1},
	{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
	{0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}}
