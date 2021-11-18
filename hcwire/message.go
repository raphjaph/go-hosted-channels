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
	MsgInvokeHostedChannel    MessageType = 65535
	MsgInitHostedChannel      MessageType = 65533
	MsgLastCrossedSignedState             = 65531
	MsgStateUpdate                        = 65529
	MsgStateOverride                      = 65527
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
	default:
		return "<unknown>"
	}
}

type UnknownMessage struct {
	messageType MessageType
}

func (u *UnknownMessage) Error() string {
	return fmt.Sprintf("unable to parse message of unknown type: %v", u.messageType)
}

type Serializable interface {
	Decode(io.Reader, uint32) error
	Encode(*bytes.Buffer, uint32) error
}

type Message interface {
	Serializable
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

		/*
			case MsgLastCrossedSignedState:
				msg = &LastCrossSignedState{}
			case MsgStateUpdate:
				msg = &StateUpdate{}
			case MsgStateOverride:
				msg = &StateOverride{}
		*/

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
