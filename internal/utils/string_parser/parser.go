package string_parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	// tag_name identifies a struct tag that contains parameters for parsing text
	tag_name = "txt"
)

var (
	// time_replacer maps formats thar can be identified on tags for formating Time fields
	time_replacer = map[string]string{
		"yyyymmdd":   "20060102",
		"yyyy-mm-dd": "2006-01-02",
		"yymmdd":     "060102",
		"yy-mm-dd":   "06-01-02",
		"ddmmyyyy":   "02012006",
		"dd-mm-yyyy": "02-01-2006",
	}
	// time_type maps struct fields that this program can support for parsing txt
	// others types should be implemented here
	field_type = map[string]string{
		"int":    "d",
		"int8":   "d",
		"int16":  "d",
		"int32":  "d",
		"int64":  "d",
		"string": "s",
		"Time":   "t",
	}
)

// StringParser has ability to Parse text strings into the fields of generic struct
type StringParser struct{
	parserType string
}

// UnmarshalField try to find a structure field and parse the value of a string based on the parameters of this field
//
// source has a structure that possible have the field
// fieldName has the name of field to be found
// txt has the string to be parsed based on the parameters of this field
// txtPos has the position of the string to start to parse
//
// returns the next string field position based on the start position and the field length and a possible error
func (s StringParser) ParseField(source interface{}, fieldName string, txt string, txtPos int) (int, error) {
	// verify source
	if err := verifyValidInterface(source); err != nil {
		return txtPos, err
	}
	// find fields
	fieldType, fieldIndex, fieldTag, err := getFieldByName(source, fieldName)
	if err != nil {
		return txtPos, err
	}
	// set values
	var fieldLen int
	switch fieldIndex {
	case "d":
		var dval int64
		dval, fieldLen, err = getDecimal(fieldType, fieldIndex, fieldTag, s.parserType, txt, txtPos)
		if err != nil {
			return txtPos, err
		}
		reflect.ValueOf(source).Elem().FieldByName(fieldName).SetInt(dval)
	case "s":
		var sVal string
		sVal, fieldLen, err = getString(fieldIndex, fieldTag, s.parserType, txt, txtPos)
		if err != nil {
			return txtPos, err
		}
		reflect.ValueOf(source).Elem().FieldByName(fieldName).SetString(sVal)
	case "t":
		var tVal time.Time
		tVal, fieldLen, err = getTime(fieldIndex, fieldTag, s.parserType, txt, txtPos)
		if err != nil {
			return txtPos, err
		}
		reflect.ValueOf(source).Elem().FieldByName(fieldName).Set(reflect.ValueOf(tVal))

	default:
		return txtPos, fmt.Errorf("invalid tag type")
	}
	return txtPos + fieldLen, nil
}

// Unmarshal try to find all structure fields values on a sequenced string based on this parameters (types and tags)
// the possibles tags values are: a numeric value that represents the substring length if the field is integer or string or
// a date format (ex yyyymmdd) if the field is a Time
//
// source has a structure that possible have the field
// txt has the string to be parsed based on the parameters of this field
//
// returns a error if there is one, otherwise fills source structure with the txt sequenced values
func (s StringParser) Parse(source interface{}, txt string) error {
	// verify source
	if err := verifyValidInterface(source); err != nil {
		return err
	}
	// unmarshall all fields
	var strPosition int = 0
	var err error
	fields := reflect.ValueOf(source).Elem()
	for i := 0; i < fields.NumField(); i++ {
		fieldName := fields.Type().Field(i).Name
		strPosition, err = s.ParseField(source, fieldName, txt, strPosition)
		if err != nil {
			err = errors.Wrap(err, fieldName)
			return err
		}
	}
	return nil
}

func NewStringParser(parserType string) *StringParser {
	return &StringParser{parserType: strings.ToLower(parserType)}
}

// getValue parse string and returns substring and its length based on field parameters (Index(type) and tag)
//
// fieldIndex describes field type (D - decimal/integer, S - string, T - Time)
// tagValue has the tag value of struct field
// txt has a string that is parsed
// txtPos has the position of txt to start to parse
//
// returns a substring parsed, the len of this string and a possible error
func getValue(fieldIndex string, tagValue string, ptype string, txt string, txtPos int) (string, int, error) {
	var value string
	var addPos int
	switch ptype {
	case "position":
		var fieldLen int
		var err error
		if fieldIndex == "t" {
			fieldLen = len(tagValue)
		} else {
			fieldLen, err = strconv.Atoi(tagValue)
			if err != nil {
				return "", txtPos, fmt.Errorf("invalid tag value (should be numeric)")
			}
		}
		if txtPos+fieldLen > len(txt) {
			return "", txtPos, fmt.Errorf("unexpected end of txt for parsing this field")
		}
		value = txt[txtPos : txtPos+fieldLen]
		addPos = fieldLen
	case "csv":
		txtSplit := strings.Split(txt, ",")
		if txtPos >= len(txtSplit) {
			return "", txtPos, fmt.Errorf("unexpected end of csv for parsing this field")
		}
		value = txtSplit[txtPos]
		addPos = 1
	default:
		return "", txtPos, fmt.Errorf("unexpected type of parser")
	}
	return value, addPos, nil
}

