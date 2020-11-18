package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

// пустой файл и файл с пустой первой строкой должны обрабатываться по разному,
// предложенный в задании тип Environment для этого не подходит.
type EnvValue struct {
	Value  string
	Remove bool
}

type Environment map[string]EnvValue

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make(Environment)
	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			continue
		}
		if !file.Mode().IsRegular() {
			continue
		}

		if file.Size() == 0 {
			result[file.Name()] = EnvValue{Remove: true}
			continue
		}

		value, err := readValue(dir, file.Name())
		if err != nil {
			return nil, err
		}
		result[file.Name()] = EnvValue{Value: value}
	}

	return result, nil
}

func readValue(dir, fileName string) (string, error) {
	f, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return "", err
	}
	defer func() { mustNil(f.Close()) }()

	s := bufio.NewScanner(f)
	if !s.Scan() {
		return "", nil
	}
	line := s.Text()
	line = strings.ReplaceAll(line, "\x00", "\n")
	line = strings.TrimRightFunc(line, unicode.IsSpace)
	return line, nil
}

func mustNil(err error) {
	if err != nil {
		panic(err)
	}
}
