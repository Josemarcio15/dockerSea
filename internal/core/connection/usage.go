package connection

import (
	"fmt"
	"strconv"
	"strings"
)

// SystemUsage armazena informações de uso de memória, swap, disco e sistema
type SystemUsage struct {
	// Memória RAM (em Bytes)
	MemTotal     uint64  `json:"memTotal"`
	MemUsed      uint64  `json:"memUsed"`
	MemFree      uint64  `json:"memFree"`
	MemAvailable uint64  `json:"memAvailable"`
	MemUsagePerc float64 `json:"memUsagePerc"`

	// Memória Swap (em Bytes)
	SwapTotal     uint64  `json:"swapTotal"`
	SwapUsed      uint64  `json:"swapUsed"`
	SwapFree      uint64  `json:"swapFree"`
	SwapUsagePerc float64 `json:"swapUsagePerc"`

	// Disco Root / (em Bytes)
	DiskTotal     uint64  `json:"diskTotal"`
	DiskUsed      uint64  `json:"diskUsed"`
	DiskFree      uint64  `json:"diskFree"`
	DiskUsagePerc float64 `json:"diskUsagePerc"`

	// Informações adicionais leves
	Uptime      string `json:"uptime"`
	OSInfo      string `json:"osInfo"`
	HostName    string `json:"hostname"`
	KernelArch  string `json:"kernelArch"`
	CPUCount    int    `json:"cpuCount"`
}

const (
	// CmdGetSystemUsage executa uma leitura rápida de /proc e metadados de sistema em batch
	CmdGetSystemUsage = `free -b; echo "---DOCKSEA-SEP---"; df -B1 /; echo "---DOCKSEA-SEP---"; uptime -p 2>/dev/null || uptime; echo "---DOCKSEA-SEP---"; uname -srm; echo "---DOCKSEA-SEP---"; hostname 2>/dev/null || cat /etc/hostname 2>/dev/null || echo ""; echo "---DOCKSEA-SEP---"; nproc 2>/dev/null || grep -c ^processor /proc/cpuinfo 2>/dev/null || echo 1`
)

// FetchSystemUsage obtém uso de RAM, Swap, Disco e Uptime de forma ultra-rápida (batch single-turn)
func (c *Client) FetchSystemUsage() (*SystemUsage, error) {
	out, err := c.ExecCommand(CmdGetSystemUsage, false)
	if err != nil && strings.TrimSpace(out) == "" {
		return nil, fmt.Errorf("falha ao consultar recursos do sistema: %w", err)
	}

	return parseSystemUsage(out)
}

func parseSystemUsage(raw string) (*SystemUsage, error) {
	usage := &SystemUsage{}
	parts := strings.Split(raw, "---DOCKSEA-SEP---")
	if len(parts) < 3 {
		return usage, fmt.Errorf("formato de resposta inesperado do sistema")
	}

	// 1. Parse free -b
	freeLines := strings.Split(strings.TrimSpace(parts[0]), "\n")
	for _, line := range freeLines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if strings.HasPrefix(fields[0], "Mem:") {
			if len(fields) >= 7 {
				usage.MemTotal, _ = strconv.ParseUint(fields[1], 10, 64)
				usage.MemUsed, _ = strconv.ParseUint(fields[2], 10, 64)
				usage.MemFree, _ = strconv.ParseUint(fields[3], 10, 64)
				usage.MemAvailable, _ = strconv.ParseUint(fields[6], 10, 64)
				if usage.MemTotal > 0 {
					// Memória usada real = Total - Available (se available disponível)
					if usage.MemAvailable > 0 && usage.MemTotal >= usage.MemAvailable {
						actualUsed := usage.MemTotal - usage.MemAvailable
						usage.MemUsed = actualUsed
						usage.MemUsagePerc = (float64(actualUsed) / float64(usage.MemTotal)) * 100
					} else {
						usage.MemUsagePerc = (float64(usage.MemUsed) / float64(usage.MemTotal)) * 100
					}
				}
			}
		} else if strings.HasPrefix(fields[0], "Swap:") {
			if len(fields) >= 4 {
				usage.SwapTotal, _ = strconv.ParseUint(fields[1], 10, 64)
				usage.SwapUsed, _ = strconv.ParseUint(fields[2], 10, 64)
				usage.SwapFree, _ = strconv.ParseUint(fields[3], 10, 64)
				if usage.SwapTotal > 0 {
					usage.SwapUsagePerc = (float64(usage.SwapUsed) / float64(usage.SwapTotal)) * 100
				}
			}
		}
	}

	// 2. Parse df -B1 /
	dfLines := strings.Split(strings.TrimSpace(parts[1]), "\n")
	if len(dfLines) >= 2 {
		dfFields := strings.Fields(dfLines[len(dfLines)-1])
		if len(dfFields) >= 5 {
			usage.DiskTotal, _ = strconv.ParseUint(dfFields[1], 10, 64)
			usage.DiskUsed, _ = strconv.ParseUint(dfFields[2], 10, 64)
			usage.DiskFree, _ = strconv.ParseUint(dfFields[3], 10, 64)
			if usage.DiskTotal > 0 {
				usage.DiskUsagePerc = (float64(usage.DiskUsed) / float64(usage.DiskTotal)) * 100
			}
		}
	}

	// 3. Parse Uptime
	if len(parts) > 2 {
		usage.Uptime = strings.TrimSpace(parts[2])
	}

	// 4. Parse uname / OS
	if len(parts) > 3 {
		usage.KernelArch = strings.TrimSpace(parts[3])
		usage.OSInfo = usage.KernelArch
	}

	// 5. Parse Hostname
	if len(parts) > 4 {
		usage.HostName = strings.TrimSpace(parts[4])
	}

	// 6. Parse CPU Count
	if len(parts) > 5 {
		cpuStr := strings.TrimSpace(parts[5])
		cpuCount, _ := strconv.Atoi(cpuStr)
		if cpuCount < 1 {
			cpuCount = 1
		}
		usage.CPUCount = cpuCount
	}

	return usage, nil
}
