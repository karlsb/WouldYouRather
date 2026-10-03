package main

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// seedDBPath is the database shipped in the image and used directly in local dev.
const seedDBPath = "./build-database/wouldyourather.db"

// ensureDatabase copies seedPath to path unless a database already exists there,
// so data on a persistent volume survives redeploys.
func ensureDatabase(path, seedPath string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	src, err := os.Open(seedPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Copy to a temp file and rename it into place, so an interrupted copy never
	// leaves a partial database that the next boot would treat as real.
	tmp, err := os.CreateTemp(filepath.Dir(path), ".seed-*.db")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name()) // no-op once renamed

	if _, err := io.Copy(tmp, src); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
