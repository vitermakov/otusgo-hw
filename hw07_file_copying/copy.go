package main

import (
	"errors"
	"io"
	"os"

	"github.com/vitermakov/otusgo-hw/hw07_file_copying/progress"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	fileIn, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = fileIn.Close()
	}()

	stats, err := fileIn.Stat()
	if err != nil {
		return err
	}
	if stats.IsDir() || !stats.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	size := stats.Size()
	if offset >= size {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 || limit+offset > size {
		limit = size - offset
	}

	fileOut, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = fileOut.Close()
	}()

	_, err = fileIn.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	bar := progress.New(fileIn, limit)
	_, err = io.CopyN(fileOut, bar, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
