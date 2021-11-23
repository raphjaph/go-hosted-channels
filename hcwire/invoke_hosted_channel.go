// Copyright (C) 2015-2017 The Lightning Network Developers
// code derived from https:// github.com/lightningnetwork/lnd/blob/master/lnwire/

package hcwire

import (
	"bytes"
	"fmt"
	"io"
)

type InvokeHostedChannel struct {
	ChainHash          [32]byte
	RefundScriptPubKey []byte
	Secret             []byte // optional data which can be used by Host to tweak channel parameters (non-zero initial Client balance, larger capacity, only allow Clients with secrets etc)
}

func NewInvokeHostedChannel() *InvokeHostedChannel {
	return &InvokeHostedChannel{}
}

// A compile time check to ensure InvokeHostedChannel implements the hcwire.Message
// interface.
var _ Message = (*InvokeHostedChannel)(nil)

func (c *InvokeHostedChannel) Decode(r io.Reader, pver uint32) error {
	_, err := io.ReadFull(r, c.ChainHash[:])
	if err != nil {
		return fmt.Errorf("could not parse chain_hash: %v", err)
	}

	// p2wsh is 34 bytes long (max allowed length)
	c.RefundScriptPubKey, err = ReadVarBytes(r, 34, "refund_scriptpubkey")
	if err != nil {
		return err
	}
	// read the custom TLV field (secret)
	// Secret should not be longer than 64 bytes
	c.Secret, err = ReadVarBytes(r, 64, "secret")

	return err
}

func (c *InvokeHostedChannel) Encode(buf *bytes.Buffer, pver uint32) error {
	if _, err := buf.Write(c.ChainHash[:]); err != nil {
		return err
	}

	if err := WriteVarBytes(buf, c.RefundScriptPubKey); err != nil {
		return err
	}

	return WriteVarBytes(buf, c.Secret)
}

func (c *InvokeHostedChannel) MsgType() MessageType {
	return MsgInvokeHostedChannel
}
