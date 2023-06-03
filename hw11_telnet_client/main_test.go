package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommandLine(t *testing.T) {
	str, err := parseCommandLine([]string{}, true)
	require.Equal(t, "", str)
	require.NotNil(t, err)

	str, err = parseCommandLine([]string{"go-telnet", "127.0.0.1", "4242"}, true)
	require.Equal(t, "", str)
	require.NotNil(t, err)

	str, err = parseCommandLine([]string{"go-telnet", "timeout", "127.0.0.1", "4242"}, true)
	require.Equal(t, "127.0.0.1:4242", str)
	require.Nil(t, err)

	str, err = parseCommandLine([]string{"go-telnet", "timeout", "127.0.0.1", "4242"}, false)
	require.Equal(t, "", str)
	require.NotNil(t, err)

	str, err = parseCommandLine([]string{"go-telnet", "127.0.0.1", "4242"}, false)
	require.Equal(t, "127.0.0.1:4242", str)
	require.Nil(t, err)
}
