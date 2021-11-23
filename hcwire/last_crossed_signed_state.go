package hcwire

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lightningnetwork/lnd/lnwire"
)

type LastCrossSignedState struct {
	IsHost                 bool
	LastRefundScriptPubKey []byte
	InitHostedChannel      InitHostedChannel
	Blockday               uint32
	LocalBalanceMSat       uint64
	RemoteBalanceMSat      uint64
	LocalUpdates           uint32
	RemoteUpdates          uint32
	IncomingHTLCs          []lnwire.UpdateAddHTLC
	OutgoingHTLCs          []lnwire.UpdateAddHTLC
	RemoteSigOfLocal       [64]byte
	LocalSigOfRemote       [64]byte
}

func NewLastCrossedSignedState() *LastCrossSignedState {
	return &LastCrossSignedState{}
}

var _ Message = (*LastCrossSignedState)(nil)

func (c *LastCrossSignedState) Decode(r io.Reader, pver uint32) (err error) {

	c.LastRefundScriptPubKey, err = ReadVarBytes(r, 34, "last_refund_scriptpubkey")
	if err != nil {
		return err
	}

	initHC := NewInitHostedChannel()
	err = initHC.Decode(r, 1)
	if err != nil {
		return err
	}
	c.InitHostedChannel = *initHC

	if err := ReadElement(r, &c.Blockday); err != nil {
		return err
	}

	if err := ReadElement(r, &c.LocalBalanceMSat); err != nil {
		return err
	}

	if err := ReadElement(r, &c.RemoteBalanceMSat); err != nil {
		return err
	}

	if err := ReadElement(r, &c.LocalUpdates); err != nil {
		return err
	}

	if err := ReadElement(r, &c.RemoteUpdates); err != nil {
		return err
	}

	var num uint16
	if err := ReadElement(r, &num); err != nil {
		return err
	}
	c.IncomingHTLCs = make([]lnwire.UpdateAddHTLC, num)
	for _, incoming := range c.IncomingHTLCs {
		err := incoming.Decode(r, 1)
		if err != nil {
			return err
		}
	}

	if err := lnwire.ReadElement(r, &num); err != nil {
		return err
	}
	c.OutgoingHTLCs = make([]lnwire.UpdateAddHTLC, num)
	for _, outgoing := range c.OutgoingHTLCs {
		err := outgoing.Decode(r, 1)
		if err != nil {
			return err
		}
	}

	_, err = io.ReadFull(r, c.RemoteSigOfLocal[:])
	if err != nil {
		return fmt.Errorf("could not parse remote_sig_of_local: %v", err)
	}

	_, err = io.ReadFull(r, c.LocalSigOfRemote[:])
	if err != nil {
		return fmt.Errorf("could not parse local_sig_of_remote: %v", err)
	}

	return err
}

func (c *LastCrossSignedState) Encode(buf *bytes.Buffer, pver uint32) (err error) {
	if err := WriteVarBytes(buf, c.LastRefundScriptPubKey); err != nil {
		return err
	}

	if err := c.InitHostedChannel.Encode(buf, 1); err != nil {
		return err
	}

	if err := lnwire.WriteUint32(buf, c.Blockday); err != nil {
		return err
	}

	if err := lnwire.WriteUint64(buf, c.LocalBalanceMSat); err != nil {
		return err
	}

	if err := lnwire.WriteUint64(buf, c.RemoteBalanceMSat); err != nil {
		return err
	}

	if err := lnwire.WriteUint32(buf, c.LocalUpdates); err != nil {
		return err
	}

	if err := lnwire.WriteUint32(buf, c.RemoteUpdates); err != nil {
		return err
	}

	if err := lnwire.WriteUint16(buf, uint16(len(c.IncomingHTLCs))); err != nil {
		return err
	}
	for _, incoming := range c.IncomingHTLCs {
		if err := incoming.Encode(buf, 1); err != nil {
			return err
		}
	}

	if err := lnwire.WriteUint16(buf, uint16(len(c.OutgoingHTLCs))); err != nil {
		return err
	}
	for _, outgoing := range c.OutgoingHTLCs {
		if err := outgoing.Encode(buf, 1); err != nil {
			return err
		}
	}

	if _, err := buf.Write(c.RemoteSigOfLocal[:]); err != nil {
		return err
	}

	if _, err := buf.Write(c.LocalSigOfRemote[:]); err != nil {
		return err
	}

	return err
}

func (c *LastCrossSignedState) MsgType() MessageType {
	return MsgLastCrossedSignedState
}
