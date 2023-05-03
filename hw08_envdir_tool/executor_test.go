package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("test without env", func(t *testing.T) {
		file, _ := os.Create("testdata/test.sh")
		defer os.Remove("testdata/test.sh")

		text := `#!/usr/bin/env bash
		echo test script`
		file.WriteString(text)

		code := RunCmd([]string{"sh", "testdata/test.sh"}, nil)
		require.Equal(t, 0, code)
	})

	t.Run("test with non zero code", func(t *testing.T) {
		file, _ := os.Create("testdata/test.sh")
		defer os.Remove("testdata/test.sh")

		text := `#!/usr/bin/env bash
exit 1`
		file.WriteString(text)

		code := RunCmd([]string{"sh", "testdata/test.sh"}, nil)
		require.Equal(t, 1, code)
	})

	t.Run("test without arguments", func(t *testing.T) {
		file, _ := os.Create("testdata/test.sh")
		defer os.Remove("testdata/test.sh")

		text := `#!/usr/bin/env bash
exit 1`
		file.WriteString(text)

		code := RunCmd([]string{"sh"}, nil)
		require.Equal(t, 0, code)
	})
}

func TestSetEnvironments(t *testing.T) {
	t.Run("test set evnironment variable", func(t *testing.T) {
		environment := make(Environment)
		environment["TEST"] = EnvValue{Value: "Test value", NeedRemove: false}

		script, _ := os.Create("testdata/script.sh")
		text := `#!/usr/bin/env bash

echo -n ${TEST}`

		script.WriteString(text)

		out := os.Stdout
		test, _ := os.Create("testdata/testOut")
		os.Stdout = test
		defer func() {
			os.Stdout = out
			os.Remove("testdata/testOut")
			os.Remove("testdata/script.sh")
		}()

		code := RunCmd([]string{"sh", "testdata/script.sh"}, environment)
		content, _ := os.ReadFile(os.Stdout.Name())
		require.Equal(t, 0, code)
		require.Equal(t, "Test value", string(content))
	})

	t.Run("test unset evnironment variable", func(t *testing.T) {
		os.Setenv("TEST", "Test value")
		environment := make(Environment)
		environment["TEST"] = EnvValue{Value: "New Value", NeedRemove: true}

		script, _ := os.Create("testdata/script.sh")
		text := `#!/usr/bin/env bash

echo -n ${TEST}`

		script.WriteString(text)

		out := os.Stdout
		test, _ := os.Create("testdata/testOut")
		os.Stdout = test
		defer func() {
			os.Stdout = out
			os.Remove("testdata/testOut")
			os.Remove("testdata/script.sh")
		}()

		code := RunCmd([]string{"sh", "testdata/script.sh"}, environment)
		content, _ := os.ReadFile(os.Stdout.Name())
		require.Equal(t, 0, code)
		require.Equal(t, "", string(content))
	})
}
