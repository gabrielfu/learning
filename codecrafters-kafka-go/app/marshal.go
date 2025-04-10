package main

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"reflect"
)

func marshalInt8BigEndian(value int8) []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, value)
	return b.Bytes()
}

func marshalInt32BigEndian(value int32) []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, value)
	return b.Bytes()
}

func unmarshalInt32BigEndian(data []byte) (int32, error) {
	var value int32
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func marshalInt16BigEndian(value int16) []byte {
	b := new(bytes.Buffer)
	binary.Write(b, binary.BigEndian, value)
	return b.Bytes()
}

func unmarshalInt16BigEndian(data []byte) (int16, error) {
	var value int16
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.BigEndian, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func marshalBool(value bool) byte {
	if value {
		return 1
	}
	return 0
}

func marshal(value interface{}) ([]byte, error) {
	switch value := value.(type) {
	case int8:
		return marshalInt8BigEndian(value), nil
	case int16:
		return marshalInt16BigEndian(value), nil
	case int32:
		return marshalInt32BigEndian(value), nil
	case bool:
		return []byte{marshalBool(value)}, nil
	case byte:
		return []byte{value}, nil
	case encoding.BinaryMarshaler:
		return value.MarshalBinary()
	}

	r := reflect.ValueOf(value)
	switch r.Kind() {
	case reflect.Struct:
		return marshalStruct(value)
	default:
		return nil, fmt.Errorf("unsupported type: %s", r.Kind())
	}
}

func marshalStruct(value interface{}) ([]byte, error) {
	var out []byte
	r := reflect.ValueOf(value)
	numFields := r.NumField()
	for i := 0; i < numFields; i++ {
		field := r.Field(i)
		b, err := marshal(field.Interface())
		if err != nil {
			return nil, err
		}
		out = append(out, b...)
	}
	return out, nil
}
