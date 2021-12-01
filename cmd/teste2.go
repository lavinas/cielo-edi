package main

import (
	"fmt"
	"log"
	"reflect"
	"time"
	"strings"
	"strconv"
)

type Foo struct {
	FirstName string `tag_name:"tag 1"`
	LastName  string `tag_name:"tag 2"`
	Age       int    `tag_name:"tag 3"`
}

type HeaderX struct {
	RegisterType    int8      `format:"d,1"` //%1d
	Headquarter     int64     `format:"d,10"` //%10d
	ProcDate        time.Time `format:"t,yyyymmdd"` //%8s
	PeriodInitDate  time.Time `format:"t,yyyymmdd"` //%8s
	PeriodEndDate   time.Time `format:"t,yyyymmdd"`//%8s
	Sequence        int       `format:"d,7"` //%7d
	Acquirer        string    `format:"s,5"` //%5s
	StatementOption int8      `format:"d,2"` //%2d
	Transmition     string    `format:"s,1"` //%1s
	PostalBox       string    `format:"s,20"` //%20s
	LayoutVersion   int8      `format:"d, 3"` //%3d
	Filler          string    `format:"s,177"` //%177s
}


func (f *Foo) Reflect() {
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		x := typeField.Type
		fmt.Println(x)
		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("tag_name"))
	}
}

func (f *HeaderX) Parse(txt string) error {
	var strPosition int = 0
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		tag_val := tag.Get("format")
		if tag_val == "" {
			return fmt.Errorf("no format annotation in field %s", typeField.Name)
		}
		val := strings.Split(tag_val, ",")
		if len(val) != 2 {
			return fmt.Errorf("invalid annotation '%s' in field %s", tag_val, typeField.Name)
		}
		if val[0] == "d" {
			var dimN int = 32
			var err error
			if len(typeField.Type.Name()) < 3 {
				return fmt.Errorf("invalid annotation '%s' in field %s for type %s", tag_val, typeField.Name, typeField.Type)
			}
			if typeField.Type.Name()[:3] != "int" {
				return fmt.Errorf("invalid annotation '%s' in field %s for type %s", tag_val, typeField.Name, typeField.Type)
			}
			dim := typeField.Type.Name()[3:]
			if dim != "" {
				dimN, err = strconv.Atoi(dim)
				if err != nil {
					return fmt.Errorf("invalid annotation '%s' in field %s for type %s", tag_val, typeField.Name, typeField.Type)
				}
			}
			lag, err := strconv.Atoi(val[1])
			if err != nil {
				return fmt.Errorf("invalid annotation '%s' in field %s", tag_val, typeField.Name)
			}
			if strPosition + lag >= len(txt) {
				return fmt.Errorf("invalid string length")
			}
			sub := txt[strPosition: strPosition + lag]
			dec, err := strconv.ParseInt(sub, 10, dimN)
			if err != nil {
				return fmt.Errorf("invalid dimention")
			}
			reflect.ValueOf(f).Elem().FieldByName(typeField.Name).SetInt(dec)
		} 
		
		
	}
	return nil
}


func (f *HeaderX) Reflect() {
	val := reflect.ValueOf(f).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		x := typeField.Type
		fmt.Println(x)
		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Tag Value: %s\n", typeField.Name, valueField.Interface(), tag.Get("format"))
	}
}

func main() {
	/*
	f := &Foo{
		FirstName: "Drew",
		LastName:  "Olson",
		Age:       30,
	}
	*/

	f := &HeaderX{
		RegisterType: 1,
		Headquarter: 123213,
	}

	// f.Reflect()
	err := f.Parse("abcssss")
	if err != nil {
		log.Fatal(err)
	}
}