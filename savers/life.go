package savers

import (
	"github.com/gdamore/tcell/v2"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

var nbrX = []int{1, -1, 0, 1, -1, 0, 1, -1}
var nbrY = []int{0, 0, -1, -1, -1, 1, 1, 1}

var aliveColorsBlue = []tcell.Color{
	tcell.ColorDeepSkyBlue,
	tcell.ColorBlue,
}

type LifeSaver struct {
	screen tcell.Screen
	style  tcell.Style

	width  int
	height int

	periodic   bool
	colors     []tcell.Color
	aliveCells [][]int
}

func NewLifeSaver(opts shared.ScreensaverOpts) (shared.Screensaver, error) {
	p := &LifeSaver{colors: aliveColorsBlue}

	if err := p.Initialize(opts); err != nil {
		return nil, err
	}

	p.aliveCells = make([][]int, p.width)
	for i := range p.aliveCells {
		p.aliveCells[i] = make([]int, p.height)
	}

	//for i := 0; i < p.width; i++ {
	//for j := 0; j < p.height; j++ {
	//r := rand.Intn(15)
	//if r < 2 {
	//p.aliveCells[i][j] = 1
	//}
	//}
	//}

	// two R-pentominos

	//tX := p.width / 3
	//hY := p.height / 2
	//for i := 0; i < 3; i++ {
	//for j := 0; j < 3; j++ {
	//p.aliveCells[tX+i][hY+j] = rPentomino[j][i]
	//p.aliveCells[2*tX+i][hY+j] = rPentomino[i][j]
	//}
	//}

	// colliding glider guns!!!
	for i := 0; i < 36; i++ {
		for j := 0; j < 9; j++ {
			p.aliveCells[10+j][5+i] = glidergun[i][j]
			p.aliveCells[p.width-20+j][5+i] = glidergun[i][9-j-1]
		}
	}

	return p, nil
}

func (p *LifeSaver) Clear() {}

func (p *LifeSaver) Initialize(opts shared.ScreensaverOpts) error {
	p.screen = opts.Screen
	p.style = opts.Style
	p.width, p.height = p.screen.Size()

	return nil
}

func (p *LifeSaver) Inputs() map[string]shared.SaverInput {
	return map[string]shared.SaverInput{}
}

func (p *LifeSaver) SetInputs(inputs map[string]string) error {
	return nil
}

func (p *LifeSaver) isOnGrid(x, y int) bool {
	return x >= 0 && y >= 0 && x < p.width && y < p.height
}

func (p *LifeSaver) Update() error {
	for i := 0; i < p.width; i++ {
		for j := 0; j < p.height; j++ {
			if p.aliveCells[i][j] > 0 {
				// cell is alive
				for k := 0; k < 8; k++ {
					ni := (i + nbrX[k] + p.width) % p.width
					nj := (j + nbrY[k] + p.height) % p.height
					if p.aliveCells[ni][nj] > 0 {
						p.aliveCells[i][j]++
					}
				}
			} else {
				// cell is dead
				for k := 0; k < 8; k++ {
					ni := (i + nbrX[k] + p.width) % p.width
					nj := (j + nbrY[k] + p.height) % p.height
					if p.aliveCells[ni][nj] > 0 {
						p.aliveCells[i][j]--
					}
				}
			}
		}
	}

	// next generation
	for i := 0; i < p.width; i++ {
		for j := 0; j < p.height; j++ {
			n := p.aliveCells[i][j]
			switch {
			case n == -3 || n == 3:
				drawStr(p.screen, i, j, p.style.Foreground(p.colors[0]), "*")
				p.aliveCells[i][j] = 1
			case n == 4:
				drawStr(p.screen, i, j, p.style.Foreground(p.colors[1]), "#")
				p.aliveCells[i][j] = 1
			default:
				p.aliveCells[i][j] = 0
				drawStr(p.screen, i, j, p.style, " ")
			}
		}
	}

	return nil
}

var rPentomino = [3][3]int{
	{0, 1, 0},
	{1, 1, 0},
	{0, 1, 1}}

var glider = [3][3]int{
	{1, 0, 0},
	{1, 0, 1},
	{1, 1, 0}}

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
