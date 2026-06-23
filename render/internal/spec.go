package render

type Options struct {
	RepoRoot string
	BaseRoot string
	OutRoot  string
	Clean    bool
	Notes    bool
}

type Config struct {
	RepoRoot string
	BaseRoot string
	OutRoot  string
	Clean    bool
	Notes    bool
}

type RenderedDeck struct {
	Name   string
	Title  string
	OutDir string
}

type deck struct {
	name   string
	dir    string
	slides []string
}
