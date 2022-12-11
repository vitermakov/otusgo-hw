package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		var value EnvValue
		if err != nil {
			return err
		}

		name := info.Name()
		if !info.Mode().IsRegular() ||
			strings.ContainsAny(name, ";=") {
			return nil
		}
		// по заданию имеется в виду, что если файл полностью пустой, переменная помечается
		// на удаление, т.е. мы не учитываем, что значение может стать пустым после TrimRight
		if info.Size() == 0 {
			value.NeedRemove = true
		} else {
			v, err := readFirstLine(path)
			if err != nil {
				return err
			}
			value.Value = clearValue(v)
		}

		env[name] = value

		return nil
	})
	if err != nil {
		return nil, err
	}

	return env, nil
}

func clearValue(value string) string {
	value = strings.TrimRight(value, " \t\r")
	return strings.ReplaceAll(value, "\x00", "\n")
}

func readFirstLine(filePath string) (string, error) {
	handle, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = handle.Close()
	}()

	reader := bufio.NewReader(handle)
	bs, _, err := reader.ReadLine()
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}

	return string(bs), nil
}
