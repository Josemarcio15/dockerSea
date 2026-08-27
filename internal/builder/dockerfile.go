package builder

import "os"

func ReadDockerfile(path string) ([]byte, error) { return os.ReadFile(path) }
