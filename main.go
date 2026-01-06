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

	// 2. Process Ignore List
	ignoredExts := make(map[string]bool)
	if ignoreList != "" {
		parts := strings.Split(ignoreList, ",")
		for _, p := range parts {
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

	outFile, err := os.Create(outputFilename)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	fmt.Printf("Scanning: %s\nOutput: %s\nIgnoring: %v\n", root, outputFilename, ignoredExts)

	// Call the logic function
	if err := walkAndWrite(root, outFile, outputFilename, ignoredExts); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Done!")
	}
}

// walkAndWrite is the core logic, extracted so we can test it
func walkAndWrite(root string, w io.Writer, skipFilename string, ignoredExts map[string]bool) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		if info.Name() == skipFilename {
			return nil
		}

		ext := filepath.Ext(path)
		if ignoredExts[ext] {
			return nil
		}

		isBin, err := isBinary(path)
		if err != nil {
			return nil
		}
		if isBin {
			return nil
		}

		header := fmt.Sprintf("\n========================================\nPATH: %s\n---\n", path)
		if _, err := w.Write([]byte(header)); err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer srcFile.Close()

		if _, err := io.Copy(w, srcFile); err != nil {
			return err
		}

		w.Write([]byte("\n"))
		return nil
	})
}

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