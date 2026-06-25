package render

type Options struct {
	RepoRoot            string
	BaseRoot            string
	OutRoot             string
	Clean               bool
	Notes               bool
	R2Bucket            string
	R2PublicBaseURL     string
	WranglerBin         string
	LargeAssetThreshold int64
}

type Config struct {
	RepoRoot            string
	BaseRoot            string
	OutRoot             string
	Clean               bool
	Notes               bool
	R2Bucket            string
	R2PublicBaseURL     string
	WranglerBin         string
	LargeAssetThreshold int64
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
