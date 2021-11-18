package hcwire

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvokeHostedChannel(t *testing.T) {
	var genesisHash [32]byte
	hash, _ := hex.DecodeString("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	copy(genesisHash[:], hash)

	invokeHC := &InvokeHostedChannel{
		ChainHash:          genesisHash,
		RefundScriptPubKey: []byte{8},
		Secret:             []byte{10},
	}

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

	initHC := &InitHostedChannel{
		MaxHTLCValueInFlightMSat:           100000000,
		HTLCMinimumMSat:                    1000,
		MaxAcceptedHTLCs:                   30,
		ChannelCapacityMSat:                1000000000,
		LiabilityDeadlineBlockdays:         360,
		MinimalOnChainRefundAmountSatoshis: 100000,
		InitialClientBalanceMSat:           0,
		Features:                           []byte{},
	}

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
