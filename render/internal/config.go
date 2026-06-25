package render

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultLargeAssetThreshold = 25 * 1024 * 1024

func NewConfig(opts Options) (Config, error) {
	repoRoot := opts.RepoRoot
	if repoRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return Config{}, err
		}
		repoRoot = cwd
		if filepath.Base(cwd) == "render" && exists(filepath.Join(filepath.Dir(cwd), "present-assets")) {
			repoRoot = filepath.Dir(cwd)
		}
	}

	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return Config{}, err
	}
	if !exists(filepath.Join(repoRoot, "present-assets")) {
		return Config{}, fmt.Errorf("repo root %q does not contain present-assets", repoRoot)
	}

	if opts.R2Bucket == "" {
		opts.R2Bucket = "slides"
	}
	if opts.R2PublicBaseURL == "" {
		opts.R2PublicBaseURL = os.Getenv("SLIDES_R2_PUBLIC_BASE_URL")
	}
	opts.R2PublicBaseURL = strings.TrimRight(opts.R2PublicBaseURL, "/")
	if opts.WranglerBin == "" {
		opts.WranglerBin = os.Getenv("WRANGLER_BIN")
	}
	if opts.WranglerBin == "" {
		opts.WranglerBin = "wrangler"
	}
	if opts.LargeAssetThreshold <= 0 {
		opts.LargeAssetThreshold = defaultLargeAssetThreshold
	}

	return Config{
		RepoRoot:            repoRoot,
		BaseRoot:            resolvePath(repoRoot, opts.BaseRoot),
		OutRoot:             resolvePath(repoRoot, opts.OutRoot),
		Clean:               opts.Clean,
		Notes:               opts.Notes,
		R2Bucket:            opts.R2Bucket,
		R2PublicBaseURL:     opts.R2PublicBaseURL,
		WranglerBin:         opts.WranglerBin,
		LargeAssetThreshold: opts.LargeAssetThreshold,
	}, nil
}

func prepareOutput(cfg Config) error {
	if samePath(cfg.RepoRoot, cfg.OutRoot) {
		return errors.New("refusing to use the repo root as the output directory")
	}
	if cfg.Clean {
		if err := validateCleanTarget(cfg); err != nil {
			return err
		}
		if err := os.RemoveAll(cfg.OutRoot); err != nil {
			return err
		}
	}
	return os.MkdirAll(cfg.OutRoot, 0o755)
}

func validateCleanTarget(cfg Config) error {
	outRoot, err := filepath.Abs(cfg.OutRoot)
	if err != nil {
		return err
	}
	if filepath.Dir(outRoot) == outRoot {
		return fmt.Errorf("refusing to clean filesystem root %q", outRoot)
	}

	if isWithin(cfg.RepoRoot, outRoot) {
		return nil
	}
	if isWithin(os.TempDir(), outRoot) {
		return nil
	}
	return fmt.Errorf("refusing to clean output outside the repo or temp dir: %s", outRoot)
}
