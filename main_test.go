package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIsBinary verifies that our binary detection logic works
func TestIsBinary(t *testing.T) {
	// 1. Create a dummy text file
	tmpText, _ := os.CreateTemp("", "test_*.txt")
	defer os.Remove(tmpText.Name())
	tmpText.WriteString("Hello world")
	tmpText.Close()

	// 2. Create a dummy binary file (contains null byte)
	tmpBin, _ := os.CreateTemp("", "test_*.bin")
	defer os.Remove(tmpBin.Name())
	tmpBin.Write([]byte{0x00, 0x01, 0x02}) // Write raw binary
	tmpBin.Close()

	// Assertions
	isBin, _ := isBinary(tmpText.Name())
	if isBin {
		t.Errorf("Expected text file to be detected as text, got binary")
	}

	isBin, _ = isBinary(tmpBin.Name())
	if !isBin {
		t.Errorf("Expected binary file with null bytes to be detected as binary, got text")
	}
}

// TestWalkAndWrite verifies the full integration
func TestWalkAndWrite(t *testing.T) {
	// 1. Setup a Temporary Directory for testing
	tempDir, err := os.MkdirTemp("", "filepecker_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	// 2. Create sample files inside tempDir
	
	// A valid Go file (Should be included)
	os.WriteFile(filepath.Join(tempDir, "main.go"), []byte("package main"), 0644)
	
	// A valid Text file (Should be included)
	os.WriteFile(filepath.Join(tempDir, "read.txt"), []byte("some text"), 0644)

	// An Ignored JSON file (Should be ignored by flag)
	os.WriteFile(filepath.Join(tempDir, "data.json"), []byte("{}"), 0644)

	// A Binary file (Should be ignored by detection)
	os.WriteFile(filepath.Join(tempDir, "image.png"), []byte{0x00, 0x00, 0x00}, 0644)

	// A .git directory (Should be ignored by logic)
	os.Mkdir(filepath.Join(tempDir, ".git"), 0755)
	os.WriteFile(filepath.Join(tempDir, ".git", "config"), []byte("secret"), 0644)

	// 3. Run the Logic
	var outputBuffer bytes.Buffer // We write to memory instead of a real file
	ignored := map[string]bool{".json": true}
	
	err = walkAndWrite(tempDir, &outputBuffer, "output.txt", ignored)
	if err != nil {
		t.Fatalf("walkAndWrite returned error: %v", err)
	}

	// 4. Verify Output
	result := outputBuffer.String()

	// Check positives
	if !strings.Contains(result, "package main") {
		t.Error("Output should contain content of main.go")
	}
	if !strings.Contains(result, "some text") {
		t.Error("Output should contain content of read.txt")
	}

	// Check negatives
	if strings.Contains(result, "{}") {
		t.Error("Output should NOT contain content of ignored .json file")
	}
	if strings.Contains(result, "secret") {
		t.Error("Output should NOT contain content of .git folder")
	}
}