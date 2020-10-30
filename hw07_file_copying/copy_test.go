package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNegativeOffset(t *testing.T) {
	err := Copy("testdata/input.txt", "out.txt", -10, 0)

	require.Error(t, err)
}

func TestNegativeLimit(t *testing.T) {
	err := Copy("testdata/input.txt", "out.txt", 0, -1)

	require.Error(t, err)
}

func TestNoFromFile(t *testing.T) {
	err := Copy("filename", "out.txt", 0, 0)

	require.Error(t, err)
}

func TestUnsupportedFile(t *testing.T) {
	err := Copy("/dev/urandom", "out.txt", 0, 0)

	require.Equal(t, true, errors.Is(err, ErrUnsupportedFile))
}

func TestOffsetExceedsFileSize(t *testing.T) {
	err := Copy("testdata/input.txt", "out.txt", 1000000, 0)

	require.Equal(t, true, errors.Is(err, ErrOffsetExceedsFileSize))
}

func TestNoPermFromFile(t *testing.T) {
	toFile, _ := os.OpenFile("fromFile", os.O_CREATE, 0)
	_ = toFile.Close()
	defer func() {
		_ = os.Remove("fromFile")
	}()

	err := Copy("fromFile", "out.txt", 0, 0)

	require.Error(t, err)
}

func TestNoPermToFile(t *testing.T) {
	err := Copy("testdata/input.txt", "/dev/urandom", 0, 0)

	require.Error(t, err)
}
