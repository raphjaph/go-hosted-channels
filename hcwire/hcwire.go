package hcwire

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/lightningnetwork/lnd/lnwire"
)

// assumes the length field is a uint16
// copied/modified from ReadVarBytes in "github.com/btcsuite/btcd/wire"
func ReadVarBytes(r io.Reader, maxAllowed uint16, fieldName string) ([]byte, error) {
	var length uint16
	if err := ReadElement(r, &length); err != nil {
		return nil, err
	}

	if length > maxAllowed {
		str := fmt.Sprintf("%s is larger than the max allowed size "+
			"[length %d, max %d]", fieldName, length, maxAllowed)
		return nil, fmt.Errorf("ReadVarBytes: %v", str)
	}

	b := make([]byte, length)
	_, err := io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// copied/modified from github.com/lightningnetwork/lnd/lnwire
// Important: always give this function a pointer to the data structure!
func ReadElement(r io.Reader, element interface{}) error {
	switch e := element.(type) {
	case *bool:
		var b [1]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}

		if b[0] == 1 {
			*e = true
		}

	case *uint16:
		var b [2]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = binary.BigEndian.Uint16(b[:])

	case *uint32:
		var b [4]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = binary.BigEndian.Uint32(b[:])

	case *uint64:
		var b [8]byte
		if _, err := io.ReadFull(r, b[:]); err != nil {
			return err
		}
		*e = binary.BigEndian.Uint64(b[:])

	default:
		return fmt.Errorf("unknown type in ReadElement: %T", e)
	}

	return nil
}

// writes a uint16 for the length field
func WriteVarBytes(buf *bytes.Buffer, bytes []byte) error {
	length := len(bytes)

	if length > 65535 {
		return fmt.Errorf("can not encode byte array with length larger than uint16 (65535): %v", length)
	}

	if err := lnwire.WriteUint16(buf, uint16(length)); err != nil {
		return err
	}

	if _, err := buf.Write(bytes); err != nil {
		return err
	}

	return nil
}
