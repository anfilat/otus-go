package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// подготовка тестовых данных
		dir, err := ioutil.TempDir("", "test")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		// папка с переменными окружения
		err = os.Mkdir(filepath.Join(dir, "vars"), 0777)
		require.NoError(t, err)
		err = ioutil.WriteFile(filepath.Join(dir, "vars", "BAR"), []byte("bar"), 0666)
		require.NoError(t, err)
		// баш-скрипт (распечатывает свой первый аргумент и переменную окружения BAR)
		err = ioutil.WriteFile(filepath.Join(dir, "t.sh"), []byte("#!/usr/bin/env bash\necho $1\necho $BAR\n"), 0666)
		require.NoError(t, err)
		err = os.Chmod(filepath.Join(dir, "t.sh"), 0777)
		require.NoError(t, err)
		// конец подготовки тестовых данных

		env, err := ReadDir(filepath.Join(dir, "vars"))
		require.NoError(t, err)

		var returnCode int
		result := capturer.CaptureStdout(func() {
			returnCode = RunCmd([]string{filepath.Join(dir, "t.sh"), "something"}, env)
		})
		require.Equal(t, 0, returnCode)
		require.Equal(t, "something\nbar\n", result)
	})

	t.Run("Rewrite HOME", func(t *testing.T) {
		// подготовка тестовых данных
		dir, err := ioutil.TempDir("", "test")
		require.NoError(t, err)
		defer os.RemoveAll(dir)

		// папка с переменными окружения
		err = os.Mkdir(filepath.Join(dir, "vars"), 0777)
		require.NoError(t, err)
		err = ioutil.WriteFile(filepath.Join(dir, "vars", "HOME"), []byte("42"), 0666)
		require.NoError(t, err)
		// баш-скрипт (распечатывает переменную окружения HOME)
		err = ioutil.WriteFile(filepath.Join(dir, "t.sh"), []byte("#!/usr/bin/env bash\necho $HOME\n"), 0666)
		require.NoError(t, err)
		err = os.Chmod(filepath.Join(dir, "t.sh"), 0777)
		require.NoError(t, err)
		// конец подготовки тестовых данных

		env, err := ReadDir(filepath.Join(dir, "vars"))
		require.NoError(t, err)

		var returnCode int
		result := capturer.CaptureStdout(func() {
			returnCode = RunCmd([]string{filepath.Join(dir, "t.sh")}, env)
		})
		require.Equal(t, 0, returnCode)
		require.Equal(t, "42\n", result)
	})

	t.Run("Fail with wrong command", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		require.NoError(t, err)

		var returnCode int
		result := capturer.CaptureStderr(func() {
			returnCode = RunCmd([]string{"cat", "."}, env)
		})
		require.Equal(t, 1, returnCode)
		require.Equal(t, "cat: .: Is a directory\n", result)
	})
}
