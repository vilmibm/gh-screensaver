package shared

import "github.com/gdamore/tcell/v2"

type SaverInput struct {
	Default     string
	Description string // TODO do i even want
}

type Screensaver interface {
	Initialize(opts ScreensaverOpts) error
	SetInputs(map[string]string) error
	Update() error
	Inputs() map[string]SaverInput
	Clear()
}

type SaverCreator func(ScreensaverOpts) (Screensaver, error)

type ScreensaverOpts struct {
	Screensaver string
	Repository  string
	List        bool
	Style       tcell.Style
	Screen      tcell.Screen
	Savers      map[string]SaverCreator
	SaverArgs   []string
}
