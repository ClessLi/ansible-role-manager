package file_store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStore interface {
	Write(filePath string, data []byte) error
	Read(filePath string) (data []byte, err error)
	Remove(filePath string) error
	AllFiles() ([]string, error)
	Workspace() (string, error)
}

var fileLock = new(sync.Mutex)

type fileStore struct {
	workspace string
}

func (f *fileStore) Write(filePath string, data []byte) error {
	absFilePath, err := f.checkFilePath(filePath)
	if err != nil {
		return err
	}

	err = checkDir(f.workspace)
	if err != nil {
		return err
	}

	return os.WriteFile(absFilePath, data, 0644)
}

func (f *fileStore) Read(filePath string) (data []byte, err error) {
	absFilePath, err := f.checkFilePath(filePath)
	if err != nil {
		return nil, err
	}

	err = checkDir(f.workspace)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(absFilePath)
}

func (f *fileStore) Remove(filePath string) error {
	absFilePath, err := f.checkFilePath(filePath)
	if err != nil {
		return err
	}

	err = checkDir(f.workspace)
	if err != nil {
		return err
	}

	return os.Remove(absFilePath)
}

func (f *fileStore) AllFiles() ([]string, error) {

	err := checkDir(f.workspace)
	if err != nil {
		return nil, err
	}

	// change directory to workspace
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = os.Chdir(f.workspace)
	if err != nil {
		return nil, err
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Printf("failed to change dir back to %s", dir)
		}
	}(pwd)

	filePaths, err := filepath.Glob("*")
	if err != nil {
		return nil, err
	}

	for j := 0; j < len(filePaths); j++ {
		fStat, err := os.Stat(filePaths[j])
		if err == nil && !fStat.IsDir() {
			continue

		}
		if err != nil {
			fmt.Printf("check %s error, %s\n", filePaths[j], err)
		}
		filePaths = append(filePaths[:j], filePaths[j+1:]...) // remove this file path record
	}

	if len(filePaths) == 0 {
		err = errors.New("no file in this file-store")
	}

	return filePaths, err
}

func (f *fileStore) Workspace() (string, error) {
	return f.workspace, checkDir(f.workspace)
}

func (f *fileStore) checkFilePath(filePath string) (absFilePath string, err error) {
	if filePath[0] != '/' {
		filePath = filepath.Join(f.workspace, filePath)
	}
	relativeFilePath, err := filepath.Rel(f.workspace, filePath)
	if err != nil {
		return "", err
	}

	if relativeFilePath[0] == '.' && relativeFilePath[1] == '.' {
		return "", fmt.Errorf("invalid file path: %v", filePath)
	}

	return filepath.Abs(filepath.Join(f.workspace, relativeFilePath))
}

func NewFileStore(dirPath string) (FileStore, error) {
	err := checkDir(dirPath)
	if err != nil {
		return nil, err
	}

	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}

	fs := &fileStore{workspace: absDirPath}
	return FileStore(fs), nil
}

func checkDir(dirPath string) error {
	fileLock.Lock()
	defer fileLock.Unlock()
	dirStat, err := os.Stat(dirPath)
	if err != nil {
		return err
	}
	if !dirStat.IsDir() {
		return fmt.Errorf("%s is not a directory", dirPath)
	}

	return nil
}
