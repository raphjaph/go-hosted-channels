package main

import (
	"encoding/hex"

	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/raphjaph/go-hosted-channels/hcwire"
)

type Channel struct {
	ChannelID            lnwire.ChannelID
	PeerID               string
	InitHostedChannel    hcwire.InitHostedChannel    // parameters of the channel: size, refund_addr, etc.
	LastCrossSignedState hcwire.LastCrossSignedState // current state; similar to committment transaction + revokation key
}

func getGenesisHash() [32]byte {
	var genesisHash [32]byte
	hash, _ := hex.DecodeString("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	copy(genesisHash[:], hash)
	return genesisHash
}

func getRandomShortChannelID() string {
	return "1x2x3"
}

func getNodeKey() {

}
