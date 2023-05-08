package test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/danielchalef/zep/config"
	"os"
	"path/filepath"
	"runtime"
)

const TestDsn string = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func NewTestConfig() (*config.Config, error) {
	projectRoot, err := FindProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("Failed to find project root: %v", err)
	}
	configPath := filepath.Join(projectRoot, "config.yaml")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	return cfg, nil
}

func GenerateRandomSessionID(length int) (string, error) {
	bytes := make([]byte, (length+1)/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random session ID: %w", err)
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// FindProjectRoot returns the absolute path to the project root directory.
func FindProjectRoot() (string, error) {
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not get current file path")
	}

	dir := filepath.Dir(currentFilePath)

	for {
		// Check if the current directory contains a marker file or directory that indicates the project root.
		// In this case, we use "go.mod" as an example, but you can use any other marker.
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		// If we've reached the top-level directory, the project root is not found.
		if dir == filepath.Dir(dir) {
			return "", fmt.Errorf("project root not found")
		}

		// Move up one directory level.
		dir = filepath.Dir(dir)
	}
}