package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func EnsureDir(dir string, root bool) (string, error) {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		absPath = dir
	}

	info, err := os.Stat(absPath)
	if err == nil {
		if info.IsDir() {
			return absPath, nil
		}
		return "", fmt.Errorf("path exists but is not a directory: %s", absPath)
	}

	if !os.IsNotExist(err) {
		return "", err
	}

	if root {
		err = os.MkdirAll(absPath, 0755)
	} else {
		err = os.Mkdir(absPath, 0755)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create directory %q: %w", absPath, err)
	}

	return absPath, nil
}
