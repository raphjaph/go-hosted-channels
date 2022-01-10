package hcwire

import (
	"bytes"
	"io"

	"github.com/lightningnetwork/lnd/lnwire"
)

type UpdateFulfillHTLC struct {
	lnwire.UpdateFulfillHTLC
}

func NewUpdateFulfillHTLC() *UpdateFulfillHTLC {
	return &UpdateFulfillHTLC{}
}

var _ Message = (*UpdateFulfillHTLC)(nil)

func (c *UpdateFulfillHTLC) Decode(r io.Reader, pver uint32) error {
	return c.UpdateFulfillHTLC.Decode(r, pver)
}

func (c *UpdateFulfillHTLC) Encode(buf *bytes.Buffer, pver uint32) error {
	return c.UpdateFulfillHTLC.Encode(buf, pver)
}

func (c *UpdateFulfillHTLC) MsgType() MessageType {
	return MsgUpdateFulfillHTLC
}
