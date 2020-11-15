package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Success with testdata", func(t *testing.T) {
		expectEnv := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: `"hello"`},
			"UNSET": EnvValue{Remove: true},
		}
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)
		require.Equal(t, env, expectEnv)
	})

	t.Run("Success", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		// файл, который должен быть проигнорирован
		err = ioutil.WriteFile(filepath.Join(dir, "t=t"), []byte("bar"), 0666)
		require.NoError(t, err)
		// имя файла маленькими буквами
		err = ioutil.WriteFile(filepath.Join(dir, "test"), []byte("test"), 0666)
		require.NoError(t, err)
		// файл с пустой первой строкой
		err = ioutil.WriteFile(filepath.Join(dir, "EMPTY"), []byte("\n"), 0666)
		require.NoError(t, err)

		expectEnv := Environment{
			"test":  EnvValue{Value: "test", Remove: false},
			"EMPTY": EnvValue{Value: "", Remove: false},
		}
		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, env, expectEnv)
	})

	t.Run("Success with empty dir", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Len(t, env, 0)
	})

	t.Run("Fail with dir not exists", func(t *testing.T) {
		_, err := ReadDir("some name")
		require.Error(t, err)
	})
}
