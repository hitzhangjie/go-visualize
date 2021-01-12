package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

type Marshaler interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

var Marshalers = map[Serialization]Marshaler{
	JSON: &JSONMarshaler{},
	FORM: &FORMMarshaler{},
}

type JSONMarshaler struct {
}

func (m *JSONMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m *JSONMarshaler) Unmarshal(buf []byte, v interface{}) error {
	return json.Unmarshal(buf, v)
}

type FORMMarshaler struct {
}

func (m *FORMMarshaler) Marshal(v interface{}) ([]byte, error) {
	vals := url.Values{}
	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		ch := rt.Field(i).Name[0]
		if ch >= 'a' && ch <= 'z' {
			continue
		}
		tag := rt.Field(i).Tag.Get("json")
		vals.Add(tag, fmt.Sprintf("%v", rv.Field(i)))
	}
	return []byte(vals.Encode()), nil
}

func (m *FORMMarshaler) Unmarshal(buf []byte, v interface{}) error {
	vals, err := url.ParseQuery(string(buf))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v).Elem()
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		ch := rt.Field(i).Name[0]
		if ch >= 'a' && ch <= 'z' {
			continue
		}
		tag := rt.Field(i).Tag.Get("json")
		switch rv.Field(i).Kind() {
		case reflect.Int:
			v, err := strconv.ParseInt(vals.Get(tag), 10, 64)
			if err != nil {
				return err
			}
			rv.Field(i).Set(reflect.ValueOf(int(v)))
		case reflect.String:
			rv.Field(i).Set(reflect.ValueOf(vals.Get(tag)))
		default:
			return errors.New("not handled kind")
		}
	}
	return nil
}
