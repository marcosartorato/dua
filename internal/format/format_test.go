package format_test

import (
	"bytes"
	"testing"

	"github.com/marcosartorato/dua/internal/format"
	"github.com/marcosartorato/dua/internal/scan"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// Run PrintTable.
	err := format.PrintTable(&buf, res)
	require.NoError(t, err)

	// Check output.
	assert.Equal(t,
		"Path:   /tmp\n"+
			"Size:   12345 bytes\n"+
			"Files:  7\n"+
			"Dirs:   2\n",
		buf.String(),
	)
}

func TestPrintTable_WrongWriterType(t *testing.T) {
	// Use a bytes.Buffer to capture output.
	var buf bytes.Buffer

	// Run PrintTable.
	err := format.PrintTable(&buf, "not a scan.Result")

	// Check output.
	require.Error(t, err)
	assert.Empty(t, buf.String())
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

	require.Error(t, err)
	assert.Empty(t, buf.String())
}
