package main

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

var ErrNotDirectory = errors.New("input path is not directory")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// Check input dir string and get reader.
func getDir(dirPath string) (*os.File, error) {
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return nil, err
	}
	if !dirInfo.IsDir() {
		return nil, ErrNotDirectory
	}

	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	return dir, nil
}

// Process one file to get env info.
func getEnvValue(filePath string) (EnvValue, error) {
	var envValue EnvValue

	file, err := os.Open(filePath)
	if err != nil {
		return envValue, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// read first line only
	scanner.Scan()
	value := scanner.Text()

	if len(value) == 0 {
		envValue.NeedRemove = true
	} else {
		// trim ending
		value = strings.TrimRightFunc(value, unicode.IsSpace)

		// replace terminal nulls
		value = strings.ReplaceAll(value, "\x00", "\n")

		envValue.Value = value
	}

	return envValue, nil
}

// Process all files in dir.
func getEnvs(dirPath string, files []fs.DirEntry) (Environment, error) {
	envs := Environment{}

	for _, file := range files {
		fileName := file.Name()
		// skip files with '=' in name
		if !strings.Contains(fileName, "=") {
			filePath := filepath.Join(dirPath, fileName)

			envValue, err := getEnvValue(filePath)
			if err != nil {
				return nil, err
			}

			envs[file.Name()] = envValue
		}
	}

	return envs, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dirPath string) (Environment, error) {
	dir, err := getDir(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	return getEnvs(dirPath, files)
}
