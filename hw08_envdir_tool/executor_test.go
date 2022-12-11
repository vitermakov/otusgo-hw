package main

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	testCases := []struct {
		name         string
		command      []string
		env          Environment
		expectedCode int

		checkFunc func(t *testing.T, logged string)
	}{
		{
			name:         "nil args",
			command:      nil,
			env:          nil,
			expectedCode: -1,
			checkFunc: func(t *testing.T, logged string) {
				t.Helper()
				require.Contains(t, logged, "cmd not set")
			},
		}, {
			name:         "empty cmd 1",
			command:      []string{},
			env:          nil,
			expectedCode: -1,
			checkFunc: func(t *testing.T, logged string) {
				t.Helper()
				require.Contains(t, logged, "cmd not set")
			},
		}, {
			name:         "error cmd 2",
			command:      []string{""},
			env:          nil,
			expectedCode: -1,
			checkFunc: func(t *testing.T, logged string) {
				t.Helper()
				require.True(t, len(logged) > 0)
			},
		}, {
			name:         "error cmd",
			command:      []string{"cool_program"},
			env:          nil,
			expectedCode: -1,
			checkFunc: func(t *testing.T, logged string) {
				t.Helper()
				require.True(t, len(logged) > 0)
			},
		}, {
			name:         "test ok",
			command:      []string{"test", "100", "-eq", "100"},
			env:          nil,
			expectedCode: 0,
			checkFunc: func(t *testing.T, logged string) {
				t.Helper()
				require.True(t, len(logged) == 0)
			},
		}, {
			name:         "test failed",
			command:      []string{"test", "100", "-eq", "0"},
			env:          nil,
			expectedCode: 1,
			checkFunc: func(t *testing.T, logged string) {
				t.Helper()
				require.True(t, len(logged) == 0)
			},
		}, {
			name:    "set env",
			command: []string{"test", "0", "-eq", "0"},
			env: Environment{
				"TEST": EnvValue{
					Value: "VALUE",
				},
				"EMPTY": EnvValue{
					Value: "",
				},
				"UNSET": EnvValue{
					NeedRemove: true,
				},
			},
			expectedCode: 0,
			checkFunc: func(t *testing.T, logged string) {
				t.Helper()

				require.True(t, len(logged) == 0)

				expectedEnv := map[string]string{
					"TEST":  "VALUE",
					"EMPTY": "",
				}
				actualEnv := make(map[string]string, len(expectedEnv))
				for k := range expectedEnv {
					v, exists := os.LookupEnv(k)
					if exists {
						actualEnv[k] = v
					}
				}
				require.Equal(t, expectedEnv, actualEnv)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer func() {
				log.SetOutput(os.Stderr)
			}()
			code := RunCmd(tc.command, tc.env)
			logged := buf.String()

			require.Equal(t, tc.expectedCode, code)

			if tc.checkFunc != nil {
				tc.checkFunc(t, logged)
			}
		})
	}
}
