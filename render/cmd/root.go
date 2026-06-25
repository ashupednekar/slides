package cmd

import (
	"fmt"
	"io"

	renderer "github.com/ashupednekar/slides/render/internal"
	"github.com/spf13/cobra"
)

type options struct {
	repoRoot            string
	baseRoot            string
	outRoot             string
	clean               bool
	notes               bool
	r2Bucket            string
	r2PublicBaseURL     string
	wranglerBin         string
	largeAssetThreshold int64
}

func NewRootCommand(out, errOut io.Writer) *cobra.Command {
	opts := options{
		baseRoot:            "present-assets",
		outRoot:             "dist",
		clean:               true,
		r2Bucket:            "slides",
		largeAssetThreshold: 25,
	}

	cmd := &cobra.Command{
		Use:          "render [deck-dir-or-slide...]",
		Short:        "Render Go present slide decks to static HTML",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := renderer.NewConfig(renderer.Options{
				RepoRoot:            opts.repoRoot,
				BaseRoot:            opts.baseRoot,
				OutRoot:             opts.outRoot,
				Clean:               opts.clean,
				Notes:               opts.notes,
				R2Bucket:            opts.r2Bucket,
				R2PublicBaseURL:     opts.r2PublicBaseURL,
				WranglerBin:         opts.wranglerBin,
				LargeAssetThreshold: opts.largeAssetThreshold * 1024 * 1024,
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
	flags.StringVar(&opts.r2Bucket, "r2-bucket", opts.r2Bucket, "R2 bucket for assets larger than the threshold")
	flags.StringVar(&opts.r2PublicBaseURL, "r2-public-base-url", "", "public URL prefix for R2 objects; defaults to SLIDES_R2_PUBLIC_BASE_URL")
	flags.StringVar(&opts.wranglerBin, "wrangler-bin", opts.wranglerBin, "Wrangler executable used for R2 object checks and uploads; defaults to WRANGLER_BIN or wrangler")
	flags.Int64Var(&opts.largeAssetThreshold, "large-asset-threshold-mb", opts.largeAssetThreshold, "offload static assets larger than this many MiB")

	return cmd
}
