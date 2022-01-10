// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Copyright (C) 2015-2017 The Lightning Network Developers
// code derived from https:// github.com/lightningnetwork/lnd/blob/master/lnwire/message.go

package hcwire

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/lightningnetwork/lnd/lnwire"
)

type MessageType uint16

const (
	MsgInvokeHostedChannel     MessageType = 65535
	MsgInitHostedChannel       MessageType = 65533
	MsgLastCrossedSignedState              = 65531
	MsgStateUpdate                         = 65529
	MsgStateOverride                       = 65527
	MsgInvoiceForward                      = 65525
	MsgUpdateAddHTLC                       = 63505
	MsgUpdateFulfillHTLC                   = 63503
	MsgUpdateFailHTLC                      = 63501
	MsgUpdateFailMalformedHTLC             = 63499
)

func (t MessageType) String() string {
	switch t {
	case MsgInvokeHostedChannel:
		return "invoke_hosted_channel"
	case MsgInitHostedChannel:
		return "init_hosted_channel"
	case MsgLastCrossedSignedState:
		return "last_crossed_signed_state"
	case MsgStateUpdate:
		return "state_update"
	case MsgStateOverride:
		return "state_override"
	case MsgUpdateAddHTLC:
		return "update_add_htlc"
	case MsgUpdateFulfillHTLC:
		return "update_fulfill_htlc"
	case MsgUpdateFailHTLC:
		return "update_fail_htlc"
	case MsgUpdateFailMalformedHTLC:
		return "update_fail_malformed_htlc"
	default:
		return "<unknown>"
	}
}

// Not sure if I need this:
/*
type UnknownMessage struct {
	messageType MessageType
}

func (u *UnknownMessage) Error() string {
	return fmt.Sprintf("unable to parse message of unknown type: %v", u.messageType)
}
*/

// TODO: how to wrap the lnwire.Message; here I'm just copying it because I got annoyed with all the (re)casting of types
type Message interface {
	lnwire.Serializable
	MsgType() MessageType
}

// return a pointer to the correct message type struct
func makeEmptyMessage(msgType MessageType) (Message, error) {
	var msg Message

	switch msgType {
	case MsgInvokeHostedChannel:
		msg = &InvokeHostedChannel{}
	case MsgInitHostedChannel:
		msg = &InitHostedChannel{}
	case MsgLastCrossedSignedState:
		msg = &LastCrossSignedState{}
	case MsgStateUpdate:
		msg = &StateUpdate{}
	case MsgStateOverride:
		msg = &StateOverride{}
	case MsgUpdateAddHTLC:
		msg = &UpdateAddHTLC{}
	case MsgUpdateFulfillHTLC:
		msg = &UpdateFulfillHTLC{}
	default:
		return nil, fmt.Errorf("not a hosted channel message")
	}
	return msg, nil
}

func WriteMessage(buf *bytes.Buffer, msg Message, pver uint32) (int, error) {
	// Record the size of the bytes already written in buffer.
	oldByteSize := buf.Len()

	// cleanBrokenBytes is a helper closure that helps reset the buffer to
	// its original state. It truncates all the bytes written in current
	// scope.
	var cleanBrokenBytes = func(b *bytes.Buffer) int {
		b.Truncate(oldByteSize)
		return 0
	}

	// Write the message type.
	var mType [2]byte
	binary.BigEndian.PutUint16(mType[:], uint16(msg.MsgType()))
	msgTypeBytes, err := buf.Write(mType[:])
	if err != nil {
		return cleanBrokenBytes(buf), lnwire.ErrorWriteMessageType(err)
	}

	// Use the write buffer to encode our message.
	if err := msg.Encode(buf, pver); err != nil {
		return cleanBrokenBytes(buf), lnwire.ErrorEncodeMessage(err)
	}

	// Enforce maximum overall message payload. The write buffer now has
	// the size of len(originalBytes) + len(payload) + len(type). We want
	// to enforce the payload here, so we subtract it by the length of the
	// type and old bytes.
	lenp := buf.Len() - oldByteSize - msgTypeBytes
	if lenp > lnwire.MaxMsgBody {
		return cleanBrokenBytes(buf), lnwire.ErrorPayloadTooLarge(lenp)
	}

	return buf.Len() - oldByteSize, nil
}

func ReadMessage(r io.Reader, pver uint32) (Message, error) {
	// first two bites custom message type
	var mType [2]byte
	if _, err := io.ReadFull(r, mType[:]); err != nil {
		return nil, err
	}

	msgType := MessageType(binary.BigEndian.Uint16(mType[:]))

	msg, err := makeEmptyMessage(msgType)
	if err != nil {
		return nil, err
	}

	if err := msg.Decode(r, pver); err != nil {
		return nil, err
	}

	return msg, nil
}

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
