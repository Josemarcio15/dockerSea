package extras

// NginxSitesResult contém os arquivos em sites-available e sites-enabled
type NginxSitesResult struct {
	Available []string `json:"available"`
	Enabled   []string `json:"enabled"`
}

// ExtraActionResult contém o retorno de comandos e operações de extras
type ExtraActionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Output  string `json:"output"`
}

// PortEntry representa uma porta em escuta no servidor
type PortEntry struct {
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	ProcessName string `json:"processName"`
	PID         string `json:"pid"`
	Address     string `json:"address"`
}
