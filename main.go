package main

import (
	"bytes"
	"encoding/hex"

	"github.com/fiatjaf/lightningd-gjson-rpc/plugin"
	"github.com/raphjaph/go-hosted-channels/hcwire"
)

var continueHTLC = map[string]interface{}{"result": "continue"}
var failHTLC = map[string]interface{}{"result": "fail", "failure_message": "2002"} // TODO: hosted channel specific error codes
var resolveHTLC = map[string]interface{}{"result": "resolve", "payment_key": "0000000000000000000000000000000000000000000000000000000000000000"}

func main() {
	p := plugin.Plugin{
		Name:    "hosted-channels",
		Version: "v0.0.1",
		Options: []plugin.Option{
			{
				Name:        "hosted-channel-size",
				Type:        "int",
				Default:     1000000,
				Description: "The default size in sats of a hosted channel.",
			},
		},

		// do something asynchronously; lightnind doesn't wait for response
		Subscriptions: []plugin.Subscription{
			{
				Type: "invoice_payment",
				Handler: func(p *plugin.Plugin, params plugin.Params) {
					label := params.Get("invoice_payment.label").String()
					adjective := p.Args.Get("hosted-channel-size").String()
					p.Logf("%s payment received with label %s", adjective, label)
				},
			},
		},

		// do somehting but lightningd waits for response; synchronous
		Hooks: []plugin.Hook{
			{
				Type:    "custommsg",
				Handler: handleCustomMsg,
			},
			{
				Type:    "htlc_accepted",
				Handler: handleHTLCAccepted,
			},
		},

		RPCMethods: []plugin.RPCMethod{
			{
				Name:            "hc-invoke",
				Usage:           "node_id refund_address secret",
				Description:     "Invokes a new HC with remote nodeId, if accepted your node will be a Client side. Established HC is private by default.",
				LongDescription: "",
				Handler:         hcInvoke,
			},
		},

		OnInit: func(p *plugin.Plugin) {
			p.Log("hosted-channel plugin loaded")
		},
	}

	p.Run()
}

func hcInvoke(p *plugin.Plugin, params plugin.Params) (interface{}, int, error) {

	nodeId := params.Get("node_id").String()
	//refundAddr := params.Get("refund_address").String()
	//secret := params.Get("secret").String()

	var genesisHash [32]byte
	hash, _ := hex.DecodeString("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	copy(genesisHash[:], hash)

	invokeHC := &hcwire.InvokeHostedChannel{
		ChainHash:          genesisHash,
		RefundScriptPubKey: []byte{8},
		Secret:             []byte{10},
	}

	buf := new(bytes.Buffer)
	if _, err := hcwire.WriteMessage(buf, invokeHC, 1); err != nil {
		return nil, 1, err
	}
	payload := hex.EncodeToString(buf.Bytes())

	p.Client.Call("sendcustommsg", nodeId, payload)

	return nil, 0, nil
}

func handleHTLCAccepted(p *plugin.Plugin, params plugin.Params) (resp interface{}) {
	/*
		1. parse htlc
		2. scan peer channels for next hop
		3. construct new htlc and forward to peer
		4. wait for htlc_fulfill
		5. pass it on to hosted channel peer
	*/

	shortChannelID := params.Get("onion.short_channel_id").String()
	if shortChannelID == "0x0x0" {
		// if this node receiver
		return continueHTLC
	}

	p.Log("scid: ", shortChannelID)

	_, err := p.Client.Call("listpeers")
	if err != nil {
		p.Log("couldn't find any peers, error: ", err)
		return continueHTLC
	}

	return continueHTLC
}

func handleCustomMsg(p *plugin.Plugin, params plugin.Params) (resp interface{}) {

	peer := params.Get("peer_id").String()

	payload := params.Get("payload").String()
	b, err := hex.DecodeString(payload)
	if err != nil {
		p.Log("error decoding []byte from hex string: ", err)
		return continueHTLC
	}

	r := bytes.NewReader(b)
	msg, err := hcwire.ReadMessage(r, 1)
	if err != nil {
		p.Log("error reading custom message: ", err)
		return continueHTLC
	}

	switch msg.MsgType() {
	case hcwire.MsgInvokeHostedChannel:
		// Type assertions: https://golang.org/ref/spec#Type_assertions
		invokeHC, ok := msg.(*hcwire.InvokeHostedChannel)
		if !ok {
			p.Log("unable to assert InvokeHostedChannel type")
			return continueHTLC
		}
		p.Logf("got %v from %v", invokeHC.MsgType(), peer)

		p.Log("sending init_hosted_channel message")
		initHC := &hcwire.InitHostedChannel{
			MaxHTLCValueInFlightMSat:           100000000,
			HTLCMinimumMSat:                    1000,
			MaxAcceptedHTLCs:                   30,
			ChannelCapacityMSat:                1000000000,
			LiabilityDeadlineBlockdays:         360,
			MinimalOnChainRefundAmountSatoshis: 100000,
			InitialClientBalanceMSat:           0,
			Features:                           []byte{},
		}

		buf := new(bytes.Buffer)
		if _, err := hcwire.WriteMessage(buf, initHC, 1); err != nil {
			return continueHTLC
		}
		payload := hex.EncodeToString(buf.Bytes())

		p.Client.Call("sendcustommsg", peer, payload)

	case hcwire.MsgInitHostedChannel:
		initHC, ok := msg.(*hcwire.InitHostedChannel)
		if !ok {
			p.Log("unable to assert InitHostedChannel type")
			return continueHTLC
		}
		p.Logf("got %v from %v", initHC.MsgType(), peer)

	case hcwire.MsgLastCrossedSignedState:
		lastCSS, ok := msg.(*hcwire.LastCrossSignedState)
		if !ok {
			p.Log("unable to assert LastCrossSignedState type")
			return continueHTLC
		}
		p.Logf("got %v from %v", lastCSS.MsgType(), peer)

	case hcwire.MsgStateUpdate:
		stateUpdate, ok := msg.(*hcwire.StateUpdate)
		if !ok {
			p.Log("unable to assert StateUpdate type")
			return continueHTLC
		}
		p.Logf("got %v from %v", stateUpdate.MsgType(), peer)

	case hcwire.MsgStateOverride:
		stateOverride, ok := msg.(*hcwire.StateOverride)
		if !ok {
			p.Log("unable to assert StateOverrid type")
			return continueHTLC
		}
		p.Logf("got %v from %v", stateOverride.MsgType(), peer)

	default:
		p.Log("handeling msg type: ", msg.MsgType(), " from peer: ", peer, " with content: ", payload[:])
	}

	return continueHTLC
}
