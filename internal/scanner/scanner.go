package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

type Scanner struct{}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Scan(root string) ([]string, error) {
	return s.ScanWithExtension(root, ".go")
}

func (s *Scanner) ScanWithExtension(root, ext string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			dir := info.Name()
			if dir == "vendor" || dir == ".git" || dir == "testdata" || dir == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		if ext == ".js" {
			if strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".test.js") && !strings.HasSuffix(path, ".spec.js") {
				files = append(files, path)
			}
		} else if ext == ".go" {
			if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				files = append(files, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
