package main

import (
	"bytes"
	"errors"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrIsNotDir  = errors.New("the current path is not directory")
	ErrNoSuchDir = errors.New("no such directory")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return nil, ErrNoSuchDir
	}

	if !fileInfo.IsDir() {
		return nil, ErrIsNotDir
	}

	return getEnvironmentsInDir(dir)
}

func getEnvironmentsInDir(dir string) (Environment, error) {
	objects, readDirErr := os.ReadDir(dir)
	environment := make(Environment)

	for _, object := range objects {
		if object.IsDir() {
			continue
		}

		filename := object.Name()
		if strings.Contains(filename, "=") {
			continue
		}

		value, err := getEnvironmentValue(path.Join(dir, filename))
		if err != nil {
			continue
		}

		environment[filename] = value
	}

	return environment, readDirErr
}

func getEnvironmentValue(filename string) (EnvValue, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return EnvValue{}, err
	}
	if len(data) == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	return EnvValue{Value: getStringValue(data)}, nil
}

func getStringValue(data []byte) (strData string) {
	lines := bytes.Split(data, []byte{'\n'})
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		strData = string(bytes.ReplaceAll(line, []byte{0x00}, []byte{'\n'}))
		break
	}
	strData = strings.TrimRight(strData, "\t ")
	return
}
