package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cm = NewCommandLine()
)


func TestExec(t *testing.T) {
	args := []string{"program", "rename", "cielovendas", "/home/paulo/Desktop/nespresso/arquivos-x"}
	err := cm.Run(args)
	assert.Nil(t, err)
	args = []string{"program", "rename", "cielovendas", ""}
	err = cm.Run(args)
	assert.NotNil(t, err)
	args = []string{"program", "gaps", "cielovendas", "/home/paulo/Desktop/nespresso/arquivos-x"}
	err = cm.Run(args)
	assert.NotNil(t, err)
	args = []string{"program", "gaps", "cielovendas", "/home/paulo/Desktop/nespresso/arquivos-x", "01/01/2021", "10/01/2021"}
	err = cm.Run(args)
	assert.Nil(t, err)
	args = []string{"program", "gaps", "cielovendas", "/home/paulo/Desktop/nespresso/arquivos-x", "10/01/2021", "01/01/2021"}
	err = cm.Run(args)
	assert.NotNil(t, err)
	assert.Equal(t, "initDate after endDate", err.Error())
	args = []string{"program", "periods", "cielovendas", "/home/paulo/Desktop/nespresso/arquivos-x", "01/01/2021", "10/01/2021"}
	err = cm.Run(args)
	assert.Nil(t, err)
}

func TestExecErrorParam(t *testing.T) {
	args := []string{"program", "xxxx", "cielovendas", "/home/paulo/Desktop/nespresso/arquivos-x"}
	err := cm.Run(args)
	assert.NotNil(t, err)
	assert.Equal(t, "command xxxx not found (should be rename, gaps or periods)", err.Error())
}

func TestExecErrorParam2(t *testing.T) {
	args := []string{"program", "/home/paulo/Desktop/nespresso/arquivos-x"}
	err := cm.Run(args)
	assert.NotNil(t, err)
	assert.Equal(t, "wrong number of parameters (should be ./command-line command acquirer path", err.Error())
}

func TestExecErrorParam3(t *testing.T) {
	args := []string{"program", "rename", "/home/paulo/Desktop/nespresso/arquivos-x", "cielovendas"}
	err := cm.Run(args)
	assert.NotNil(t, err)
	assert.Equal(t, "acquirer name /home/paulo/Desktop/nespresso/arquivos-x not found (should be cielo, redecredito, getnet)", err.Error())
}

func TestExecRedCredito(t *testing.T) {
	args := []string{"program", "rename", "redecredito", "/home/paulo/Desktop/nespresso/arquivos-x"}
	err := cm.Run(args)
	assert.Nil(t, err)
}

func TestExecGetnet(t *testing.T) {
	args := []string{"program", "rename", "getnet", "/home/paulo/Desktop/nespresso/arquivos3"}
	err := cm.Run(args)
	assert.Nil(t, err)
}

func TestExecRedDebito(t *testing.T) {
	args := []string{"program", "rename", "rededebito", "/home/paulo/Desktop/nespresso/arquivos-x"}
	err := cm.Run(args)
	assert.Nil(t, err)
}

