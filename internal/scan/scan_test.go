package scan_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/marcosartorato/dua/internal/scan"
)

func TestRun_BasicTotals(t *testing.T) {
	t.Parallel()

	// Use a temporary directory.
	root := t.TempDir()

	// Create structure:
	// root/
	//   a.txt (10 bytes)
	//   sub/
	//     b.txt (5 bytes)
	if err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("1234567890"), 0644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(root, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "b.txt"), []byte("12345"), 0644); err != nil {
		t.Fatal(err)
	}

	fixedNow := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	nowFn := func() time.Time { return fixedNow }

	// Run scan.
	res, warnings, err := scan.Run(context.Background(), root, scan.Options{}, nowFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}

	// Check totals.
	if res.TotalFiles != 2 {
		t.Fatalf("TotalFiles = %d, want 2", res.TotalFiles)
	}
	if res.TotalDirs != 2 { // root + sub
		t.Fatalf("TotalDirs = %d, want 2", res.TotalDirs)
	}
	if res.TotalSize != 15 {
		t.Fatalf("TotalSize = %d, want 15", res.TotalSize)
	}
	if !res.Generated.Equal(fixedNow) {
		t.Fatalf("Generated = %v, want %v", res.Generated, fixedNow)
	}
}

func TestRun_EmptyDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	fixedNow := time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)
	nowFn := func() time.Time { return fixedNow }

	res, warnings, err := scan.Run(context.Background(), root, scan.Options{}, nowFn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", warnings)
	}

	if res.TotalFiles != 0 {
		t.Fatalf("TotalFiles = %d, want 0", res.TotalFiles)
	}
	if res.TotalDirs != 1 { // root only
		t.Fatalf("TotalDirs = %d, want 1", res.TotalDirs)
	}
	if res.TotalSize != 0 {
		t.Fatalf("TotalSize = %d, want 0", res.TotalSize)
	}
	if !res.Generated.Equal(fixedNow) {
		t.Fatalf("Generated = %v, want %v", res.Generated, fixedNow)
	}
}

func TestRun_InvalidRoot(t *testing.T) {
	t.Parallel()

	nowFn := func() time.Time { return time.Unix(0, 0).UTC() }

	_, _, err := scan.Run(context.Background(), "/path/does/not/exist", scan.Options{}, nowFn)
	if err == nil {
		t.Fatalf("expected error for invalid root, got nil")
	}
}

func TestRun_ContextCancel(t *testing.T) {
	// This test should not be run in parallel because it uses context cancellation.
	// filepath.WalkDir may not handle concurrent cancellations well.

	root := t.TempDir()

	// Create many files to give WalkDir something to do
	for i := 0; i < 2000; i++ {
		name := filepath.Join(root, fmt.Sprintf("f%06d.txt", i))
		_ = os.WriteFile(name, []byte("data"), 0644)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	nowFn := func() time.Time { return time.Unix(0, 0).UTC() }

	_, _, err := scan.Run(ctx, root, scan.Options{}, nowFn)
	if err == nil {
		t.Fatalf("expected context cancellation error, got nil")
	}
}
