package scanner

// FileInfo represents metadata for a scanned file.
type FileInfo struct {
	Path  string `json:"path"`
	Name  string `json:"name"`
	Lines int    `json:"lines"`
	Size  int64  `json:"size"`
}
