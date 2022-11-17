package main

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	inFile := "testdata/input.txt"
	cmpFileFmt := "testdata/out_offset%d_limit%d.txt"

	testCases := []struct {
		Limit, Offset int64
	}{
		{
			Offset: 0,
			Limit:  0,
		}, {
			Offset: 0,
			Limit:  10,
		}, {
			Offset: 0,
			Limit:  1000,
		}, {
			Offset: 0,
			Limit:  10000,
		}, {
			Offset: 100,
			Limit:  1000,
		}, {
			Offset: 6000,
			Limit:  1000,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d/%d", tc.Offset, tc.Limit), func(t *testing.T) {
			outFile, err := os.CreateTemp("", "out")
			require.NoError(t, err)

			defer func() {
				_ = os.Remove(outFile.Name())
			}()

			err = Copy(inFile, outFile.Name(), tc.Offset, tc.Limit)
			require.NoError(t, err)

			out, err := os.ReadFile(outFile.Name())
			require.NoError(t, err)

			expected, err := os.ReadFile(fmt.Sprintf(cmpFileFmt, tc.Offset, tc.Limit))
			require.NoError(t, err)

			require.Equal(t, expected, out)
		})
	}

	t.Run("empty in", func(t *testing.T) {
		err := Copy("", "", 0, 0)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("empty out", func(t *testing.T) {
		err := Copy(inFile, "", 0, 0)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("not exists", func(t *testing.T) {
		err := Copy("not_exists", "", 0, 0)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("not regular", func(t *testing.T) {
		err := Copy("/dev/urandom", "", 0, 0)
		require.True(t, errors.Is(err, ErrUnsupportedFile))
	})

	t.Run("src is directory", func(t *testing.T) {
		err := Copy(os.TempDir(), "", 0, 0)
		require.True(t, errors.Is(err, ErrUnsupportedFile))
	})

	t.Run("wrong offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "", 10000, 0)
		require.True(t, errors.Is(err, ErrOffsetExceedsFileSize))
	})
}
