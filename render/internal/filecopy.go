package render

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == src {
			return nil
		}
		if skipCopiedPath(entry) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		info, err := entry.Info()
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return nil
		}
		return copyFile(path, target, info.Mode())
	})
}

func skipCopiedPath(entry fs.DirEntry) bool {
	name := entry.Name()
	if name == ".DS_Store" || name == "__pycache__" {
		return true
	}
	if strings.HasSuffix(name, ".pyc") || strings.HasSuffix(name, ".slide") {
		return true
	}
	return false
}

func copyFile(src, dst string, mode fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode.Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
