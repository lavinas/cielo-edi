package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	tag_name      = "parse_txt"
	tag_delimiter = ','
)

func splitTag(tag string) (string, string, error) {
	if tag == "" {
		return "", "", fmt.Errorf("tag is empty")
	}
	val := strings.Split(tag, ",")
	if len(val) != 2 {
		return "", "", fmt.Errorf("invalid annotation")
	}
	return val[0], val[1], nil
}

func getField(fType string, format string, line string, linePos int) (string, int, error) {
	var fLen int
	var err error
	if fType == "t" {
		fLen = len(format)
	} else {
		fLen, err = strconv.Atoi(format)
		if err != nil {
			return "", linePos, fmt.Errorf("invalid tag format")
		}
	}
	if linePos + fLen > len(line) {
		return "", linePos, fmt.Errorf("line is shorter than field length")
	}
	return line[linePos: linePos+fLen], fLen, nil
}

func getDecimal(fType string, value string) (int64, error) {
	var dimN int = 32
	var err error
	if len(fType) < 3 || fType[:3] != "int" {
		return 0, fmt.Errorf("invalid field type (should be int, int8, int16, int32, int64)")
	}
	dim := fType[3:]
	if dim != "" {
		dimN, err = strconv.Atoi(dim)
		if err != nil {
			return 0, fmt.Errorf("invalid field type (should be int, int8, int16, int32, int64)")
		}
	}
	if !(dimN == 8 || dimN == 16 || dimN == 32 || dimN == 64) {
		return 0, fmt.Errorf("invalid field type (should be int, int8, int16, int32, int64)")
	}
	dec, err := strconv.ParseInt(value, 10, dimN)
	if err != nil {
		return 0, fmt.Errorf("parsing integer error")
	}
	return dec, nil
}

func getString(fType, value string) (string, error) {
	if fType != "string" {
		return "", fmt.Errorf("invalid field type (should be string)")
	}
	return value, nil
}

func getTime(fType, format, value string)(time.Time, error) {
	if fType != "Time" {
		return time.Time{}, fmt.Errorf("invalid field type (should be Time)")
	}
	t, err :=time.Parse(format, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format (should by yyyymmdd, yymmdd, etc")
	}
	return t, nil
}

func UnmarshalField(source interface{}, fName string, line string, linePos int) (int, error) {
	// find fields
	if !reflect.ValueOf(source).Elem().FieldByName(fName).CanSet() {
		return linePos, fmt.Errorf("invalid field name")
	}
	field, _ := reflect.ValueOf(source).Elem().Type().FieldByName(fName)
	fType := field.Type.Name()
	tag := field.Tag.Get(tag_name)
	// format tag
	tagType, tagFormat, err := splitTag(tag)
	if err != nil {
		return linePos, fmt.Errorf("invalid tag structure (should have <type, format>): field %s", fName)
	}
	if !(tagType == "d" || tagType == "s" || tagType == "t") {
		return linePos, fmt.Errorf("invalid tag type")
	} 
	// set values
	var fValue string
	var fLen int 
	switch tagType {
	case "d":
		fValue, fLen, err = getField(tagType, tagFormat, line, linePos)
		if err != nil {
			return linePos, err
		}
		dval, err := getDecimal(fType, fValue)
		if err != nil {
			return linePos, err
		}
		reflect.ValueOf(source).Elem().FieldByName(fName).SetInt(dval)
	case "s":
		fValue, fLen, err = getField(tagType, tagFormat, line, linePos)
		if err != nil {
			return linePos, err
		}
		sVal, err := getString(fType, fValue)
		if err != nil {
			return linePos, err
		}
		reflect.ValueOf(source).Elem().FieldByName(fName).SetString(sVal)
	case "t":
		var tVal time.Time
		fValue, fLen, err = getField(tagType, tagFormat, line, linePos)
		if err != nil {
			return linePos, err
		}
		tVal, err = getTime(fType, tagFormat, fValue)
		if err != nil {
			return linePos, err
		}
		reflect.ValueOf(source).Elem().FieldByName(fName).Set(tVal)

	default:
		return linePos, fmt.Errorf("invalid tag type")
	}
	return linePos + fLen, nil
}

func UnMarshal(source interface{}, line string) error {
	var strPosition int = 0
	var err error

	fields := reflect.ValueOf(source).Elem()
	for i := 0; i < fields.NumField(); i++ {
		fieldName := fields.Type().Field(i).Type.Name()
		strPosition, err = UnmarshalField(source, fieldName, line, strPosition)
		if err != nil {
			err = errors.Wrap(err, fmt.Sprintf("field: %s", fieldName))
			return err
		}
	}
	return nil
}
