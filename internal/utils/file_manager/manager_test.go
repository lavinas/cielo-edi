package file_manager

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	path = "./temp"
)

var (
	listName = []string{"file1.txt", "file2.txt"}
)

func initPath() {
	err := os.Mkdir(path, 0755)
	if err != nil {
		panic(err)
	}
}

func endPath() {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err)
	}
}

func TestGetFilesOk(t *testing.T) {
	initPath()
	for _, f := range listName {
		fn := filepath.Join(path, f)
		os.WriteFile(fn, []byte("abc"), 0755)
	}
	fm := NewFileManager()
	fl, err := fm.GetFiles(path)
	assert.Nil(t, err)
	assert.Equal(t, len(listName), len(fl))
	for _, f := range fl {
		assert.Contains(t, listName, f.Name())
	}
	endPath()
}

func TestGetFilesEmptyPath(t *testing.T) {
	initPath()
	fm := NewFileManager()
	fl, err := fm.GetFiles(path)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(fl))
	assert.Equal(t, []fs.FileInfo{}, fl)
	endPath()
}

func TestGetFirstLine(t *testing.T) {
	initPath()
	first := "abc"
	fn := filepath.Join(path, listName[0])
	os.WriteFile(fn, []byte(first), 0755)
	files, err := ioutil.ReadDir(path)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(files))
	fm := NewFileManager()
	line, err := fm.GetFirstLine(path, files[0])
	assert.Nil(t, err)
	assert.Equal(t, first, line)
	endPath()
}

func TestRenameFile(t *testing.T) {
	initPath()
	first := "abc"
	fn := filepath.Join(path, listName[0])
	os.WriteFile(fn, []byte(first), 0755)
	fm := NewFileManager()
	err := fm.RenameFile(path, listName[0], listName[1])
	assert.Nil(t, err)
	files, err := ioutil.ReadDir(path)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(files))
	assert.Equal(t, listName[1], files[0].Name())
	endPath()
}
