package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func DeleteProjectDirectory(path string) error {
	root := os.Getenv("ROOT_PATH")
	if root == "" {
		return fmt.Errorf("ROOT_PATH environment variable is not set")
	}

	// Resolve absolute paths and clean ., .., etc.
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("failed to resolve ROOT_PATH: %w", err)
	}

	targetAbs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve target path: %w", err)
	}

	projectsRoot := filepath.Join(rootAbs, "projects")

	// 1. Must be inside $ROOT_PATH/projects
	if !strings.HasPrefix(targetAbs, projectsRoot+string(os.PathSeparator)) {
		return fmt.Errorf("path is not under allowed directory: %s", targetAbs)
	}

	// 2. Must be exactly one level below /projects
	rel, err := filepath.Rel(projectsRoot, targetAbs)
	if err != nil {
		return fmt.Errorf("failed to compute relative path: %w", err)
	}

	if rel == "." || strings.Contains(rel, string(os.PathSeparator)) {
		return fmt.Errorf("only /projects/{project} can be deleted, not: %s", targetAbs)
	}

	// 3. Extra safety: disallow symlinks
	info, err := os.Lstat(targetAbs)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to delete symlink")
	}

	// 4. Delete
	return os.RemoveAll(targetAbs)
}

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
