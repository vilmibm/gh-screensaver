package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vilmibm/gh-screensaver/savers"
	"github.com/vilmibm/gh-screensaver/savers/shared"
)

func runScreensaver(opts shared.ScreensaverOpts) error {
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

	// TODO ignore unused parameters
	providedInputs := map[string]string{}
	fs := pflag.FlagSet{}
	for inputName, input := range saver.Inputs() {
		fs.String(inputName, input.Default, input.Description)
	}
	err = fs.Parse(opts.SaverArgs)
	if err != nil {
		return fmt.Errorf("could not parse input args: %w", err)
	}
	for inputName := range saver.Inputs() {
		providedValue, _ := fs.GetString(inputName)
		providedInputs[inputName] = providedValue
	}

	err = saver.SetInputs(providedInputs)
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

		saver.Clear()
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
	opts := shared.ScreensaverOpts{}
	cmd := &cobra.Command{
		Use:   "screensaver",
		Short: "Watch a terminal saver animation",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.SaverArgs = args
			if opts.Repository == "" {
				repo, _ := resolveRepository()
				// Not erroring here; if a saver requires to know the repository it
				// will have to error itself if opts.Repository is ""
				opts.Repository = repo
			}
			opts.Savers = map[string]shared.SaverCreator{
				"marquee":   savers.NewMarqueeSaver,
				"fireworks": savers.NewFireworksSaver,
				"pipes":     savers.NewPipesSaver,
				"starfield": savers.NewStarfieldSaver,
				"pollock":   savers.NewPollockSaver,
				// TODO aquarium
				// TODO noise
				// TODO game of life
				// TODO issues/pr float by?
			}
			if opts.Screensaver == "" {
				opts.Screensaver = pickRandom(opts.Savers)
			}
			if opts.List {
				for _, k := range saverKeys(opts.Savers) {
					fmt.Println(k)
				}
				return nil
			}

			return runScreensaver(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository to contribute to")
	cmd.Flags().StringVarP(&opts.Screensaver, "saver", "s", "", "Screensaver to play")
	cmd.Flags().BoolVarP(&opts.List, "list", "l", false, "List available screensavers and exit")

	return cmd
}

func saverKeys(savers map[string]shared.SaverCreator) []string {
	keys := []string{}
	for k := range savers {
		keys = append(keys, k)
	}

	return keys
}

func pickRandom(savers map[string]shared.SaverCreator) string {
	rand.Seed(time.Now().UTC().UnixNano())
	keys := saverKeys(savers)
	ix := rand.Intn(len(keys))
	return keys[ix]
}

func main() {
	rc := rootCmd()

	if err := rc.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
