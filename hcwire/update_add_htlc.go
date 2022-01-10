package hcwire

import (
	"bytes"
	"io"

	"github.com/lightningnetwork/lnd/lnwire"
)

// wrapper around lnwire.UpdateAddHTLC
// not sure if this is the way to do this
type UpdateAddHTLC struct {
	lnwire.UpdateAddHTLC
}

func NewUpdateAddHTLC() *UpdateAddHTLC {
	return &UpdateAddHTLC{}
}

var _ Message = (*UpdateAddHTLC)(nil)

func (c *UpdateAddHTLC) Decode(r io.Reader, pver uint32) error {
	return c.UpdateAddHTLC.Decode(r, pver)
}

func (c *UpdateAddHTLC) Encode(buf *bytes.Buffer, pver uint32) error {
	return c.UpdateAddHTLC.Encode(buf, pver)
}

func (c *UpdateAddHTLC) MsgType() MessageType {
	return MsgUpdateAddHTLC
}
