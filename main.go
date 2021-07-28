package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

type Screensaver interface {
	Initialize(opts screensaverOpts) error
	Update() error
}

type SaverCreator func(screensaverOpts) (Screensaver, error)

type screensaverOpts struct {
	Screensaver string
	Repository  string
	List        bool
	Style       tcell.Style
	Screen      tcell.Screen
	Savers      map[string]SaverCreator
}

func drawStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func runScreensaver(opts screensaverOpts) error {
	style := tcell.StyleDefault
	opts.Style = style

	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err = screen.Init(); err != nil {
		return err
	}
	screen.SetStyle(style)

	opts.Screen = screen

	saver, err := opts.Savers[opts.Screensaver](opts)
	if err != nil {
		return err
	}

	quit := make(chan struct{})
	go func() {
		for {
			ev := screen.PollEvent()
			switch ev.(type) {
			case *tcell.EventKey:
				close(quit)
				return
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()

	var saverErr error
loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 100):
		}

		screen.Clear()
		if err := saver.Update(); err != nil {
			saverErr = err
			break loop
		}
		screen.Show()
	}

	screen.Fini()

	return saverErr
}

func rootCmd() *cobra.Command {
	opts := screensaverOpts{}
	cmd := &cobra.Command{
		Use:   "screensaver",
		Short: "Watch a terminal saver animation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Repository == "" {
				repo, err := resolveRepository()
				if err != nil {
					return err
				}
				opts.Repository = repo
				opts.Savers = map[string]SaverCreator{
					"basic": NewBasicSaver,
				}
				// TODO respect -l
				// TODO support random by default
			}
			return runScreensaver(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository to contribute to")
	cmd.Flags().StringVarP(&opts.Screensaver, "saver", "s", "basic", "Screensaver to play")
	cmd.Flags().BoolVarP(&opts.List, "list", "l", false, "List available screensavers and exit")

	return cmd
}

func main() {
	rc := rootCmd()

	if err := rc.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
