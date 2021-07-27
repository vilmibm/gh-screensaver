package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cli/safeexec"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

type screensaverOpts struct {
	Screensaver string
	Repository  string
	List        bool
	Style       tcell.Style
	Screen      tcell.Screen
	Savers      map[string]SaverCreator
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

func resolveRepository() (string, error) {
	sout, _, err := gh("repo", "view")
	if err != nil {
		return "", err
	}
	viewOut := strings.Split(sout.String(), "\n")[0]
	repo := strings.TrimSpace(strings.Split(viewOut, ":")[1])

	return repo, nil
}

func main() {
	rc := rootCmd()

	if err := rc.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

type Screensaver interface {
	Initialize(opts screensaverOpts) error
	Update() error
}

type BasicSaver struct {
	screen tcell.Screen
	style  tcell.Style
	frame  int
}

type SaverCreator func(screensaverOpts) (Screensaver, error)

func NewBasicSaver(opts screensaverOpts) (Screensaver, error) {
	bs := &BasicSaver{}
	if err := bs.Initialize(opts); err != nil {
		return nil, err
	}
	return bs, nil
}

func (bs *BasicSaver) Initialize(opts screensaverOpts) error {
	bs.screen = opts.Screen
	bs.style = opts.Style

	return nil
}

func (bs *BasicSaver) Update() error {
	bs.frame++
	x := 0
	y := 0
	drawStr(bs.screen, x, y, bs.style, "HELLO WORLD")
	return nil
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

// gh shells out to gh, returning STDOUT/STDERR and any error
func gh(args ...string) (sout, eout bytes.Buffer, err error) {
	ghBin, err := safeexec.LookPath("gh")
	if err != nil {
		err = fmt.Errorf("could not find gh. Is it installed? error: %w", err)
		return
	}

	cmd := exec.Command(ghBin, args...)
	cmd.Stderr = &eout
	cmd.Stdout = &sout

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("failed to run gh. error: %w, stderr: %s", err, eout.String())
		return
	}

	return
}
