package file_store

import (
	"fmt"
	"github.com/marmotedu/errors"
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
		return errors.Wrapf(err, "failed to write file '%s'", filePath)
	}

	err = checkDir(f.workspace)
	if err != nil {
		return errors.Wrapf(err, "failed to write file '%s'", filePath)
	}

	return errors.Wrapf(os.WriteFile(absFilePath, data, 0644), "write file '%s' error", filePath)
}

func (f *fileStore) Read(filePath string) (data []byte, err error) {
	absFilePath, err := f.checkFilePath(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file '%s'", filePath)
	}

	err = checkDir(f.workspace)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file '%s'", filePath)
	}

	data, err = os.ReadFile(absFilePath)
	return data, errors.Wrapf(err, "read file '%s' error", filePath)
}

func (f *fileStore) Remove(filePath string) error {
	absFilePath, err := f.checkFilePath(filePath)
	if err != nil {
		return errors.Wrapf(err, "failed to remove file '%s'", filePath)
	}

	err = checkDir(f.workspace)
	if err != nil {
		return errors.Wrapf(err, "failed to remove file '%s'", filePath)
	}

	return errors.Wrapf(os.Remove(absFilePath), "remove file '%s' error", filePath)
}

func (f *fileStore) AllFiles() ([]string, error) {

	err := checkDir(f.workspace)
	if err != nil {
		return nil, errors.Wrap(err, "get all files error")
	}

	// change directory to workspace
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	err = os.Chdir(f.workspace)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Printf("failed to change dir back to %s", dir)
		}
	}(pwd)

	filePaths, err := filepath.Glob("*")
	if err != nil {
		return nil, errors.New(err.Error())
	}

	errs := make([]error, 0)
	for j := 0; j < len(filePaths); j++ {
		fStat, err := os.Stat(filePaths[j])
		if err == nil && !fStat.IsDir() {
			continue

		}
		if err != nil {
			errs = append(errs, errors.New(err.Error()))
			fmt.Printf("check %s error, %s\n", filePaths[j], err)
		}
		filePaths = append(filePaths[:j], filePaths[j+1:]...) // remove this file path record
	}

	if len(filePaths) == 0 {
		//err = errors.New("no file in this file-store")
		err = errors.New("no file in this FileStore")
		errs = append(errs, err)
	}

	return filePaths, errors.NewAggregate(errs)
}

func (f *fileStore) Workspace() (string, error) {
	return f.workspace, errors.Wrap(checkDir(f.workspace), "workspace error")
}

func (f *fileStore) checkFilePath(filePath string) (absFilePath string, err error) {
	if filePath[0] != '/' {
		filePath = filepath.Join(f.workspace, filePath)
	}
	relativeFilePath, err := filepath.Rel(f.workspace, filePath)
	if err != nil {
		return "", errors.New(err.Error())
	}

	if relativeFilePath[0] == '.' && relativeFilePath[1] == '.' {
		return "", errors.Errorf("invalid file path: %v", filePath)
	}

	return filepath.Abs(filepath.Join(f.workspace, relativeFilePath))
}

func NewFileStore(dirPath string) (FileStore, error) {
	err := checkDir(dirPath)
	if err != nil {
		return nil, errors.Wrap(err, "build FileStore error")
	}

	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, errors.Wrap(err, "build FileStore error")
	}

	fs := &fileStore{workspace: absDirPath}
	return FileStore(fs), nil
}

func checkDir(dirPath string) error {
	fileLock.Lock()
	defer fileLock.Unlock()
	dirStat, err := os.Stat(dirPath)
	if err != nil {
		return errors.New(err.Error())
	}
	if !dirStat.IsDir() {
		return errors.Errorf("%s is not a directory", dirPath)
	}

	return nil
}
