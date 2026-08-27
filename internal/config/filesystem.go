package config

import "os"

func ReadConfigFile(path string) ([]byte, error) { return os.ReadFile(path) }
