package dao

type Metadata struct {
	Name    string `json:"name"`
	Version int    `json:"version"`
	Size    int64  `json:"size"`
	Hash    string `json:"hash"`
}
