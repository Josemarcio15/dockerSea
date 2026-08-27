package volumes

import "strings"

// NormalizeVolumeCreateRequest applies the domain defaults before persistence
// or Docker-specific code is invoked.
func NormalizeVolumeCreateRequest(req VolumeCreateRequest) VolumeCreateRequest {
	req.Name = strings.TrimSpace(req.Name)
	req.Driver = strings.TrimSpace(req.Driver)
	return req
}

func ValidateVolumeCreateRequest(req VolumeCreateRequest) string {
	if req.Name == "" {
		return "Nome do volume não informado"
	}
	return ""
}
