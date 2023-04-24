package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	defer os.RemoveAll("results")
	require.ErrorIs(t, Copy("", "results/res.txt", 0, 100), ErrUnsupportedFile)
	require.ErrorIs(t, Copy("testdata/input.txt", "", 0, 100), ErrUnsupportedFile)
	require.ErrorIs(t, Copy("testdata/", "testdata/input.txt", 0, 100), ErrIsDirectory)
	require.ErrorIs(t, Copy("testdata/input.txt", "results/out.txt", 10000, 100), ErrOffsetExceedsFileSize)
	require.ErrorIs(t, Copy("testdata/input.txt", "testdata/input.txt", 0, 100), ErrIsSamePaths)
}
