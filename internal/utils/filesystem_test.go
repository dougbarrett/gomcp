package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirExists(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-dir-exists")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"existing directory", tmpDir, true},
		{"non-existing directory", filepath.Join(tmpDir, "nonexistent"), false},
		{"file not directory", "", false}, // will be set below
	}

	// Create a file for testing
	testFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	tests[2].path = testFile

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DirExists(tt.path)
			if got != tt.want {
				t.Errorf("DirExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-file-exists")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{"existing file", testFile, true},
		{"non-existing file", filepath.Join(tmpDir, "nonexistent.txt"), false},
		{"directory not file", tmpDir, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FileExists(tt.path)
			if got != tt.want {
				t.Errorf("FileExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestEnsureDir(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-ensure-dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"new directory", filepath.Join(tmpDir, "newdir"), false},
		{"nested directory", filepath.Join(tmpDir, "a", "b", "c"), false},
		{"existing directory", tmpDir, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureDir(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("EnsureDir(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !DirExists(tt.path) {
				t.Errorf("EnsureDir(%q) did not create directory", tt.path)
			}
		})
	}
}

func TestReadWriteFileString(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-read-write")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!\nLine 2"

	// Test write
	err = WriteFileString(testFile, content, false)
	if err != nil {
		t.Fatalf("WriteFileString failed: %v", err)
	}

	// Test read
	got, err := ReadFileString(testFile)
	if err != nil {
		t.Fatalf("ReadFileString failed: %v", err)
	}
	if got != content {
		t.Errorf("ReadFileString() = %q, want %q", got, content)
	}

	// Test overwrite=false with existing file
	err = WriteFileString(testFile, "new content", false)
	if err == nil {
		t.Error("WriteFileString with overwrite=false should fail for existing file")
	}

	// Test overwrite=true
	newContent := "Updated content"
	err = WriteFileString(testFile, newContent, true)
	if err != nil {
		t.Fatalf("WriteFileString with overwrite=true failed: %v", err)
	}

	got, err = ReadFileString(testFile)
	if err != nil {
		t.Fatalf("ReadFileString failed: %v", err)
	}
	if got != newContent {
		t.Errorf("ReadFileString() = %q, want %q", got, newContent)
	}
}

func TestWriteFileString_CreatesParentDir(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-write-parent")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write to nested path that doesn't exist
	testFile := filepath.Join(tmpDir, "a", "b", "c", "test.txt")
	content := "nested content"

	err = WriteFileString(testFile, content, false)
	if err != nil {
		t.Fatalf("WriteFileString failed to create parent dirs: %v", err)
	}

	got, err := ReadFileString(testFile)
	if err != nil {
		t.Fatalf("ReadFileString failed: %v", err)
	}
	if got != content {
		t.Errorf("ReadFileString() = %q, want %q", got, content)
	}
}

func TestCopyFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-copy")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	content := "Copy me!"
	if err := os.WriteFile(srcFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy to destination
	dstFile := filepath.Join(tmpDir, "dest.txt")
	err = CopyFile(srcFile, dstFile)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	// Verify content
	got, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(got) != content {
		t.Errorf("Copied content = %q, want %q", string(got), content)
	}

	// Test copy to nested path
	dstNested := filepath.Join(tmpDir, "nested", "dir", "dest.txt")
	err = CopyFile(srcFile, dstNested)
	if err != nil {
		t.Fatalf("CopyFile to nested path failed: %v", err)
	}
}

func TestListFiles(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-list-files")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := []string{
		"file1.go",
		"file2.go",
		"file3.txt",
		"subdir/file4.go",
	}

	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// List .go files
	goFiles, err := ListFiles(tmpDir, "*.go")
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	// Should find 3 .go files (including in subdirectory)
	if len(goFiles) != 3 {
		t.Errorf("ListFiles(*.go) found %d files, want 3", len(goFiles))
	}

	// List non-existent directory
	files2, err := ListFiles(filepath.Join(tmpDir, "nonexistent"), "*.go")
	if err != nil {
		t.Fatalf("ListFiles on non-existent dir should not error: %v", err)
	}
	if files2 != nil && len(files2) != 0 {
		t.Errorf("ListFiles on non-existent dir should return nil or empty slice")
	}
}

func TestListDirs(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-list-dirs")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectories
	dirs := []string{"dir1", "dir2", "dir3"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, d), 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
	}

	// Create a file (should not be included)
	if err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// List directories
	got, err := ListDirs(tmpDir)
	if err != nil {
		t.Fatalf("ListDirs failed: %v", err)
	}

	if len(got) != 3 {
		t.Errorf("ListDirs found %d dirs, want 3", len(got))
	}
}

