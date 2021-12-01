package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	pathFrom   string = "/home/paulo/Desktop/nespresso/arquivos1"
	pathTo     string = "/home/paulo/Desktop/nespresso/cielo"
	dateFormat string = "20060102"
)

type Header struct {
	RegisterType    int8      //%1d
	Headquarter     int64     //%10d
	ProcDate        time.Time //%8s
	PeriodInitDate  time.Time //%8s
	PeriodEndDate   time.Time //%8s
	Sequence        int       //%7d
	Acquirer        string    //%5s
	StatementOption int8      //%2d
	Transmition     string    //%1s
	PostalBox       string    //%20s
	LayoutVersion   int8      //%3d
	Filler          string    //%177s
}

func (h *Header) Parse(txt string) error {
	var procDate string
	var periodInitDate string
	var periodEndDate string
	scan := "%1d%10d%8s%8s%8s%7d%5s%2d%1s%20s%3d%177s"
	_, err := fmt.Sscanf(string(txt), scan, &h.RegisterType, &h.Headquarter, &procDate, &periodInitDate, &periodEndDate,
		&h.Sequence, &h.Acquirer, &h.StatementOption, &h.Transmition, &h.PostalBox, &h.LayoutVersion, &h.Filler)
	if err != nil {
		return err
	}
	h.ProcDate, err = time.Parse(dateFormat, procDate)
	if err != nil {
		return err
	}
	h.PeriodInitDate, err = time.Parse(dateFormat, periodInitDate)
	if err != nil {
		return err
	}
	h.PeriodEndDate, err = time.Parse(dateFormat, periodEndDate)
	if err != nil {
		return err
	}
	return nil
}

func GetFiles(path string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(pathFrom)
	if err != nil {
		return make([]fs.FileInfo, 0), err
	}
	return files, nil
}

func ReadFile(path string, filename string) (*bufio.Scanner, error) {
	fileIO, err := os.OpenFile(filepath.Join(path, filename), os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(fileIO)
	return scanner, nil
}

func main2() {
	log.Println("Starting")
	files, err := GetFiles(pathFrom)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		buf, err := ReadFile(pathFrom, file.Name())
		if err != nil {
			log.Fatal(err)
			return
		}
		if buf.Scan() {
			var header Header
			hDate := buf.Text()
			if err := header.Parse(hDate); err != nil {
				log.Fatal(err)
			}
			fmt.Println(header)
		}
	}
	log.Println("End")
}
