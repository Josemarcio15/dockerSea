package dockerapi

// PortMapping representa mapeamento de portas do container
type PortMapping struct {
	IP          string `json:"ip,omitempty"`
	PrivatePort uint16 `json:"privatePort"`
	PublicPort  uint16 `json:"publicPort,omitempty"`
	Type        string `json:"type"` // "tcp" | "udp"
}

// MountInfo representa volumes e binds montados no container
type MountInfo struct {
	Type        string `json:"type"`        // "volume" | "bind" | "tmpfs"
	Name        string `json:"name"`        // Nome do volume (se houver)
	Source      string `json:"source"`      // Caminho no host
	Destination string `json:"destination"` // Caminho no container
	ReadOnly    bool   `json:"readOnly"`
}

// NetworkEndpoint representa as configurações de rede do container
type NetworkEndpoint struct {
	NetworkID string `json:"networkId,omitempty"`
	IPAddress string `json:"ipAddress"`
	Gateway   string `json:"gateway"`
	MacAddress string `json:"macAddress,omitempty"`
}

// Container representa os dados consolidados e limpos de um container Docker
type Container struct {
	ID            string                     `json:"id"`
	Names         []string                   `json:"names"`
	Name          string                     `json:"name"`
	Image         string                     `json:"image"`
	ImageID       string                     `json:"imageId"`
	Command       string                     `json:"command"`
	Created       int64                      `json:"created"`
	State         string                     `json:"state"`         // "running" | "exited" | "paused" | "created"
	Status        string                     `json:"status"`        // "Up 2 hours", "Exited (0) 5 minutes ago"
	Ports         string                     `json:"ports"`         // Formato amigável: "0.0.0.0:80->80/tcp"
	PortList      []PortMapping              `json:"portList"`
	Networks      map[string]NetworkEndpoint `json:"networks"`      // Rede -> {ipAddress, gateway}
	Mounts        []MountInfo                `json:"mounts"`
	RestartPolicy string                     `json:"restartPolicy"` // "always", "unless-stopped", "on-failure", "no"
	Labels        map[string]string          `json:"labels"`
	Env           []string                   `json:"env,omitempty"`
}

// RawDockerContainer representa a estrutura padrão retornada pelo Docker Engine API (GET /containers/json)
type RawDockerContainer struct {
	ID      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	ImageID string   `json:"ImageID"`
	Command string   `json:"Command"`
	Created int64    `json:"Created"`
	State   string   `json:"State"`
	Status  string   `json:"Status"`
	Ports   []struct {
		IP          string `json:"IP"`
		PrivatePort uint16 `json:"PrivatePort"`
		PublicPort  uint16 `json:"PublicPort"`
		Type        string `json:"Type"`
	} `json:"Ports"`
	Labels  map[string]string `json:"Labels"`
	Mounts  []struct {
		Type        string `json:"Type"`
		Name        string `json:"Name"`
		Source      string `json:"Source"`
		Destination string `json:"Destination"`
		RW          bool   `json:"RW"`
	} `json:"Mounts"`
	NetworkSettings struct {
		Networks map[string]struct {
			NetworkID string `json:"NetworkID"`
			IPAddress string `json:"IPAddress"`
			Gateway   string `json:"Gateway"`
			MacAddress string `json:"MacAddress"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
	HostConfig struct {
		RestartPolicy struct {
			Name string `json:"Name"`
		} `json:"RestartPolicy"`
	} `json:"HostConfig"`
}

// RawDockerInspect representa detalhes adicionais obtidos com GET /containers/{id}/json
type RawDockerInspect struct {
	Config struct {
		Env []string `json:"Env"`
	} `json:"Config"`
	HostConfig struct {
		RestartPolicy struct {
			Name string `json:"Name"`
		} `json:"RestartPolicy"`
	} `json:"HostConfig"`
}

// DockerEvent representa um evento emitido em tempo real pelo Docker Daemon
type DockerEvent struct {
	Type   string `json:"Type"`   // "container", "image", "volume", "network"
	Action string `json:"Action"` // "start", "stop", "die", "restart", "destroy", "create", etc.
	Actor  struct {
		ID         string            `json:"ID"`
		Attributes map[string]string `json:"Attributes"`
	} `json:"Actor"`
	Time     int64 `json:"time"`
	TimeNano int64 `json:"timeNano"`
}

