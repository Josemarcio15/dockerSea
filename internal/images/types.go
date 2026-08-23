package images

// DockerImage representa uma imagem Docker com metadados formatados
type DockerImage struct {
	ID              string   `json:"id"`
	Repo            string   `json:"repo"`
	Tag             string   `json:"tag"`
	Size            string   `json:"size"`
	RawSizeBytes    int64    `json:"rawSizeBytes"`
	Created         int64    `json:"created"`
	ContainersUsing []string `json:"containersUsing"`
	VirtualSize     int64    `json:"virtualSize,omitempty"`
	SharedSize      int64    `json:"sharedSize,omitempty"`
}

// ImageActionResult armazena o resultado de operações em imagens
type ImageActionResult struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Count   int      `json:"count,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// RawDockerImage representa a estrutura retornada por GET /images/json
type RawDockerImage struct {
	ID          string            `json:"Id"`
	ParentID    string            `json:"ParentId"`
	RepoTags    []string          `json:"RepoTags"`
	RepoDigests []string          `json:"RepoDigests"`
	Created     int64             `json:"Created"`
	Size        int64             `json:"Size"`
	VirtualSize int64             `json:"VirtualSize"`
	SharedSize  int64             `json:"SharedSize"`
	Labels      map[string]string `json:"Labels"`
	Containers  int               `json:"Containers"`
}

// DockerPullProgress representa o payload enviado linha a linha pelo Docker Engine durante POST /images/create
type DockerPullProgress struct {
	ID             string `json:"id,omitempty"`
	Status         string `json:"status,omitempty"`
	Progress       string `json:"progress,omitempty"`
	ProgressDetail struct {
		Current int64 `json:"current,omitempty"`
		Total   int64 `json:"total,omitempty"`
	} `json:"progressDetail,omitempty"`
	Error       string `json:"error,omitempty"`
	Line        string `json:"line,omitempty"`
	PullPercent int    `json:"pullPercent,omitempty"`
}

// RawImageDeleteResponse representa a resposta de DELETE /images/{id}
type RawImageDeleteResponse struct {
	Untagged string `json:"Untagged,omitempty"`
	Deleted  string `json:"Deleted,omitempty"`
}

// DockerTransferProgress representa o progresso em tempo real da transferência de imagens
type DockerTransferProgress struct {
	Stage               string `json:"stage"` // "preparing" | "transferring" | "loading" | "complete"
	BytesSent           int64  `json:"bytesSent"`
	FormattedBytes      string `json:"formattedBytes"`
	TotalBytes          int64  `json:"totalBytes"`
	FormattedTotalBytes string `json:"formattedTotalBytes"`
	Speed               string `json:"speed"`
	Percent             int    `json:"percent"`
	Message             string `json:"message"`
	Error               string `json:"error,omitempty"`
}
