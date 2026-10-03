package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestEnsureDatabaseCopiesSeedWhenMissing(t *testing.T) {
	dir := t.TempDir()
	seed := filepath.Join(dir, "seed.db")
	live := filepath.Join(dir, "live.db")
	writeFile(t, seed, "seed contents")

	if err := ensureDatabase(live, seed); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, live); got != "seed contents" {
		t.Errorf("live db = %q, want seed contents", got)
	}
}

func TestEnsureDatabaseKeepsExistingDatabase(t *testing.T) {
	dir := t.TempDir()
	seed := filepath.Join(dir, "seed.db")
	live := filepath.Join(dir, "live.db")
	writeFile(t, seed, "seed contents")
	writeFile(t, live, "votes since last deploy")

	if err := ensureDatabase(live, seed); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, live); got != "votes since last deploy" {
		t.Errorf("existing db was overwritten: got %q", got)
	}
}

func TestEnsureDatabaseFailsWithoutSeed(t *testing.T) {
	dir := t.TempDir()
	live := filepath.Join(dir, "live.db")

	if err := ensureDatabase(live, filepath.Join(dir, "missing.db")); err == nil {
		t.Fatal("expected an error when the seed file is missing")
	}
	if _, err := os.Stat(live); !os.IsNotExist(err) {
		t.Errorf("live db should not exist after a failed seed, stat err = %v", err)
	}
}

func TestEnsureDatabaseLeavesNoTempFiles(t *testing.T) {
	dir := t.TempDir()
	seed := filepath.Join(dir, "seed.db")
	live := filepath.Join(dir, "live.db")
	writeFile(t, seed, "seed contents")

	if err := ensureDatabase(live, seed); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		names := []string{}
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Errorf("want only seed.db and live.db, got %v", names)
	}
}
