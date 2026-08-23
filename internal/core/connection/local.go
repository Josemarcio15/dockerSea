package connection

import (
	"os/exec"
	"strings"
)

func execLocalCommand(cmd string) (string, error) {
	c := exec.Command("bash", "-c", cmd)
	out, err := c.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
