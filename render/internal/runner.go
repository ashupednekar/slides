package render

import (
	"fmt"

	"golang.org/x/tools/present"
)

func Run(cfg Config, args []string) ([]RenderedDeck, error) {
	present.PlayEnabled = false
	present.NotesEnabled = cfg.Notes

	tmpl, err := slideTemplate(cfg.BaseRoot)
	if err != nil {
		return nil, err
	}

	decks, err := discoverDecks(cfg.RepoRoot, args)
	if err != nil {
		return nil, err
	}
	if len(decks) == 0 {
		return nil, fmt.Errorf("no slide directories found")
	}

	if err := prepareOutput(cfg); err != nil {
		return nil, err
	}

	rendered := make([]RenderedDeck, 0, len(decks))
	for _, d := range decks {
		rd, err := renderDeck(cfg, tmpl, d)
		if err != nil {
			return nil, err
		}
		rendered = append(rendered, rd)
	}

	if err := writeIndex(cfg.OutRoot, rendered); err != nil {
		return nil, err
	}
	return rendered, nil
}
