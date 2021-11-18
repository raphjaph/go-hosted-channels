// code derived from https://github.com/lightningnetwork/lnd/blob/master/lnwire

package hcwire

import (
	"bytes"
	"io"
)

type InitHostedChannel struct {
	MaxHTLCValueInFlightMSAT           uint64
	HTLCMinimumMSAT                    uint64
	MaxAcceptedHTLCs                   uint16
	ChannelCapacityMSAT                uint64
	LiabilityDeadlineBlockdays         uint16
	MinimalOnChainRefundAmountSatoshis uint64
	InitialClientBalanceMSAT           uint64
	Features                           []byte
}

func NewInitHostedChannel() *InitHostedChannel {
	return &InitHostedChannel{}
}

var _ Message = (*InitHostedChannel)(nil)

func (c *InitHostedChannel) Decode(r io.Reader, pver uint32) error {

	return nil
}

/*
func (t *InitHostedChannel) decode(b []byte) {
	r := NewReader(b)
	t.MaxHTLCValueInFlightMSAT = r.readUint64()
	t.HTLCMinimumMSAT = r.readUint64()
	t.MaxAcceptedHTLCs = r.readUint16()
	t.ChannelCapacityMSAT = r.readUint64()
	t.LiabilityDeadlineBlockdays = r.readUint16()
	t.MinimalOnChainRefundAmountSatoshis = r.readUint64()
	t.InitialClientBalanceMSAT = r.readUint64()
	t.Features = r.readDynamic()
}
*/

func (c *InitHostedChannel) Encode(w *bytes.Buffer, pver uint32) error {

	return nil
}

func (c *InitHostedChannel) MsgType() MessageType {
	return MsgInitHostedChannel
}
