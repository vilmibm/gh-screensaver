package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cli/safeexec"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/cobra"
)

type screensaverOpts struct {
	Screensaver string
	Repository  string
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
			}
			return runScreensaver(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Repository, "repo", "R", "", "Repository to contribute to")

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

func runScreensaver(opts screensaverOpts) error {
	// TODO reset terminal
	style := tcell.StyleDefault

	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err = s.Init(); err != nil {
		return err
	}
	s.SetStyle(style)

	// TODO do stuff

	return nil
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
