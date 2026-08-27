package builder

import "os"

func BrowseDirectory(path string) ([]os.DirEntry, error) { return os.ReadDir(path) }
