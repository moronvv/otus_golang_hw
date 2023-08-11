package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Get src file while checking edge cases.
func getSrcFile(path string, offset int64, limit *int64) (*os.File, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// offset must be less then file size
	if offset > fileSize {
		return nil, ErrOffsetExceedsFileSize
	}
	// TODO: check inf file like /dev/urandom

	// correct limit
	if *limit == 0 {
		*limit = fileSize
	}

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

	if _, err := io.CopyN(toFile, fromFile, limit); err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}
	}

	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, err := getSrcFile(fromPath, offset, &limit)
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
