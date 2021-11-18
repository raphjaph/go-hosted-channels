// code derived from https://github.com/lightningnetwork/lnd/blob/master/lnwire

package hcwire

import (
	"bytes"
	"io"

	"github.com/btcsuite/btcd/wire"
	"github.com/lightningnetwork/lnd/lnwire"
)

type InitHostedChannel struct {
	MaxHTLCValueInFlightMSat           uint64
	HTLCMinimumMSat                    uint64
	MaxAcceptedHTLCs                   uint16
	ChannelCapacityMSat                uint64
	LiabilityDeadlineBlockdays         uint16
	MinimalOnChainRefundAmountSatoshis uint64
	InitialClientBalanceMSat           uint64
	Features                           []byte
}

func NewInitHostedChannel() *InitHostedChannel {
	return &InitHostedChannel{}
}

var _ Message = (*InitHostedChannel)(nil)

func (c *InitHostedChannel) Decode(r io.Reader, pver uint32) error {
	//TODO: make for loop or something more elegant

	if err := lnwire.ReadElement(r, &c.MaxHTLCValueInFlightMSat); err != nil {
		return err
	}

	if err := lnwire.ReadElement(r, &c.HTLCMinimumMSat); err != nil {
		return err
	}

	if err := lnwire.ReadElement(r, &c.MaxAcceptedHTLCs); err != nil {
		return err
	}

	if err := lnwire.ReadElement(r, &c.ChannelCapacityMSat); err != nil {
		return err
	}

	if err := lnwire.ReadElement(r, &c.LiabilityDeadlineBlockdays); err != nil {
		return err
	}

	if err := lnwire.ReadElement(r, &c.MinimalOnChainRefundAmountSatoshis); err != nil {
		return err
	}

	if err := lnwire.ReadElement(r, &c.InitialClientBalanceMSat); err != nil {
		return err
	}

	// TODO:  what is max size of features field
	var err error
	c.Features, err = wire.ReadVarBytes(r, 1, 13, "features")

	return err
}

func (c *InitHostedChannel) Encode(w *bytes.Buffer, pver uint32) error {

	if err := lnwire.WriteUint64(w, c.MaxHTLCValueInFlightMSat); err != nil {
		return err
	}

	if err := lnwire.WriteUint64(w, c.HTLCMinimumMSat); err != nil {
		return err
	}

	if err := lnwire.WriteUint16(w, c.MaxAcceptedHTLCs); err != nil {
		return err
	}

	if err := lnwire.WriteUint64(w, c.ChannelCapacityMSat); err != nil {
		return err
	}

	if err := lnwire.WriteUint16(w, c.LiabilityDeadlineBlockdays); err != nil {
		return err
	}

	if err := lnwire.WriteUint64(w, c.MinimalOnChainRefundAmountSatoshis); err != nil {
		return err
	}

	if err := lnwire.WriteUint64(w, c.InitialClientBalanceMSat); err != nil {
		return err
	}

	err := wire.WriteVarBytes(w, 1, c.Features)

	return err
}

func (c *InitHostedChannel) MsgType() MessageType {
	return MsgInitHostedChannel
}
