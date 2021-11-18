package hcwire

import (
	"github.com/lightningnetwork/lnd/lnwire"
)

type LastCrossSignedState struct {
	IsHost             bool
	RefundScriptPubKey []byte
	InitHostedChannel  InitHostedChannel
	Blockday           uint32
	LocalBalanceMSAT   uint64
	RemoteBalanceMSAT  uint64
	LocalUpdates       uint32
	RemoteUpdates      uint32
	IncomingHTLCs      []lnwire.UpdateAddHTLC
	OutgoingHTLC       []lnwire.UpdateAddHTLC
	RemoteSigOfLocal   [64]byte
	LocalSigOfRemote   [64]byte
}
