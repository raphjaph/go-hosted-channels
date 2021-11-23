package hcwire

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/stretchr/testify/assert"
)

func getTestInvokeHC() *InvokeHostedChannel {
	var genesisHash [32]byte
	hash, _ := hex.DecodeString("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	copy(genesisHash[:], hash)

	return &InvokeHostedChannel{
		ChainHash:          genesisHash,
		RefundScriptPubKey: []byte{8},
		Secret:             []byte{10},
	}
}

func getTestInitHC() *InitHostedChannel {
	return &InitHostedChannel{
		MaxHTLCValueInFlightMSat:           100000000,
		HTLCMinimumMSat:                    1000,
		MaxAcceptedHTLCs:                   30,
		ChannelCapacityMSat:                1000000000,
		LiabilityDeadlineBlockdays:         360,
		MinimalOnChainRefundAmountSatoshis: 100000,
		InitialClientBalanceMSat:           0,
		Features:                           []byte{},
	}
}

func getTestLassCSS() *LastCrossSignedState {
	var sig [64]byte

	return &LastCrossSignedState{
		LastRefundScriptPubKey: []byte{5},
		InitHostedChannel:      *getTestInitHC(),
		Blockday:               120,
		LocalBalanceMSat:       100000,
		RemoteBalanceMSat:      1000000,
		LocalUpdates:           2,
		RemoteUpdates:          3,
		IncomingHTLCs:          []lnwire.UpdateAddHTLC{},
		OutgoingHTLCs:          []lnwire.UpdateAddHTLC{},
		RemoteSigOfLocal:       sig,
		LocalSigOfRemote:       sig,
	}
}

func getTestStateUpdate() *StateUpdate {
	var sig [64]byte

	return &StateUpdate{
		Blockday:         122,
		LocalUpdates:     1,
		RemoteUpdates:    2,
		LocalSigOfRemote: sig,
	}
}

func getTestStateOverride() *StateOverride {
	var sig [64]byte

	return &StateOverride{
		Blockday:         122,
		LocalBalanceMSat: 88888888,
		LocalUpdates:     1,
		RemoteUpdates:    2,
		LocalSigOfRemote: sig,
	}
}

func TestInvokeHostedChannel(t *testing.T) {
	invokeHC := getTestInvokeHC()

	// learning stuff: &bytes.Buffer{} == new(bytes.Buffer)
	b := new(bytes.Buffer)
	WriteMessage(b, invokeHC, 1)

	// take wire message b and wrap in reader
	r := bytes.NewReader(b.Bytes())
	msg, err := ReadMessage(r, 1)
	if err != nil {
		fmt.Println("error: ", err)
	}

	decodedInvokeHC, ok := msg.(*InvokeHostedChannel)
	if !ok {
		fmt.Println("could not do type assertion")
	}

	assert.Equal(t, invokeHC, decodedInvokeHC)
}

func TestInitHostedChannel(t *testing.T) {
	initHC := getTestInitHC()

	b := new(bytes.Buffer)
	WriteMessage(b, initHC, 1)

	r := bytes.NewReader(b.Bytes())
	msg, err := ReadMessage(r, 1)
	if err != nil {
		fmt.Println("error: ", err)
	}

	decodedInitHC, ok := msg.(*InitHostedChannel)
	if !ok {
		fmt.Println("could not do type assertion")
	}

	assert.Equal(t, initHC, decodedInitHC)
}

func TestLastCrossedSignedState(t *testing.T) {
	lastCSS := getTestLassCSS()

	b := new(bytes.Buffer)
	WriteMessage(b, lastCSS, 1)

	r := bytes.NewReader(b.Bytes())
	msg, err := ReadMessage(r, 1)
	if err != nil {
		fmt.Println("error: ", err)
	}

	decodedLastCSS, ok := msg.(*LastCrossSignedState)
	if !ok {
		fmt.Println("could not do type assertion")
	}

	assert.Equal(t, lastCSS, decodedLastCSS)

}

func TestStateUpdate(t *testing.T) {
	stateUpdate := getTestStateUpdate()

	b := new(bytes.Buffer)
	WriteMessage(b, stateUpdate, 1)

	r := bytes.NewReader(b.Bytes())
	msg, err := ReadMessage(r, 1)
	if err != nil {
		fmt.Println("error: ", err)
	}

	decodedStateUpdate, ok := msg.(*StateUpdate)
	if !ok {
		fmt.Println("could not do type assertion")
	}

	assert.Equal(t, stateUpdate, decodedStateUpdate)

}

func TestStateOverride(t *testing.T) {
	stateOverride := getTestStateOverride()

	b := new(bytes.Buffer)
	WriteMessage(b, stateOverride, 1)

	r := bytes.NewReader(b.Bytes())
	msg, err := ReadMessage(r, 1)
	if err != nil {
		fmt.Println("error: ", err)
	}

	decodedStateOverride, ok := msg.(*StateOverride)
	if !ok {
		fmt.Println("could not do type assertion")
	}

	assert.Equal(t, stateOverride, decodedStateOverride)

}