func TestRelativePath(t *testing.T) {
	tests := []struct {
		base   string
		target string
		want   string
	}{
		{"/a/b", "/a/b/c/d", "c/d"},
		{"/a/b/c", "/a/b/c", "."},
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			got, err := RelativePath(tt.base, tt.target)
			if err != nil {
				t.Fatalf("RelativePath failed: %v", err)
			}
			if got != tt.want {
				t.Errorf("RelativePath(%q, %q) = %q, want %q", tt.base, tt.target, got, tt.want)
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-delete-file")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Delete existing file
	err = DeleteFile(testFile)
	if err != nil {
		t.Fatalf("DeleteFile failed: %v", err)
	}

	if FileExists(testFile) {
		t.Error("File still exists after DeleteFile")
	}

	// Delete non-existing file (should not error)
	err = DeleteFile(testFile)
	if err != nil {
		t.Errorf("DeleteFile on non-existing file should not error: %v", err)
	}
}

func TestDeleteDir(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-delete-dir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create nested structure
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Delete directory
	err = DeleteDir(tmpDir)
	if err != nil {
		t.Fatalf("DeleteDir failed: %v", err)
	}

	if DirExists(tmpDir) {
		t.Error("Directory still exists after DeleteDir")
	}

	// Delete non-existing directory (should not error)
	err = DeleteDir(tmpDir)
	if err != nil {
		t.Errorf("DeleteDir on non-existing dir should not error: %v", err)
	}
}

func TestIsEmpty(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "test-is-empty")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Empty directory
	emptyDir := filepath.Join(tmpDir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatalf("Failed to create empty dir: %v", err)
	}

	empty, err := IsEmpty(emptyDir)
	if err != nil {
		t.Fatalf("IsEmpty failed: %v", err)
	}
	if !empty {
		t.Error("IsEmpty should return true for empty directory")
	}

	// Non-empty directory
	if err := os.WriteFile(filepath.Join(emptyDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	empty, err = IsEmpty(emptyDir)
	if err != nil {
		t.Fatalf("IsEmpty failed: %v", err)
	}
	if empty {
		t.Error("IsEmpty should return false for non-empty directory")
	}

	// Non-existent directory
	empty, err = IsEmpty(filepath.Join(tmpDir, "nonexistent"))
	if err != nil {
		t.Fatalf("IsEmpty on non-existent should not error: %v", err)
	}
	if !empty {
		t.Error("IsEmpty should return true for non-existent directory")
	}
}

func TestFindProjectRoot(t *testing.T) {
	// Create temp directory with go.mod
	tmpDir, err := os.MkdirTemp("", "test-find-root")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create nested directory
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested dirs: %v", err)
	}

	// Find project root from nested dir
	root, err := FindProjectRoot(nestedDir)
	if err != nil {
		t.Fatalf("FindProjectRoot failed: %v", err)
	}

	absRoot, _ := filepath.Abs(tmpDir)
	if root != absRoot {
		t.Errorf("FindProjectRoot() = %q, want %q", root, absRoot)
	}

	// Find from root itself
	root, err = FindProjectRoot(tmpDir)
	if err != nil {
		t.Fatalf("FindProjectRoot failed: %v", err)
	}
	if root != absRoot {
		t.Errorf("FindProjectRoot() = %q, want %q", root, absRoot)
	}
}

func TestGetModulePath(t *testing.T) {
	// Create temp directory with go.mod
	tmpDir, err := os.MkdirTemp("", "test-get-module")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create go.mod
	goModContent := "module github.com/test/project\n\ngo 1.21\n"
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Get module path
	got, err := GetModulePath(tmpDir)
	if err != nil {
		t.Fatalf("GetModulePath failed: %v", err)
	}

	want := "github.com/test/project"
	if got != want {
		t.Errorf("GetModulePath() = %q, want %q", got, want)
	}
}

func TestCreateTempDir(t *testing.T) {
	dir, err := CreateTempDir("test-prefix")
	if err != nil {
		t.Fatalf("CreateTempDir failed: %v", err)
	}
	defer os.RemoveAll(dir)

	if !DirExists(dir) {
		t.Error("CreateTempDir did not create directory")
	}
}

func TestCreateTempFile(t *testing.T) {
	file, err := CreateTempFile("test-*.txt")
	if err != nil {
		t.Fatalf("CreateTempFile failed: %v", err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if !FileExists(file.Name()) {
		t.Error("CreateTempFile did not create file")
	}
}
