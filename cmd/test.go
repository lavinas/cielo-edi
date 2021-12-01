package main

import (
	"fmt"
	"log"
	"time"
)

const (
	dateFormat2 string = "20060102"
)

type Header2 struct {
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

func main5() {
	var str string = "010238632322021063020210630202106300008358CIELO04I                    014                                                                                                                                                                                 "
	var scan string = "%1d%10d%8s%8s%8s%7d%5s%2d%1s"
	var procDate string
	var periodInitDate string
	var periodEndDate string
	h := Header2{}
	str2 := str[:51]
	var x string

	fmt.Println("*" + str2 + "*")

	_, err := fmt.Sscanf(str2, scan, &h.RegisterType, &h.Headquarter, &procDate, &periodInitDate, &periodEndDate, &h.Sequence, &h.Acquirer, &h.StatementOption, &h.Transmition, &x)
	if err != nil {
		log.Fatal(err)
	}

	h.ProcDate, err = time.Parse(dateFormat2, procDate)
	if err != nil {
		log.Fatal(err)
	}
	
	h.PeriodInitDate, err = time.Parse(dateFormat2, periodInitDate)
	if err != nil {
		log.Fatal(err)
	}

	h.PeriodEndDate, err = time.Parse(dateFormat2, periodEndDate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", h)
}
