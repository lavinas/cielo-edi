package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRename(t *testing.T) {
	path := "/home/paulo/Desktop/nespresso/arquivos-x"
	args := make([]string, 0)
	err := rename(path, args)
	assert.Nil(t, err)	
}

func TestGaps(t *testing.T) {
	path := "/home/paulo/Desktop/nespresso/arquivos-x"
	args := []string{"1", "2", "01/01/2021", "10/01/2021"}
	err := gaps(path, args)
	assert.Nil(t, err)
	args = []string{"1", "2", "22/01/2021", "10/01/2021"}
	err = gaps(path, args)
	assert.NotNil(t, err)
}

func TestExec(t *testing.T) {
	args := []string{"program", "rename", "/home/paulo/Desktop/nespresso/arquivos-x"}
	err := exec(args)
	assert.Nil(t, err)
	args = []string{"program", "rename", ""}
	err = exec(args)
	assert.NotNil(t, err)
	args = []string{"program", "gaps", "/home/paulo/Desktop/nespresso/arquivos-x"}
	err = exec(args)
	assert.NotNil(t, err)
	args = []string{"program", "gaps", "/home/paulo/Desktop/nespresso/arquivos-x", "01/01/2021", "10/01/2021"}
	err = exec(args)
	assert.Nil(t, err)
	args = []string{"program", "gaps", "/home/paulo/Desktop/nespresso/arquivos-x", "10/01/2021", "01/01/2021"}
	err = exec(args)
	assert.NotNil(t, err)

}
