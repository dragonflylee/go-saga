package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"
)

var (
	zero      = []byte{0}
	byteSlice = reflect.TypeOf(zero)
)

// marshalValue 序列化
func marshalValue(w io.Writer, v reflect.Value) (int, error) {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		if v.Type() == byteSlice {
			return w.Write(v.Bytes())
		}
		l := v.Len()
		off := 1
		w.Write([]byte{byte(l)})
		for i := 0; i < l; i++ {
			n, err := marshalValue(w, v.Index(i))
			if err != nil {
				return 0, err
			}
			off += n
		}
		return off, nil

	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		off := 0
		for i := 0; i < l; i++ {
			if t.Field(i).Tag.Get("saga") == "-" {
				continue
			}
			if f := v.Field(i); f.IsValid() {
				n, err := marshalValue(w, f)
				if err != nil {
					return 0, err
				}
				off += n
			}
		}
		return off, nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		err := binary.Write(w, binary.BigEndian, v.Interface())
		return int(v.Type().Size()), err

	case reflect.String:
		l := v.Len()
		if l <= 0 {
			return w.Write(zero)
		}
		w.Write([]byte{byte(l + 1)})
		n, err := w.Write(append([]byte(v.String()), byte(0)))
		return n + 1, err
	}
	return 0, fmt.Errorf("unsupport type `%s`", v.Type())
}

// unmarshalValue 反序列化
func unmarshalValue(data []byte, v reflect.Value) (int, error) {
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		l := int(data[0])
		off := 1
		if v.Kind() == reflect.Slice {
			v.SetLen(l)
		} else {
			l = v.Type().Len()
		}
		for i := 0; i < l; i++ {
			n, err := unmarshalValue(data[off:], v.Index(i))
			if err != nil {
				return 0, err
			}
			off += n
		}
		return off, nil

	case reflect.Struct:
		off := 0
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			if t.Field(i).Tag.Get("saga") == "-" {
				continue
			}
			if f := v.Field(i); f.CanSet() {
				n, err := unmarshalValue(data[off:], f)
				if err != nil {
					return 0, err
				}
				// glog.Infof("unmarshal %s: %s off %d", t.Field(i).Name, f.Type(), n)
				off += n
			}
		}
		return off, nil

	case reflect.Uint8:
		if len(data) < 1 {
			return 0, io.ErrShortBuffer
		}
		v.SetUint(uint64(data[0]))
		return 1, nil

	case reflect.Uint16:
		if len(data) < 2 {
			return 0, io.ErrShortBuffer
		}
		v.SetUint(uint64(binary.BigEndian.Uint16(data[:2])))
		return 2, nil

	case reflect.Uint32:
		if len(data) < 4 {
			return 0, io.ErrShortBuffer
		}
		v.SetUint(uint64(binary.BigEndian.Uint32(data[:4])))
		return 4, nil

	case reflect.Uint64:
		if len(data) < 8 {
			return 0, io.ErrShortBuffer
		}
		v.SetUint(binary.BigEndian.Uint64(data[:8]))
		return 8, nil

	case reflect.String:
		l := data[0]
		if len(data) < int(l+1) {
			return 0, io.ErrShortBuffer
		}
		if l > 0 {
			v.SetString(strings.Trim(string(data[1:l+1]), "\x00"))
		}
		return int(l + 1), nil
	}
	return 0, fmt.Errorf("unsupport type `%s`", v.Type())
}

// Marshal 序列化
func Marshal(w io.Writer, obj interface{}) (int, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	if !v.IsValid() {
		return 0, fmt.Errorf("%v is invalid", obj)
	}
	return marshalValue(w, v)
}

// Unmarshal 反序列化
func Unmarshal(data []byte, v interface{}) (int, error) {
	objV := reflect.ValueOf(v)
	if objV.Kind() != reflect.Ptr || objV.Elem().Kind() != reflect.Struct {
		return 0, fmt.Errorf("%v must be struct pointer", v)
	}
	return unmarshalValue(data, objV.Elem())
}

// Packet 协议包
func Packet(w io.Writer, size, id uint16, v ...interface{}) error {
	binary.Write(w, binary.BigEndian, size)
	binary.Write(w, binary.BigEndian, id)
	off := 2
	for _, o := range v {
		n, err := Marshal(w, o)
		if err != nil {
			return err
		}
		off += n
	}
	if off < int(size) {
		w.Write(bytes.Repeat(zero, int(size)-off))
	}
	return nil
}
