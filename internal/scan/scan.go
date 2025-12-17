package scan

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

type Result struct {
	Root       string    `json:"root"`
	TotalSize  int64     `json:"total_size"`
	TotalFiles int64     `json:"total_files"`
	TotalDirs  int64     `json:"total_dirs"`
	Generated  time.Time `json:"generated_at"`
}

// TODO support options like MaxDepth, etc.
type Options struct {
	MaxDepth int // -1 = unlimited
	TopN     int
}

// Run scans the directory tree rooted at root and returns disk usage statistics.
// It returns the total size in bytes, total number of files, and total number of directories.
func Run(ctx context.Context, root string, opts Options, nowFunction func() time.Time) (Result, []string, error) {
	// Ensure root exists
	if _, err := os.Stat(root); err != nil {
		return Result{}, nil, err
	}

	var (
		totalSize  int64
		totalFiles int64
		totalDirs  int64
		warnings   []string
	)

	// filepath.WalkDir is used to traverse the directory tree.
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		// Allow cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Handle errors from WalkDir
		if walkErr != nil {
			return walkErr
		}

		// Directory
		if d.IsDir() {
			totalDirs++
			return nil
		}

		// File (or symlink/etc.)
		info, err := d.Info()
		if err != nil {
			return err
		}
		totalSize += info.Size()
		totalFiles++
		return nil
	})
	if err != nil {
		return Result{}, warnings, err
	}

	res := Result{
		Root:       root,
		TotalSize:  totalSize,
		TotalFiles: totalFiles,
		TotalDirs:  totalDirs,
		Generated:  nowFunction(),
	}
	return res, warnings, nil

}
