package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrToAndFromAreSameFiles = errors.New("to and from must be different files")
)

// Check files edge cases
func checkFiles(fromPath, toPath string, offset int64, limit *int64) error {
	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	toFileInfo, err := os.Stat(toPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else {
		// check from and to are not the same file
		if os.SameFile(fromFileInfo, toFileInfo) {
			return ErrToAndFromAreSameFiles
		}
	}

	fileSize := fromFileInfo.Size()
	// incorrect file size
	if fileSize == 0 {
		return ErrUnsupportedFile
	}

	// offset must be less then file size
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	// correct limit
	if *limit == 0 {
		*limit = fileSize
	}

	return nil
}

// Get src file.
func getSrcFile(path string) (*os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Get dst file. Create dst path if not exists.
func getDstFile(path string) (*os.File, error) {
	file, err := os.Create(path)

	// if dir not exists
	if os.IsNotExist(err) {
		// create path automatically
		dir, _ := filepath.Split(path)
		if err := os.MkdirAll(dir, os.ModeDir); err != nil {
			return nil, err
		}

		// try to create file again
		file, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}

	return file, err
}

// Copy from one file to another with offset/limit.
func copyContent(fromFile, toFile *os.File, offset, limit int64) error {
	if _, err := fromFile.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	bar := pb.Full.Start64(limit)
	defer bar.Finish()
	fromFileWithBar := bar.NewProxyReader(fromFile)

	if _, err := io.CopyN(toFile, fromFileWithBar, limit); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := checkFiles(fromPath, toPath, offset, &limit); err != nil {
		return err
	}

	fromFile, err := getSrcFile(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := getDstFile(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	return copyContent(fromFile, toFile, offset, limit)
}
