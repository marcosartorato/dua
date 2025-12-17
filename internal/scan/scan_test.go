package scan_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/marcosartorato/dua/internal/scan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, os.WriteFile(filepath.Join(root, "a.txt"), []byte("1234567890"), 0644))
	sub := filepath.Join(root, "sub")
	require.NoError(t, os.Mkdir(sub, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(sub, "b.txt"), []byte("12345"), 0644))
	fixedNow := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	nowFn := func() time.Time { return fixedNow }

	// Run scan.
	res, warnings, err := scan.Run(context.Background(), root, scan.Options{}, nowFn)
	require.NoError(t, err)
	assert.Empty(t, warnings)

	// Check totals.
	assert.Equal(t, int64(2), res.TotalFiles)
	assert.Equal(t, int64(2), res.TotalDirs) // root + sub
	assert.Equal(t, int64(15), res.TotalSize)
	assert.True(t, res.Generated.Equal(fixedNow))
}

func TestRun_EmptyDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	fixedNow := time.Date(2024, 2, 2, 0, 0, 0, 0, time.UTC)
	nowFn := func() time.Time { return fixedNow }

	res, warnings, err := scan.Run(context.Background(), root, scan.Options{}, nowFn)

	require.NoError(t, err)
	assert.Empty(t, warnings)

	assert.Equal(t, int64(0), res.TotalFiles)
	assert.Equal(t, int64(1), res.TotalDirs) // root only
	assert.Equal(t, int64(0), res.TotalSize)
	assert.Equal(t, fixedNow, res.Generated)
}

func TestRun_InvalidRoot(t *testing.T) {
	t.Parallel()

	nowFn := func() time.Time { return time.Unix(0, 0).UTC() }

	_, _, err := scan.Run(context.Background(), "/path/does/not/exist", scan.Options{}, nowFn)
	assert.Error(t, err)
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
	assert.Error(t, err)
}
