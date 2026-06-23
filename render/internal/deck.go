package render

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func discoverDecks(repoRoot string, args []string) ([]deck, error) {
	if len(args) > 0 {
		decks := make([]deck, 0, len(args))
		for _, arg := range args {
			d, err := deckFromArg(repoRoot, arg)
			if err != nil {
				return nil, err
			}
			decks = append(decks, d)
		}
		sortDecks(decks)
		return decks, nil
	}

	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return nil, err
	}

	var decks []deck
	for _, entry := range entries {
		if !entry.IsDir() || shouldSkipTopLevel(entry.Name()) {
			continue
		}
		dir := filepath.Join(repoRoot, entry.Name())
		slides, err := slideFiles(dir)
		if err != nil {
			return nil, err
		}
		if len(slides) == 0 {
			continue
		}
		decks = append(decks, deck{name: filepath.ToSlash(entry.Name()), dir: dir, slides: slides})
	}
	sortDecks(decks)
	return decks, nil
}

func deckFromArg(repoRoot, arg string) (deck, error) {
	path := resolvePath(repoRoot, arg)
	info, err := os.Stat(path)
	if err != nil {
		return deck{}, err
	}
	if !info.IsDir() {
		if filepath.Ext(path) != ".slide" {
			return deck{}, fmt.Errorf("%s is not a .slide file or directory", arg)
		}
		dir := filepath.Dir(path)
		name, err := filepath.Rel(repoRoot, dir)
		if err != nil {
			return deck{}, err
		}
		return deck{name: filepath.ToSlash(name), dir: dir, slides: []string{path}}, nil
	}

	slides, err := slideFiles(path)
	if err != nil {
		return deck{}, err
	}
	if len(slides) == 0 {
		return deck{}, fmt.Errorf("%s contains no .slide files", arg)
	}
	name, err := filepath.Rel(repoRoot, path)
	if err != nil {
		return deck{}, err
	}
	return deck{name: filepath.ToSlash(name), dir: path, slides: slides}, nil
}

func slideFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var slides []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".slide" {
			continue
		}
		slides = append(slides, filepath.Join(dir, entry.Name()))
	}
	sort.Strings(slides)
	return slides, nil
}

func sortDecks(decks []deck) {
	sort.Slice(decks, func(i, j int) bool {
		return decks[i].name < decks[j].name
	})
}

func shouldSkipTopLevel(name string) bool {
	return strings.HasPrefix(name, ".") ||
		name == "dist" ||
		name == "present-assets" ||
		name == "render"
}
