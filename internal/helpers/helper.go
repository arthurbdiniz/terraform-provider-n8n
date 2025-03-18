// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package helpers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// findRepoRoot traverses upwards from the current working directory to find the root of the repository.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	marker := "README.md"

	// Traverse upwards through the directories until we find the marker file or reach the root
	for {
		markerPath := filepath.Join(dir, marker)
		if _, err := os.Stat(markerPath); !os.IsNotExist(err) {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			break
		}
		dir = parentDir
	}

	return "", fmt.Errorf("repository root not found")
}

// getTestDataFilePath returns the absolute path to files under testdata folder, relative to the root of the repo.
func getTestDataFilePath(filename string) (string, error) {
	root, err := findRepoRoot()
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(root, "testdata", filename)

	return filepath.Abs(configPath)
}

// CreateTestContainer creates and starts a container with the given configuration paths.
func CreateTestContainer() (testcontainers.Container, string, error) {
	configPath, err := getTestDataFilePath("config")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get config path: %w", err)
	}

	databasePath, err := getTestDataFilePath("database.sqlite")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get database path: %w", err)
	}

	req := testcontainers.ContainerRequest{
		Image:        "docker.n8n.io/n8nio/n8n:1.84.1",
		ExposedPorts: []string{"5678/tcp"},
		Files: []testcontainers.ContainerFile{
			{
				HostFilePath:      configPath,
				ContainerFilePath: "/home/node/.n8n/config",
				FileMode:          0o777,
			},
			{
				HostFilePath:      databasePath,
				ContainerFilePath: "/home/node/.n8n/database.sqlite",
				FileMode:          0o777,
			},
		},
		Env: map[string]string{
			"N8N_ENFORCE_SETTINGS_FILE_PERMISSIONS": "false",
			"N8N_RUNNERS_ENABLED":                   "true",
		},
		WaitingFor: wait.ForLog("http://localhost:5678").WithStartupTimeout(30 * time.Second),
	}

	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("could not start container: %w", err)
	}

	// Get the mapped host
	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("could not get container host: %s", err)
	}

	// Get the mapped port
	port, err := container.MappedPort(ctx, "5678")
	if err != nil {
		return nil, "", fmt.Errorf("could not get mapped port: %s", err)
	}

	url := fmt.Sprintf("http://%s:%s", host, port.Port())

	return container, url, nil
}

// TerminateContainer terminates the container and handles errors.
func TerminateContainer(ctx context.Context, container testcontainers.Container) error {
	if err := container.Terminate(ctx); err != nil {
		return fmt.Errorf("error when terminating container: %w", err)
	}
	return nil
}
