package builder

// ValidateBuildInput keeps build-input validation in a dedicated boundary.
func ValidateBuildInput(folderPath, projectName string) string {
	if message := ValidateProjectName(projectName); message != "" {
		return message
	}
	if folderPath == "" {
		return "pasta do projeto é obrigatória"
	}
	return ""
}
