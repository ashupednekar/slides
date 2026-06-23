package main

import (
	"fmt"
	"io"
	"os"

	renderer "github.com/ashupednekar/slides/render/internal"
	"github.com/spf13/cobra"
)

type options struct {
	repoRoot string
	baseRoot string
	outRoot  string
	clean    bool
	notes    bool
}

func newRootCommand(out, errOut io.Writer) *cobra.Command {
	opts := options{
		baseRoot: "present-assets",
		outRoot:  "dist",
		clean:    true,
	}

	cmd := &cobra.Command{
		Use:          "render [deck-dir-or-slide...]",
		Short:        "Render Go present slide decks to static HTML",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := renderer.NewConfig(renderer.Options{
				RepoRoot: opts.repoRoot,
				BaseRoot: opts.baseRoot,
				OutRoot:  opts.outRoot,
				Clean:    opts.clean,
				Notes:    opts.notes,
			})
			if err != nil {
				return err
			}

			rendered, err := renderer.Run(cfg, args)
			if err != nil {
				return err
			}

			for _, deck := range rendered {
				fmt.Fprintf(out, "rendered %s -> %s\n", deck.Name, deck.OutDir)
			}
			return nil
		},
	}

	cmd.SetOut(out)
	cmd.SetErr(errOut)

	flags := cmd.Flags()
	flags.StringVar(&opts.repoRoot, "repo", "", "repository root; defaults to the current directory, or its parent when run from render/")
	flags.StringVar(&opts.baseRoot, "base", opts.baseRoot, "present template/static asset directory, relative to repo unless absolute")
	flags.StringVar(&opts.outRoot, "out", opts.outRoot, "output directory, relative to repo unless absolute")
	flags.BoolVar(&opts.clean, "clean", opts.clean, "remove the output directory before rendering")
	flags.BoolVar(&opts.notes, "notes", opts.notes, "render presenter notes support")

	return cmd
}

func execute() {
	if err := newRootCommand(os.Stdout, os.Stderr).Execute(); err != nil {
		os.Exit(1)
	}
}
