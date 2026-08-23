package extras

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"go-walis/internal/core/connection"
)

// listListeningPorts consulta portas TCP e UDP em escuta no servidor usando ss / netstat / lsof
func listListeningPorts(client *connection.Client) ([]PortEntry, error) {
	cmd := `
		if command -v ss >/dev/null 2>&1; then
			ss -tulpn -H 2>/dev/null || ss -tuln -H 2>/dev/null
		elif command -v netstat >/dev/null 2>&1; then
			netstat -tulpn 2>/dev/null || netstat -tuln 2>/dev/null
		elif command -v lsof >/dev/null 2>&1; then
			lsof -iTCP -iUDP -sTCP:LISTEN -P -n 2>/dev/null
		else
			echo ""
		fi
	`
	out, err := client.ExecCommand(cmd, true)
	if err != nil && strings.TrimSpace(out) == "" {
		return nil, fmt.Errorf("falha ao obter portas em escuta: %w", err)
	}

	var ports []PortEntry
	lines := strings.Split(out, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "Proto") || strings.HasPrefix(trimmed, "Active") || strings.HasPrefix(trimmed, "COMMAND") {
			continue
		}

		fields := strings.Fields(trimmed)
		if len(fields) < 4 {
			continue
		}

		proto := strings.ToLower(fields[0])
		if !strings.HasPrefix(proto, "tcp") && !strings.HasPrefix(proto, "udp") {
			continue
		}

		var localAddr string
		var processInfo string

		if strings.Contains(line, "LISTEN") || strings.Contains(line, "UNCONN") {
			if len(fields) >= 5 {
				localAddr = fields[4]
			}
			if len(fields) >= 7 {
				processInfo = strings.Join(fields[6:], " ")
			}
		} else if strings.HasPrefix(proto, "tcp") || strings.HasPrefix(proto, "udp") {
			if len(fields) >= 4 {
				localAddr = fields[3]
			}
			if len(fields) >= 7 {
				processInfo = fields[6]
			}
		}

		if localAddr == "" {
			continue
		}

		lastColon := strings.LastIndex(localAddr, ":")
		if lastColon == -1 {
			continue
		}

		portStr := localAddr[lastColon+1:]
		addrStr := localAddr[:lastColon]
		if addrStr == "" {
			addrStr = "*"
		}

		portNum, err := strconv.Atoi(portStr)
		if err != nil || portNum <= 0 {
			continue
		}

		procName := "Desconhecido"
		pidStr := ""

		if processInfo != "" {
			if strings.Contains(processInfo, "users:((") {
				start := strings.Index(processInfo, `("`)
				if start != -1 {
					rest := processInfo[start+2:]
					endQuote := strings.Index(rest, `"`)
					if endQuote != -1 {
						procName = rest[:endQuote]
					}
				}
				pidIdx := strings.Index(processInfo, "pid=")
				if pidIdx != -1 {
					pidRest := processInfo[pidIdx+4:]
					pidEnd := strings.IndexAny(pidRest, ",)")
					if pidEnd != -1 {
						pidStr = pidRest[:pidEnd]
					}
				}
			} else if strings.Contains(processInfo, "/") {
				parts := strings.SplitN(processInfo, "/", 2)
				pidStr = parts[0]
				procName = parts[1]
			} else {
				procName = processInfo
			}
		}

		ports = append(ports, PortEntry{
			Port:        portNum,
			Protocol:    strings.ToUpper(proto),
			ProcessName: procName,
			PID:         pidStr,
			Address:     addrStr,
		})
	}

	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Port < ports[j].Port
	})

	return ports, nil
}
