// Copyright (c) Arthur Diniz <arthurbdiniz@gmail.com>
// SPDX-License-Identifier: Apache-2.0

package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindRepoRoot_Success(t *testing.T) {
	// Create a temp dir with README.md
	tmpDir := t.TempDir()
	readmePath := filepath.Join(tmpDir, "README.md")
	err := os.WriteFile(readmePath, []byte("Test Readme"), 0644)
	assert.NoError(t, err)

	// Change to that dir
	originalWD, _ := os.Getwd()
	defer func() {
		err := os.Chdir(originalWD)
		assert.NoError(t, err)
	}()

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	root, err := findRepoRoot()
	assert.NoError(t, err)
	assert.Equal(t, tmpDir, root)
}

func TestFindRepoRoot_NotFound(t *testing.T) {
	// Create a temp dir with no README.md
	tmpDir := t.TempDir()

	originalWD, _ := os.Getwd()
	defer func() {
		err := os.Chdir(originalWD)
		assert.NoError(t, err)
	}()

	err := os.Chdir(tmpDir)
	assert.NoError(t, err)

	root, err := findRepoRoot()
	assert.Error(t, err)
	assert.Equal(t, "", root)
}

func TestGetTestDataFilePath_Success(t *testing.T) {
	// Create a temp dir with README.md and testdata/file.txt
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(""), 0644)
	assert.NoError(t, err)

	testdataDir := filepath.Join(tmpDir, "testdata")
	err = os.Mkdir(testdataDir, 0755)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(testdataDir, "file.txt"), []byte("ok"), 0644)
	assert.NoError(t, err)

	originalWD, _ := os.Getwd()
	defer func() {
		err := os.Chdir(originalWD)
		assert.NoError(t, err)
	}()
	err = os.Chdir(testdataDir) // simulate nested inside repo
	assert.NoError(t, err)

	path, err := getTestDataFilePath("file.txt")
	assert.NoError(t, err)
	assert.FileExists(t, path)
}
