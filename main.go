package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 1. Define Command Line Flags
	outputPtr := flag.String("o", "file.txt", "Output filename")
	ignorePtr := flag.String("ignore", "", "Comma-separated extensions to ignore (e.g. .json,.md)")
	flag.Parse()

	outputFilename := *outputPtr
	ignoreList := *ignorePtr

	// 2. Process Ignore List into a Map for fast lookup
	ignoredExts := make(map[string]bool)
	if ignoreList != "" {
		parts := strings.Split(ignoreList, ",")
		for _, p := range parts {
			// Ensure extension starts with a dot and trim whitespace
			cleanExt := strings.TrimSpace(p)
			if !strings.HasPrefix(cleanExt, ".") {
				cleanExt = "." + cleanExt
			}
			ignoredExts[cleanExt] = true
		}
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// 3. Create the output file
	outFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	fmt.Printf("Scanning: %s\nOutput: %s\nIgnoring: %v\n", root, outputFilename, ignoredExts)

	// 4. Start Walking
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip hidden directories like .git, .idea, .vscode
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// Check 1: Skip the output file itself
		if info.Name() == outputFilename {
			return nil
		}

		// Check 2: Skip ignored extensions
		ext := filepath.Ext(path)
		if ignoredExts[ext] {
			return nil
		}

		// Check 3: Check if file is binary
		isBin, err := isBinary(path)
		if err != nil {
			// If we can't read it, skip it
			return nil
		}
		if isBin {
			// Silent skip for binaries to keep console clean
			return nil
		}

		// 5. Write Header
		header := fmt.Sprintf("\n========================================\nPATH: %s\n---\n", path)
		if _, err := outFile.WriteString(header); err != nil {
			return err
		}

		// 6. Write Content
		srcFile, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer srcFile.Close()

		if _, err := io.Copy(outFile, srcFile); err != nil {
			return err
		}

		outFile.WriteString("\n")
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
	} else {
		fmt.Println("Done!")
	}
}

// isBinary checks the first 512 bytes to guess if a file is binary or text.
func isBinary(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buffer := make([]byte, 512)
	n, err := f.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}
	if n == 0 {
		return false, nil
	}

	if bytes.Contains(buffer[:n], []byte{0}) {
		return true, nil
	}

	contentType := http.DetectContentType(buffer[:n])
	if strings.Contains(contentType, "octet-stream") {
		return true, nil
	}

	return false, nil
}