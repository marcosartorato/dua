package scan

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type FileEntry struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type Result struct {
	Root       string      `json:"root"`
	TotalSize  int64       `json:"total_size"`
	TotalFiles int64       `json:"total_files"`
	TotalDirs  int64       `json:"total_dirs"`
	Files      []FileEntry `json:"top_files,omitempty"`
	Generated  time.Time   `json:"generated_at"`
}

type Options struct {
	TopN         int
	IncludeFiles bool
}

// Run scans the directory tree rooted at root and returns disk usage statistics.
// It returns the total size in bytes, total number of files, and total number of directories.
func Run(
	ctx context.Context, root string, opts Options, nowFunction func() time.Time,
) (Result, []string, error) {
	// Ensure root exists
	if _, err := os.Stat(root); err != nil {
		return Result{}, nil, err
	}

	var (
		totalSize  int64 // in Bytes
		totalFiles int64
		totalDirs  int64
		warnings   []string
		files      []FileEntry
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
		size := info.Size()
		totalSize += size
		totalFiles++

		/*
			Append to files if needed.
			This is a simple implementation; in a real scenario, you might want
			to keep only the top N largest files making it slower but less
			memory intensive.
		*/
		if opts.IncludeFiles {
			files = append(files, FileEntry{Path: path, Size: size})
		}

		return nil
	})
	if err != nil {
		return Result{}, warnings, err
	}

	// If IncludeFiles is set, sort files by size descending
	if opts.IncludeFiles {
		sort.Slice(files, func(i, j int) bool { return files[i].Size > files[j].Size })

		// If TopN is set, truncate the slice
		if opts.TopN > 0 && len(files) > opts.TopN {
			files = files[:opts.TopN]
		}
	}

	res := Result{
		Root:       root,
		TotalSize:  totalSize,
		TotalFiles: totalFiles,
		TotalDirs:  totalDirs,
		Files:      files,
		Generated:  nowFunction(),
	}
	return res, warnings, nil

}
