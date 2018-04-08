package main

import (
	"fmt"
	"reflect"
)

type Header struct {
	TotalSize int16
	Id        int16
}

type Login struct {
	header Header
	Id     [16]byte
	Pw     [16]byte
}

var HeaderSizeROnly = protocolInitHeaderSize()

func protocolInitHeaderSize() int16 {
	var packetHeader Header
	headerSize := sizeof(reflect.TypeOf(packetHeader))
	return (int16)(headerSize)
}

func sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		fmt.Println("reflect.Array")
		if s := sizeof(t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		fmt.Println("reflect.Struct")
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		//fmt.Println("reflect.int")
		return int(t.Size())
	case reflect.Slice:
		//fmt.Println("reflect.Slice:", sizeof(t.Elem()))
		return 0
	}

	return -1

}
