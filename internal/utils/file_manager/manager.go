package file_manager

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileManager struct{}

func NewFileManager() *FileManager {
	return &FileManager{}
}

func (f FileManager) GetFiles(path string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return make([]fs.FileInfo, 0), err
	}
	return files, nil
}

func (f FileManager) GetFileScanner(path string, filename string) (*bufio.Scanner, error) {
	fileIO, err := os.OpenFile(filepath.Join(path, filename), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(fileIO)
	return scanner, nil
}

func (f FileManager) GetFirstLine(path string, file fs.FileInfo) (string, error) {
	if file.IsDir() {
		return "", fmt.Errorf("%s is a directory", file.Name())
	}
	buf, err := f.GetFileScanner(path, file.Name())
	if err != nil {
		return "", err
	}
	if !buf.Scan() {
		return "", fmt.Errorf("error scanning %s", file.Name())
	}
	hDate := buf.Text()
	return hDate, nil
}

func (f FileManager) RenameFile(path string, nameFrom string, nameTo string) error {
	from := filepath.Join(path, nameFrom)
	to := filepath.Join(path, nameTo)
	if err := os.Rename(from, to); err != nil {
		return err
	}
	return nil
}
