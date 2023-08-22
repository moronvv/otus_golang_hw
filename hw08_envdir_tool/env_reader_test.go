package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("dir not found", func(t *testing.T) {
		_, err := ReadDir("./not_exist")
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("not dir", func(t *testing.T) {
		_, err := ReadDir("./testdata/env/FOO")
		require.ErrorIs(t, err, ErrNotDirectory)
	})

	t.Run("valid", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env/")
		require.NoError(t, err)

		expectedEnvs := Environment{
			"FOO": EnvValue{
				Value: "   foo\nwith new line",
			},
			"BAR": EnvValue{
				Value: "bar",
			},
			"HELLO": EnvValue{
				Value: "\"hello\"",
			},
			"EMPTY": EnvValue{
				Value: "",
			},
			"UNSET": EnvValue{
				NeedRemove: true,
			},
		}
		require.Equal(t, expectedEnvs, envs)
	})
}
