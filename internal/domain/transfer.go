package domain

type FileEntry struct {
	RelPath string `json:"path"`
	Size    int64  `json:"size"`
	Mode    uint32 `json:"mode"`
	IsDir   bool   `json:"is_dir,omitempty"`

	AbsPath string `json:"-"`
}

type Manifest struct {
	Sender     string      `json:"sender"`
	TotalBytes int64       `json:"total_bytes"`
	Files      []FileEntry `json:"files"`
}

type AcceptResponse struct {
	Accepted bool   `json:"accepted"`
	Reason   string `json:"reason,omitempty"`
}
