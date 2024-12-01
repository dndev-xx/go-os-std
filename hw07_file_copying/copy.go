package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3" //nolint:depguard
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidLimit          = errors.New("invalid limit")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 || offset >= limit {
		return ErrOffsetExceedsFileSize
	}
	if limit < 0 {
		return ErrInvalidLimit
	}
	fromInfo, err := os.Stat(fromPath)
	if err != nil {
		return getErr(from, errors.Unwrap(err))
	}
	if !fromInfo.Mode().IsRegular() {
		return getErr(from, ErrUnsupportedFile)
	}
	fromSize := fromInfo.Size()
	if offset > fromSize {
		return ErrOffsetExceedsFileSize
	}
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return getErr(from, errors.Unwrap(err))
	}
	if offset > 0 {
		_, err = fromFile.Seek(offset, io.SeekStart)
		if err != nil {
			_ = fromFile.Close()
			return getErr(from, err)
		}
	}
	defer fromFile.Close()

	if fileExists(toPath) {
		err := os.Remove(toPath)
		if err != nil {
			return getErr(from, errors.Unwrap(err))
		}
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return getErr(from, errors.Unwrap(err))
	}
	defer toFile.Close()

	currentLimit := getLimit(offset, limit, fromSize)
	bar := pb.Full.Start64(currentLimit)
	barReader := bar.NewProxyReader(fromFile)
	_, err = io.CopyN(toFile, barReader, currentLimit)
	bar.Finish()
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func getErr(from string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("from file '%s': %w", from, err)
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func getLimit(offset, limit, fileSize int64) int64 {
	if limit == 0 || offset+limit > fileSize {
		return fileSize - offset
	}
	return limit
}
