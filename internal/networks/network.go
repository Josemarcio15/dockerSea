package networks

import "strings"

func NormalizeNetworkCreateRequest(req NetworkCreateRequest) NetworkCreateRequest {
	req.Name = strings.TrimSpace(req.Name)
	if req.Driver == "" {
		req.Driver = "bridge"
	}
	return req
}

func ValidateNetworkCreateRequest(req NetworkCreateRequest) string {
	if strings.TrimSpace(req.Name) == "" {
		return "Nome da rede não informado"
	}
	return ""
}