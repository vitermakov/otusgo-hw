package main

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	validData := map[string]string{
		"TEST1":     "VALUE1",
		"TEST2":     "VALUE2",
		"UNSET":     "",
		"EMPTY":     "    \t",
		"WITHNEWLN": string([]byte{'f', 'o', 'o', '\x00', 'b', 'a', 'r'}),
		"MANYLINES": "LINE1\nLINE2\nLINE3\nLINE4",
	}
	expected := Environment{
		"TEST1": EnvValue{
			Value: "VALUE1",
		},
		"TEST2": EnvValue{
			Value: "VALUE2",
		},
		"UNSET": EnvValue{
			NeedRemove: true,
		},
		"EMPTY": EnvValue{},
		"WITHNEWLN": EnvValue{
			Value: "foo\nbar",
		},
		"MANYLINES": EnvValue{
			Value: "LINE1",
		},
	}

	testPath := prepareTestPath(t, validData)
	actual, err := ReadDir(testPath)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	err = os.RemoveAll(testPath)
	require.NoError(t, err)
}

func TestReadNoDir(t *testing.T) {
	testPath := filepath.Join(os.TempDir(), "/hw08-not-exists")
	_, err := ReadDir(testPath)

	var pe *os.PathError
	require.True(t, errors.As(err, &pe))
}

func prepareTestPath(t *testing.T, validData map[string]string) string {
	t.Helper()

	var (
		err   error
		fName string
	)

	testPath := filepath.Join(os.TempDir(), "/hw08")
	subPath := filepath.Join(testPath, "/subdir")

	_ = os.RemoveAll(testPath)

	err = os.Mkdir(testPath, os.ModePerm)
	require.NoError(t, err)

	err = os.Mkdir(subPath, os.ModePerm)
	require.NoError(t, err)

	fName = filepath.Join(testPath, "/wr;ong")
	os.WriteFile(fName, []byte{'w', 'r', 'o', 'n', 'g'}, os.ModePerm)

	fName = filepath.Join(testPath, "/wr=ong")
	os.WriteFile(fName, []byte{'w', 'r', 'o', 'n', 'g'}, os.ModePerm)

	for fName, value := range validData {
		fName = filepath.Join(testPath, "/"+fName)
		os.WriteFile(fName, []byte(value), os.ModePerm)
	}

	return testPath
}
