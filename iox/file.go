package iox

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/varunamachi/libx/errx"
)

var (
	ErrTooManyIndexedFiles   = errors.New("io.indexedFile.exceededLimit")
	ErrUnknownConflictPolicy = errors.New("io.fileCreate.unknownConflictPolicy")
)

type FileConflictPolicy int

const (
	Append FileConflictPolicy = iota
	Overwrite
	RenameOriginalWithIndex
	RenameOriginalWithTimestamp
	NameNewWithIndex
	NameNewWithTimestamp
)

type namer func(string) (string, error)

func CreateFile(
	path string, conflictPolicy FileConflictPolicy) (*os.File, error) {

	parent := filepath.Dir(path)
	if !ExistsAsDir(parent) {
		if err := os.MkdirAll(parent, fs.ModePerm); err != nil {
			return nil, errx.Errf(err, "failed to create dir at '%s'", parent)
		}
	}
	if !ExistsAsFile(path) {
		return os.Create(path)
	}

	switch conflictPolicy {
	case Append:
		return os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	case Overwrite:
		return os.Create(path)
	case RenameOriginalWithIndex:
		return replace(path, indexed)
	case RenameOriginalWithTimestamp:
		return replace(path, timestamped)
	case NameNewWithIndex:
		return createWithUpdatedName(path, indexed)
	case NameNewWithTimestamp:
		return createWithUpdatedName(path, timestamped)
	}

	return nil, errx.Errf(ErrUnknownConflictPolicy,
		"unknown name confict policy used")
}

func replace(
	path string, renamer namer) (*os.File, error) {
	newPath, err := renamer(path)
	if err != nil {
		return nil, err
	}
	if err = os.Rename(path, newPath); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func createWithUpdatedName(
	path string, renamer func(string) (string, error)) (*os.File, error) {
	newPath, err := renamer(path)
	if err != nil {
		return nil, err
	}
	return os.Create(newPath)
}

func indexed(path string) (string, error) {
	const indexLimit = 1000
	parent, fileBaseName, ext := SplitPath(path)
	if !ExistsAsDir(parent) {
		if err := os.MkdirAll(parent, fs.ModePerm); err != nil {
			return "", errx.Errf(err, "failed to create dir at '%s'", parent)
		}
		return filepath.Join(parent, fileBaseName+ext), nil
	}

	for i := 0; i < indexLimit; i++ {
		newFileName := fmt.Sprintf("%s_%d.%s", fileBaseName, i, ext)
		newPath := filepath.Join(parent, newFileName)
		if !ExistsAsFile(newPath) {
			return newPath, nil
		}
	}

	return "", errx.Errf(ErrTooManyIndexedFiles,
		"more than %d files with indexed name found at %s", indexLimit, parent)
}

func timestamped(path string) (string, error) {
	parent, fileBaseName, ext := SplitPath(path)
	timeStr := time.Now().Format("_20060102_150405")
	return filepath.Join(parent, fileBaseName+timeStr+ext), nil
}

func SplitPath(path string) (dirPath, fileBaseName, ext string) {
	fileName := filepath.Base(path)
	dirPath = filepath.Dir(path)
	ext = filepath.Ext(path)
	fileBaseName = strings.TrimSuffix(fileName, ext)

	return
}
