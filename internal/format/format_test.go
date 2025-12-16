package format_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/marcosartorato/dua/internal/format"
	"github.com/marcosartorato/dua/internal/scan"
)

func TestPrintTable_OK(t *testing.T) {
	// Use a bytes.Buffer to capture output.
	var buf bytes.Buffer

	// Sample scan.Result and expected output lines.
	res := scan.Result{
		Root:       "/tmp",
		TotalSize:  12345,
		TotalFiles: 7,
		TotalDirs:  2,
	}
	wantLines := []string{
		"Path:   /tmp\n",
		"Size:   12345 bytes\n",
		"Files:  7\n",
		"Dirs:   2\n",
	}

	// Run PrintTable.
	if err := format.PrintTable(&buf, res); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Check output.
	got := buf.String()
	for _, line := range wantLines {
		if !strings.Contains(got, line) {
			t.Fatalf("output missing line %q\n--- got ---\n%s\n---", line, got)
		}
	}
}

func TestPrintTable_WrongWriterType(t *testing.T) {
	// Use a bytes.Buffer to capture output.
	var buf bytes.Buffer

	// Run PrintTable.
	err := format.PrintTable(&buf, "not a scan.Result")

	// Check output.
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output written on error, got: %q", buf.String())
	}
}

func TestPrintTable_WrongType(t *testing.T) {
	// PrintTable expects scan.Result (value), not *scan.Result
	var buf bytes.Buffer

	res := &scan.Result{
		Root:       "/tmp",
		TotalSize:  1,
		TotalFiles: 1,
		TotalDirs:  1,
	}

	err := format.PrintTable(&buf, res)
	if err == nil {
		t.Fatalf("expected error for *scan.Result, got nil")
	}
}
