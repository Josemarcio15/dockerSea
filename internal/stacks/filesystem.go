package stacks

import "os"

func ReadStackFile(path string) ([]byte, error) { return os.ReadFile(path) }
