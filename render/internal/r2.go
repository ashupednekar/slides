package render

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"mime"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const placeholderName = "placeholder.png"

var placeholderPNG = mustDecodeBase64("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/p9sAAAAASUVORK5CYII=")

type largeAsset struct {
	rel     string
	siteRel string
	path    string
	size    int64
}

type assetReplacement struct {
	rel string
	url string
}

func offloadLargeAssets(cfg Config, outDir string) error {
	assets, err := largeAssets(cfg.OutRoot, outDir, cfg.LargeAssetThreshold)
	if err != nil {
		return err
	}
	if len(assets) == 0 {
		return nil
	}

	uploader := newR2Uploader(cfg)
	replacements := make([]assetReplacement, 0, len(assets))
	usedPlaceholder := false

	for _, asset := range assets {
		targetURL, ok := uploader.publicURL(asset)
		if !ok {
			targetURL = placeholderName
			usedPlaceholder = true
		}
		replacements = append(replacements, assetReplacement{rel: asset.rel, url: targetURL})
	}

	if usedPlaceholder {
		if err := writePlaceholder(outDir); err != nil {
			return err
		}
	}
	if err := rewriteAssetReferences(outDir, replacements); err != nil {
		return err
	}
	for _, asset := range assets {
		if err := os.Remove(asset.path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func largeAssets(siteRoot, outDir string, threshold int64) ([]largeAsset, error) {
	var assets []largeAsset
	err := filepath.WalkDir(outDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Size() <= threshold {
			return nil
		}
		rel, err := filepath.Rel(outDir, path)
		if err != nil {
			return err
		}
		siteRel, err := filepath.Rel(siteRoot, path)
		if err != nil {
			return err
		}
		assets = append(assets, largeAsset{
			rel:     filepath.ToSlash(rel),
			siteRel: filepath.ToSlash(siteRel),
			path:    path,
			size:    info.Size(),
		})
		return nil
	})
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].rel < assets[j].rel
	})
	return assets, err
}

type r2Uploader struct {
	cfg      Config
	disabled bool
}

func newR2Uploader(cfg Config) *r2Uploader {
	_, err := exec.LookPath(cfg.WranglerBin)
	return &r2Uploader{
		cfg:      cfg,
		disabled: err != nil || cfg.R2PublicBaseURL == "",
	}
}

func (u *r2Uploader) publicURL(asset largeAsset) (string, bool) {
	if u.disabled {
		return "", false
	}

	key := objectKey(asset.siteRel)
	if exists, err := u.objectExists(key); err != nil {
		u.disabled = true
		return "", false
	} else if !exists {
		if err := u.putObject(key, asset.path); err != nil {
			u.disabled = true
			return "", false
		}
	}
	return u.cfg.R2PublicBaseURL + "/" + key, true
}

func objectKey(rel string) string {
	sum := sha256.Sum256([]byte(rel))
	ext := strings.ToLower(filepath.Ext(rel))
	return hex.EncodeToString(sum[:]) + ext
}

func (u *r2Uploader) objectPath(key string) string {
	return u.cfg.R2Bucket + "/" + key
}

func (u *r2Uploader) objectExists(key string) (bool, error) {
	cmd := exec.Command(u.cfg.WranglerBin, "r2", "object", "get", u.objectPath(key), "--pipe", "--remote")
	cmd.Stdout = nil
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if isMissingObject(stderr.String()) {
			return false, nil
		}
		return false, fmt.Errorf("wrangler r2 object get: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return true, nil
}

func (u *r2Uploader) putObject(key, path string) error {
	args := []string{"r2", "object", "put", u.objectPath(key), "--file", path, "--remote", "--force"}
	if contentType := mime.TypeByExtension(filepath.Ext(path)); contentType != "" {
		args = append(args, "--content-type", contentType)
	}
	cmd := exec.Command(u.cfg.WranglerBin, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wrangler r2 object put: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

func isMissingObject(stderr string) bool {
	lower := strings.ToLower(stderr)
	return strings.Contains(lower, "not found") ||
		strings.Contains(lower, "404") ||
		strings.Contains(lower, "no such key") ||
		strings.Contains(lower, "does not exist")
}

func rewriteAssetReferences(outDir string, replacements []assetReplacement) error {
	return filepath.WalkDir(outDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !isTextAsset(path) {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		s := string(b)
		next := s
		for _, replacement := range replacements {
			url := replacement.url
			if url == placeholderName {
				url = relativeURL(filepath.Dir(path), filepath.Join(outDir, placeholderName))
			}
			next = strings.ReplaceAll(next, replacement.rel, url)
		}
		if next == s {
			return nil
		}
		return os.WriteFile(path, []byte(next), 0o644)
	})
}

func isTextAsset(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".html", ".htm", ".css", ".js", ".mjs", ".json", ".svg", ".xml", ".txt":
		return true
	default:
		return false
	}
}

func relativeURL(fromDir, target string) string {
	rel, err := filepath.Rel(fromDir, target)
	if err != nil {
		return placeholderName
	}
	return filepath.ToSlash(rel)
}

func writePlaceholder(outDir string) error {
	return os.WriteFile(filepath.Join(outDir, placeholderName), placeholderPNG, 0o644)
}

func mustDecodeBase64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}
