package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
)

type TestStruct struct {
	Name  string
	Value int
}

type TypeByte struct {
	t reflect.Type
	b *bytes.Buffer
}

func StructToMap(intf interface{}) map[string]TypeByte {
	col := map[string]TypeByte{}
	s := reflect.ValueOf(intf).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		buffer := new(bytes.Buffer)
		f := s.Field(i)
		k := f.Type().Kind()
		switch k {
		case reflect.Int:
			binary.Write(buffer, binary.LittleEndian, f.Int())
		case reflect.String:
			binary.Write(buffer, binary.LittleEndian, int64(len(f.String())))
			buffer.Write([]byte(f.String()))
		default:
			log.Fatalf("Type %s is not supported yet\n", k)
		}
		tb := col[typeOfT.Field(i).Name]
		tb.b = buffer
		tb.t = f.Type()
		col[typeOfT.Field(i).Name] = tb
	}
	return col
}

func MapToStruct(m map[string]TypeByte, intf interface{}) {
	s := reflect.ValueOf(intf).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		tb := m[typeOfT.Field(i).Name]
		k := f.Type().Kind()
		switch k {
		case reflect.Int:
			var dest int64 = 0
			binary.Read(tb.b, binary.LittleEndian, &dest)
			f.SetInt(dest)
		case reflect.String:
			var l int64 = 0
			binary.Read(tb.b, binary.LittleEndian, &l)
			var text []byte = make([]byte, l)
			tb.b.Read(text)
			f.SetString(string(text))
		default:
			log.Fatalf("Type %s is not supported yet\n", k)
		}
	}
}

func main() {
	a := TestStruct{
		Name:  "GoLang",
		Value: 18,
	}

	m := StructToMap(&a)
	aRet := TestStruct{}
	MapToStruct(m, &aRet)
	fmt.Printf("output: %s, %d\n", aRet.Name, aRet.Value)
}
