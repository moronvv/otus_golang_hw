package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	envs, err := ReadDir("./testdata/env/")
	require.NoError(t, err)

	t.Run("too few args", func(t *testing.T) {
		code := RunCmd([]string{}, envs)
		require.Equal(t, 1, code)
	})

	t.Run("cmd not exists", func(t *testing.T) {
		code := RunCmd([]string{"./testdata/env/BAR", "arg=1"}, envs)
		require.Equal(t, 1, code)
	})

	t.Run("got ExitErr", func(t *testing.T) {
		code := RunCmd([]string{"sleep", "-u"}, envs)
		require.Equal(t, 1, code)
	})

	t.Run("valid", func(t *testing.T) {
		code := RunCmd([]string{"./testdata/echo.sh", "arg=1"}, envs)
		require.Equal(t, 0, code)
	})
}
