package volumes

// DockerVolume representa os dados normalizados de um volume Docker
type DockerVolume struct {
	Name        string            `json:"name"`
	Driver      string            `json:"driver"`
	Mountpoint  string            `json:"mountpoint"`
	CreatedAt   string            `json:"createdAt"`
	Labels      map[string]string `json:"labels"`
	Scope       string            `json:"scope"`
	InUse       bool              `json:"inUse"`
	Containers  [][]string        `json:"containers"` // [ [containerName, "ro" | "rw"] ]
	Size        string            `json:"size,omitempty"`
}

// VolumeActionResult armazena o resultado de operações em volumes
type VolumeActionResult struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Count   int      `json:"count,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// VolumeCreateRequest payload para criar volume
type VolumeCreateRequest struct {
	Name   string            `json:"name"`
	Driver string            `json:"driver,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

// RawDockerVolumeList representa o retorno de GET /volumes
type RawDockerVolumeList struct {
	Volumes  []RawDockerVolume `json:"Volumes"`
	Warnings []string          `json:"Warnings"`
}

// RawDockerVolume representa o volume retornado pela Docker API
type RawDockerVolume struct {
	Name       string            `json:"Name"`
	Driver     string            `json:"Driver"`
	Mountpoint string            `json:"Mountpoint"`
	CreatedAt  string            `json:"CreatedAt"`
	Status     map[string]string `json:"Status"`
	Labels     map[string]string `json:"Labels"`
	Scope      string            `json:"Scope"`
	UsageData  *struct {
		Size     int64 `json:"Size"`
		RefCount int64 `json:"RefCount"`
	} `json:"UsageData"`
}

// RawVolumePruneResponse representa a resposta de POST /volumes/prune
type RawVolumePruneResponse struct {
	VolumesDeleted []string `json:"VolumesDeleted"`
	SpaceReclaimed int64    `json:"SpaceReclaimed"`
}