// getDecimal parse string and returns a integer value and its length based on field parameters (Index(type) and tag)
//
// fieldIndex describes field type (D - decimal/integer, S - string, T - Time)
// tagValue has the tag value of struct field
// txt has a string that is parsed
// txtPos has the position of txt to start to parse
//
// returns a substring parsed transformed in a integer (8, 16, 32 or 64), the len of this string and a possible error
func getDecimal(fieldType string, fieldIndex string, tagValue string, ptype string, txt string, txtPos int) (int64, int, error) {
	var dimN int = 32
	var err error
	value, txtPos, err := getValue(fieldIndex, tagValue, ptype, txt, txtPos)
	if err != nil {
		return 0, 0, err
	}
	dim := fieldType[3:]
	if dim != "" {
		dimN, _ = strconv.Atoi(dim)
	}
	dec, err := strconv.ParseInt(value, 10, dimN)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing integer error")
	}
	return dec, txtPos, nil
}

// getStrings parse string and returns a string value and its length based on field parameters (Index(type) and tag)
//
// fieldIndex describes field type (D - decimal/integer, S - string, T - Time)
// tagValue has the tag value of struct field
// txt has a string that is parsed
// txtPos has the position of txt to start to parse
//
// returns a substring parsed, the len of this string and a possible error
func getString(fieldIndex string, tagValue string, ptype string, txt string, txtPos int) (string, int, error) {
	value, txtPos, err := getValue(fieldIndex, tagValue, ptype, txt, txtPos)
	if err != nil {
		return "", 0, err
	}
	return value, txtPos, nil
}

// getStrings parse string and returns a Time value and its length based on field parameters (Index(type) and tag)
//
// fieldIndex describes field type (D - decimal/integer, S - string, T - Time)
// tagValue has the tag value of struct field
// txt has a string that is parsed
// txtPos has the position of txt to start to parse
//
// returns a substring parsed transformed in a Time variable, the len of this string and a possible error
func getTime(fieldIndex string, tagValue string, ptype string, txt string, txtPos int) (time.Time, int, error) {
	value, txtPos, err := getValue(fieldIndex, tagValue, ptype, txt, txtPos)
	if err != nil {
		return time.Time{}, 0, err
	}
	rFormat := time_replacer[tagValue]
	if rFormat == "" {
		return time.Time{}, 0, fmt.Errorf("invalid datetime tag value (should be for ex yyyymmdd)")
	}
	t, err := time.Parse(rFormat, value)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf(fmt.Sprintf("%v", err))
	}
	return t, txtPos, nil
}

// getFieldByName try to find a structure field parameters (type, index)
//
// source has a structure that possible have the field
// fieldName has the name of field to be found
//
// returns the name of field, a field type/index (D - decimal/integer, S - string, t _time), the tag of the field and a possible error
func getFieldByName(source interface{}, fieldName string) (string, string, string, error) {
	// find fields
	if !reflect.ValueOf(source).Elem().FieldByName(fieldName).CanSet() {
		return "", "", "", fmt.Errorf("invalid field name")
	}
	field, _ := reflect.ValueOf(source).Elem().Type().FieldByName(fieldName)
	typeName := field.Type.Name()
	typeIndex := field_type[typeName]
	if typeIndex == "" {
		return "", "", "", fmt.Errorf("not supported field type")
	}
	tag := field.Tag.Get(tag_name)
	if tag == "" {
		return "", "", "", fmt.Errorf("tag is not presented")
	}

	return typeName, typeIndex, tag, nil
}

// verifyValidInterface ckecks if a interface is a structure
//
// source has a interface that should be a structure
//
// returns a error if is is not a structure
func verifyValidInterface(source interface{}) error {
	if source == nil {
		return fmt.Errorf("source interface should be a valid struct")
	}
	x := reflect.TypeOf(source).Elem().Kind()
	if x != reflect.Struct {
		return fmt.Errorf("source interface should be a valid struct")
	}
	return nil
}
