package networks

// DockerNetworkContainer representa um container conectado a uma rede
type DockerNetworkContainer struct {
	Name        string `json:"name"`
	EndpointID  string `json:"endpointId"`
	MacAddress  string `json:"macAddress"`
	IPv4Address string `json:"ip"`
	IPv6Address string `json:"ipv6Address,omitempty"`
}

// DockerNetwork representa os dados normalizados de uma rede Docker
type DockerNetwork struct {
	ID         string                   `json:"id"`
	Name       string                   `json:"name"`
	Driver     string                   `json:"driver"`
	Scope      string                   `json:"scope"`
	Subnet     string                   `json:"subnet"`
	Gateway    string                   `json:"gateway"`
	Internal   bool                     `json:"internal"`
	Attachable bool                     `json:"attachable"`
	Containers []DockerNetworkContainer `json:"containers"`
	Labels     map[string]string        `json:"labels"`
}

// NetworkActionResult armazena o resultado de operações em redes
type NetworkActionResult struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Count   int      `json:"count,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// NetworkCreateRequest representa a requisição para criação de uma rede
type NetworkCreateRequest struct {
	Name    string            `json:"name"`
	Driver  string            `json:"driver,omitempty"`
	Subnet  string            `json:"subnet,omitempty"`
	Gateway string            `json:"gateway,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

// RawDockerNetwork representa a estrutura retornada pelo Docker Engine API (GET /networks)
type RawDockerNetwork struct {
	Name   string `json:"Name"`
	ID     string `json:"Id"`
	Scope  string `json:"Scope"`
	Driver string `json:"Driver"`
	IPAM   struct {
		Driver string `json:"Driver"`
		Config []struct {
			Subnet  string `json:"Subnet"`
			Gateway string `json:"Gateway"`
		} `json:"Config"`
	} `json:"IPAM"`
	Internal   bool `json:"Internal"`
	Attachable bool `json:"Attachable"`
	Containers map[string]struct {
		Name        string `json:"Name"`
		EndpointID  string `json:"EndpointID"`
		MacAddress  string `json:"MacAddress"`
		IPv4Address string `json:"IPv4Address"`
		IPv6Address string `json:"IPv6Address"`
	} `json:"Containers"`
	Labels map[string]string `json:"Labels"`
}

// RawNetworkPruneResponse representa a resposta de POST /networks/prune
type RawNetworkPruneResponse struct {
	NetworksDeleted []string `json:"NetworksDeleted"`
}
