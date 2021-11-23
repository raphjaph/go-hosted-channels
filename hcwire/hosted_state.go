package hcwire

import (
	"bytes"
	"io"

	"github.com/lightningnetwork/lnd/lnwire"
)

//
type HostedState struct {
	ChannelID            lnwire.ChannelID
	NumNextLocalUpdates  uint16
	NextLocalUpdates     []byte // have to define new type here; some other fields still missing here as well
	LastCrossSignedState LastCrossSignedState
}

func NewHostedState() *HostedState {
	return &HostedState{}
}

var _ Message = (*HostedState)(nil)

func (c *HostedState) Decode(r io.Reader, pver uint32) error {

	return nil
}

func (c *HostedState) Encode(w *bytes.Buffer, pver uint32) error {
	return nil
}

func (c *HostedState) MsgType() MessageType {
	return 0
}
