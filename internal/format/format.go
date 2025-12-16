package format

import (
	"fmt"
	"io"

	"github.com/marcosartorato/dua/internal/scan"
)

// PrintTable prints the scan.Result in a human-readable table format to the provided writer.
// It returns an error if the result type is unsupported.
func PrintTable(w io.Writer, res interface{}) error {
	// Assert type.
	r, ok := res.(scan.Result)
	if !ok {
		return fmt.Errorf("PrintTable: unsupported result type %T", res)
	}

	// Print table.
	if _, err := fmt.Fprintf(w, "Path:   %s\n", r.Root); err != nil {
		return fmt.Errorf("PrintTable: failed to write path: %w", err)
	}
	if _, err := fmt.Fprintf(w, "Size:   %d bytes\n", r.TotalSize); err != nil {
		return fmt.Errorf("PrintTable: failed to write size: %w", err)
	}
	if _, err := fmt.Fprintf(w, "Files:  %d\n", r.TotalFiles); err != nil {
		return fmt.Errorf("PrintTable: failed to write files: %w", err)
	}
	if _, err := fmt.Fprintf(w, "Dirs:   %d\n", r.TotalDirs); err != nil {
		return fmt.Errorf("PrintTable: failed to write dirs: %w", err)
	}

	return nil
}
