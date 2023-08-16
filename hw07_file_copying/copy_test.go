package main

import (
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	const (
		inPath = "./testdata/input.txt"
	)

	t.Run("from and to are same files", func(t *testing.T) {
		err := Copy(inPath, inPath, 0, 0)
		require.ErrorIs(t, err, ErrToAndFromAreSameFiles)
	})

	t.Run("offset > file length", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "output.txt")
		require.NoError(t, err)
		outPath := outFile.Name()
		defer os.Remove(outPath)

		err = Copy(inPath, outPath, 10000, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("unsupported file", func(t *testing.T) {
		outFile, err := os.CreateTemp("", "output.txt")
		require.NoError(t, err)
		outPath := outFile.Name()
		defer os.Remove(outPath)

		err = Copy("/dev/urandom", outPath, 0, 0)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	testCases := []struct {
		name     string
		testPath string
		offset   int64
		limit    int64
	}{
		{
			name:     "input -> out_offset0_limit0",
			testPath: "./testdata/out_offset0_limit0.txt",
			offset:   0,
			limit:    0,
		},
		{
			name:     "input -> out_offset0_limit10",
			testPath: "./testdata/out_offset0_limit10.txt",
			offset:   0,
			limit:    10,
		},
		{
			name:     "input -> out_offset0_limit1000",
			testPath: "./testdata/out_offset0_limit1000.txt",
			offset:   0,
			limit:    1000,
		},
		{
			name:     "input -> out_offset0_limit10000",
			testPath: "./testdata/out_offset0_limit10000.txt",
			offset:   0,
			limit:    10000,
		},
		{
			name:     "input -> out_offset100_limit1000",
			testPath: "./testdata/out_offset100_limit1000.txt",
			offset:   100,
			limit:    1000,
		},
		{
			name:     "input -> out_offset6000_limit1000",
			testPath: "./testdata/out_offset6000_limit1000.txt",
			offset:   6000,
			limit:    1000,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// create temp file
			outFile, err := os.CreateTemp("", "output.txt")
			require.NoError(t, err)
			outPath := outFile.Name()
			defer os.Remove(outPath)

			// run copy
			err = Copy(inPath, outPath, tc.offset, tc.limit)
			require.NoError(t, err)

			// get dst file md5
			outFileHash := md5.New()
			_, err = io.Copy(outFileHash, outFile)
			require.NoError(t, err)

			// get test file md5
			testFile, err := os.Open(tc.testPath)
			require.NoError(t, err)
			defer testFile.Close()
			testFileHash := md5.New()
			_, err = io.Copy(testFileHash, testFile)
			require.NoError(t, err)

			// compare md5 of test and dst files
			require.True(t, bytes.Equal(testFileHash.Sum(nil), outFileHash.Sum(nil)))
		})
	}
}
