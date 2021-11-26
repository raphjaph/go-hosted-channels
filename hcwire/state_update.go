package hcwire

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lightningnetwork/lnd/lnwire"
)

type StateUpdate struct {
	Blockday         uint32
	LocalUpdates     uint32
	RemoteUpdates    uint32
	LocalSigOfRemote [64]byte // sig of remote last_crossed_signed_state
}

func NewStateUpdate() *StateUpdate {
	return &StateUpdate{}
}

var _ Message = (*StateUpdate)(nil)

func (c *StateUpdate) Decode(r io.Reader, pver uint32) error {
	if err := ReadElement(r, &c.Blockday); err != nil {
		return err
	}

	if err := ReadElement(r, &c.LocalUpdates); err != nil {
		return err
	}

	if err := ReadElement(r, &c.RemoteUpdates); err != nil {
		return err
	}

	if _, err := io.ReadFull(r, c.LocalSigOfRemote[:]); err != nil {
		return fmt.Errorf("could not parse local_sig_of_remote: %v", err)
	}

	return nil
}

func (c *StateUpdate) Encode(buf *bytes.Buffer, pver uint32) error {
	if err := lnwire.WriteUint32(buf, c.Blockday); err != nil {
		return err
	}

	if err := lnwire.WriteUint32(buf, c.LocalUpdates); err != nil {
		return err
	}

	if err := lnwire.WriteUint32(buf, c.RemoteUpdates); err != nil {
		return err
	}

	if _, err := buf.Write(c.LocalSigOfRemote[:]); err != nil {
		return err
	}

	return nil
}

func (c *StateUpdate) MsgType() MessageType {
	return MsgStateUpdate
}
