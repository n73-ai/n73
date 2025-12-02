package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Unzip extracts src zip file into dest directory.
// It prevents ZipSlip by ensuring each target path is within dest.
func Unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	// Ensure destination dir exists
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("mkdir dest: %w", err)
	}

	destClean := filepath.Clean(dest) + string(os.PathSeparator)

	for _, f := range r.File {
		// Normalize file path
		fname := filepath.Clean(f.Name)

		// skip empty names
		if fname == "." || fname == string(os.PathSeparator) || fname == "" {
			continue
		}

		targetPath := filepath.Join(dest, fname)

		// Prevent ZipSlip: targetPath must start with destClean
		if !strings.HasPrefix(filepath.Clean(targetPath)+string(os.PathSeparator), destClean) &&
			filepath.Clean(targetPath) != strings.TrimSuffix(destClean, string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", targetPath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, f.Mode()); err != nil {
				return fmt.Errorf("mkdir %s: %w", targetPath, err)
			}
			continue
		}

		// Make parent dirs
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return fmt.Errorf("mkdir parent %s: %w", filepath.Dir(targetPath), err)
		}

		inFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("open member %s: %w", f.Name, err)
		}

		outFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			inFile.Close()
			return fmt.Errorf("create file %s: %w", targetPath, err)
		}

		_, err = io.Copy(outFile, inFile)
		inFile.Close()
		outFile.Close()
		if err != nil {
			return fmt.Errorf("copy %s: %w", targetPath, err)
		}
	}

	return nil
}
