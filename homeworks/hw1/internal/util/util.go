package util

import (
	"context"
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetBinPath(currentWd string) (string, error) {
	cleanPath := filepath.Clean(currentWd)

	index := strings.Index(cleanPath, "hw1")

	if index == -1 {
		return "", fmt.Errorf("root directory not found")
	}

	return filepath.Join(cleanPath[:index+3], "bin"), nil
}

func ResolveFilePath(root string, filename string) (string, error) {
	cleanedRoot := filepath.Clean(root)
	nameWithoutExt := strings.TrimRight(root, filepath.Ext(filename))

	var result string

	err := filepath.WalkDir(cleanedRoot, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		name := d.Name()

		if name == filename || name == nameWithoutExt {
			result = path
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("walk fail tree fail, error: %w", err)
	}

	if result == "" {
		return "", fmt.Errorf("file %s not found in root %s", filename, root)
	}

	return result, nil
}

func GoBuild(ctx context.Context, filepath string, outputPath string) error {
	cmd := exec.CommandContext(ctx, "go", "build", "-o", outputPath, filepath)

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
