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

func TestRun_TotalsOnly(t *testing.T) {
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
	res, warnings, err := scan.Run(context.Background(), root, scan.Options{}, nowFn) // TopN ignored
	require.NoError(t, err)
	assert.Empty(t, warnings)

	// Check totals.
	assert.Equal(t, int64(2), res.TotalFiles)
	assert.Equal(t, int64(2), res.TotalDirs) // root + sub
	assert.Equal(t, int64(15), res.TotalSize)
	assert.True(t, res.Generated.Equal(fixedNow))
}

func TestRun_IncludeFiles_CollectAll(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	// sizes: 1, 2, 3 => total 6
	require.NoError(t, os.WriteFile(filepath.Join(root, "f1"), []byte("a"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "f2"), []byte("bb"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "f3"), []byte("ccc"), 0644))

	fixedNow := time.Unix(123, 0).UTC()

	res, warnings, err := scan.Run(
		context.Background(),
		root,
		scan.Options{IncludeFiles: true, TopN: 0}, // 0 => keep all
		func() time.Time { return fixedNow },
	)

	require.NoError(t, err)
	assert.Empty(t, warnings)

	assert.Equal(t, int64(6), res.TotalSize)
	assert.Equal(t, int64(3), res.TotalFiles)
	assert.Equal(t, int64(1), res.TotalDirs) // root only
	assert.Len(t, res.Files, 3)
	assert.Equal(t, fixedNow, res.Generated)

	// Verify that the returned set contains all paths (order not guaranteed when TopN==0)
	want := map[string]int64{
		filepath.Join(root, "f1"): 1,
		filepath.Join(root, "f2"): 2,
		filepath.Join(root, "f3"): 3,
	}
	got := map[string]int64{}
	for _, e := range res.Files {
		got[e.Path] = e.Size
	}
	assert.Equal(t, want, got)
}

func TestRun_IncludeFiles_TopN(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	// sizes: 10, 1, 7, 3  => Top 2 should be 10 and 7
	require.NoError(t, os.WriteFile(filepath.Join(root, "a"), []byte("1234567890"), 0644)) // 10
	require.NoError(t, os.WriteFile(filepath.Join(root, "b"), []byte("1"), 0644))          // 1
	require.NoError(t, os.WriteFile(filepath.Join(root, "c"), []byte("1234567"), 0644))    // 7
	require.NoError(t, os.WriteFile(filepath.Join(root, "d"), []byte("123"), 0644))        // 3

	topN := 2
	opt := scan.Options{IncludeFiles: true, TopN: topN}
	res, warnings, err := scan.Run(
		context.Background(),
		root,
		opt,
		time.Now,
	)

	require.NoError(t, err)
	assert.Empty(t, warnings)

	require.Len(t, res.Files, topN)
	assert.GreaterOrEqual(t, res.Files[0].Size, res.Files[1].Size) // sorted desc

	assert.Equal(t, int64(10), res.Files[0].Size)
	assert.Equal(t, int64(7), res.Files[1].Size)

	// Sanity check totals
	assert.Equal(t, int64(21), res.TotalSize)
	assert.Equal(t, int64(4), res.TotalFiles)
	assert.Equal(t, int64(1), res.TotalDirs)
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

func TestRun_TopNIgnoredWhenIncludeFilesFalse(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "a"), []byte("123"), 0644))

	res, warnings, err := scan.Run(
		context.Background(),
		root,
		scan.Options{IncludeFiles: false, TopN: 1},
		time.Now,
	)

	require.NoError(t, err)
	assert.Empty(t, warnings)
	assert.Empty(t, res.Files)
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
