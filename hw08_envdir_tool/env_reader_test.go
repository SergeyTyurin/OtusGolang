package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("correct directory path case", func(t *testing.T) {
		paths := []string{"testdata/env/", "testdata/env"}
		for _, path := range paths {
			env, err := ReadDir(path)
			require.Nil(t, err)
			require.Equal(t, 5, len(env))

			require.Equal(t, EnvValue{Value: "bar", NeedRemove: false}, env["BAR"])
			require.Equal(t, EnvValue{Value: "", NeedRemove: false}, env["EMPTY"])
			require.Equal(t, EnvValue{Value: "   foo\nwith new line", NeedRemove: false}, env["FOO"])
			require.Equal(t, EnvValue{Value: `"hello"`, NeedRemove: false}, env["HELLO"])
			require.Equal(t, EnvValue{Value: "", NeedRemove: true}, env["UNSET"])
		}
	})

	t.Run("incorrect directory path case", func(t *testing.T) {
		env, err := ReadDir("/path/to/file")
		require.Nil(t, env, "Environment variables aren't nil")
		require.ErrorIs(t, err, ErrNoSuchDir)

		env, err = ReadDir("testdata/env/BAR")
		require.Nil(t, env, "Environment variables aren't nil")
		require.ErrorIs(t, err, ErrIsNotDir)
	})

	t.Run("incorrect filename case", func(t *testing.T) {
		incorrectFile, _ := os.Create("testdata/env/TEST=INCORRECT")
		defer os.Remove("testdata/env/TEST=INCORRECT")
		incorrectFile.WriteString("test")

		env, err := ReadDir("testdata/env/")
		require.Nil(t, err)
		require.Equal(t, 5, len(env))

		_, ok := env["TEST=INCORRECT"]
		require.False(t, ok)
	})
}
