package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DirExists checks if a directory exists.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// EnsureDir creates a directory and all parent directories if they don't exist.
func EnsureDir(path string) error {
	if DirExists(path) {
		return nil
	}
	return os.MkdirAll(path, 0755)
}

// ReadFileString reads a file and returns its contents as a string.
func ReadFileString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return string(data), nil
}

// WriteFileString writes a string to a file.
// If overwrite is false and the file exists, it returns an error.
func WriteFileString(path, content string, overwrite bool) error {
	if !overwrite && FileExists(path) {
		return fmt.Errorf("file already exists: %s", path)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	// Ensure parent directory exists
	dir := filepath.Dir(dst)
	if err := EnsureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// ListFiles returns a list of files matching the pattern in the directory.
// Pattern supports simple glob patterns like "*.go" or "*.tmpl".
func ListFiles(dir, pattern string) ([]string, error) {
	if !DirExists(dir) {
		return nil, nil
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Match pattern
		matched, err := filepath.Match(pattern, info.Name())
		if err != nil {
			return err
		}
		if matched {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	return files, nil
}

// ListDirs returns a list of immediate subdirectories in the directory.
func ListDirs(dir string) ([]string, error) {
	if !DirExists(dir) {
		return nil, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

// RelativePath returns the path relative to the base directory.
func RelativePath(basePath, targetPath string) (string, error) {
	rel, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}
	return rel, nil
}

// AbsolutePath returns the absolute path.
func AbsolutePath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}
	return abs, nil
}

// JoinPath joins path elements.
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// CleanPath cleans a path.
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// DeleteFile removes a file.
func DeleteFile(path string) error {
	if !FileExists(path) {
		return nil
	}
	return os.Remove(path)
}

// DeleteDir removes a directory and its contents.
func DeleteDir(path string) error {
	if !DirExists(path) {
		return nil
	}
	return os.RemoveAll(path)
}

// IsEmpty checks if a directory is empty.
func IsEmpty(dir string) (bool, error) {
	if !DirExists(dir) {
		return true, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	return len(entries) == 0, nil
}

// FindProjectRoot finds the project root by looking for go.mod.
func FindProjectRoot(startPath string) (string, error) {
	absPath, err := AbsolutePath(startPath)
	if err != nil {
		return "", err
	}

	current := absPath
	for {
		goMod := filepath.Join(current, "go.mod")
		if FileExists(goMod) {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			// Reached filesystem root
			return "", fmt.Errorf("could not find go.mod in path hierarchy")
		}
		current = parent
	}
}

// GetModulePath reads the module path from go.mod.
func GetModulePath(projectRoot string) (string, error) {
	goModPath := filepath.Join(projectRoot, "go.mod")
	content, err := ReadFileString(goModPath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module "), nil
		}
	}

	return "", fmt.Errorf("module path not found in go.mod")
}

// CreateTempDir creates a temporary directory with a prefix.
func CreateTempDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}

// CreateTempFile creates a temporary file with a pattern.
func CreateTempFile(pattern string) (*os.File, error) {
	return os.CreateTemp("", pattern)
}

// AppendToFileIfNotContains appends content to a file if the file doesn't already contain
// the specified marker string. This is useful for adding instructions to existing files
// without duplicating them.
func AppendToFileIfNotContains(path, marker, content string) error {
	// Check if file exists
	if !FileExists(path) {
		return nil // File doesn't exist, nothing to append to
	}

	// Read existing content
	existing, err := ReadFileString(path)
	if err != nil {
		return err
	}

	// Check if marker already exists
	if strings.Contains(existing, marker) {
		return nil // Already contains the content
	}

	// Append new content
	newContent := existing
	if !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += "\n" + content

	return WriteFileString(path, newContent, true)
}

// CreateFileIfNotExists creates a file with the given content only if it doesn't exist.
func CreateFileIfNotExists(path, content string) error {
	if FileExists(path) {
		return nil
	}
	return WriteFileString(path, content, false)
}
