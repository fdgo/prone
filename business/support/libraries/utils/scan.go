package utils

import (
	"encoding"
	"errors"
	"fmt"
	"strconv"
)

func Scan(b []byte, v interface{}) error {
	switch v := v.(type) {
	case nil:
		return errors.New("Scan nil")
	case *string:
		*v = string(b)
		return nil
	case *[]byte:
		*v = b
		return nil
	case *int:
		var err error
		*v, err = strconv.Atoi(string(b))
		return err
	case *int8:
		n, err := strconv.ParseInt(string(b), 10, 8)
		if err != nil {
			return err
		}
		*v = int8(n)
		return nil
	case *int16:
		n, err := strconv.ParseInt(string(b), 10, 16)
		if err != nil {
			return err
		}
		*v = int16(n)
		return nil
	case *int32:
		n, err := strconv.ParseInt(string(b), 10, 32)
		if err != nil {
			return err
		}
		*v = int32(n)
		return nil
	case *int64:
		n, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *uint:
		n, err := strconv.ParseUint(string(b), 10, 64)
		if err != nil {
			return err
		}
		*v = uint(n)
		return nil
	case *uint8:
		n, err := strconv.ParseUint(string(b), 10, 8)
		if err != nil {
			return err
		}
		*v = uint8(n)
		return nil
	case *uint16:
		n, err := strconv.ParseUint(string(b), 10, 16)
		if err != nil {
			return err
		}
		*v = uint16(n)
		return nil
	case *uint32:
		n, err := strconv.ParseUint(string(b), 10, 32)
		if err != nil {
			return err
		}
		*v = uint32(n)
		return nil
	case *uint64:
		n, err := strconv.ParseUint(string(b), 10, 64)
		if err != nil {
			return err
		}
		*v = n
		return nil
	case *float32:
		n, err := strconv.ParseFloat(string(b), 32)
		if err != nil {
			return err
		}
		*v = float32(n)
		return err
	case *float64:
		var err error
		*v, err = strconv.ParseFloat(string(b), 64)
		return err
	case *bool:
		*v = len(b) == 1 && b[0] == '1'
		return nil
	case encoding.BinaryUnmarshaler:
		return v.UnmarshalBinary(b)
	default:
		return fmt.Errorf(
			"redis: can't unmarshal %T (consider implementing BinaryUnmarshaler)", v)
	}
}
