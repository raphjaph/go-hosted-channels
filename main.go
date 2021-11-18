package main

import (
	"bytes"
	"encoding/hex"

	"github.com/fiatjaf/lightningd-gjson-rpc/plugin"
	"github.com/raphjaph/go-hosted-channels/hcwire"
)

var continueHTLC = map[string]interface{}{"result": "continue"}

type Channel struct {
	PeerID               string
	InitHostedChannel    hcwire.InitHostedChannel
	LastCrossSignedState hcwire.LastCrossSignedState
}

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
			p.Log("listening for hooks, subscriptions and RPC methods")
		},
	}

	p.Run()
}

func hcInvoke(p *plugin.Plugin, params plugin.Params) (interface{}, int, error) {
	/*
		nodeId := params.Get("node_id").String()
		refundAddr := params.Get("refund_address").String()
		secret := params.Get("secret").String()

		payload := hcwire.NewInvokeHostedChannel()
		payload.Encode()

		p.Client.Call("sendcustommsg", nodeId)
	*/
	return nil, 0, nil
}

func handleHTLCAccepted(p *plugin.Plugin, params plugin.Params) (resp interface{}) {
	return continueHTLC
}

func handleCustomMsg(p *plugin.Plugin, params plugin.Params) (resp interface{}) {

	peer := params.Get("peer_id").String()

	p.Log(params.Get("@pretty").Raw)

	// this seems a bit convoluted
	// from hex string to []byte
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
		ihc, ok := msg.(*hcwire.InvokeHostedChannel)
		if !ok {
			p.Log("unable to cast to InvokeHostedChannel")
			return continueHTLC
		}
		p.Logf("got %v from %v", ihc.MsgType(), peer)
		p.Log(ihc.ChainHash)

	case hcwire.MsgInitHostedChannel:
		p.Logf("got %v from %v", msg.MsgType(), peer)
	}

	p.Log("handeling msg type: ", msg.MsgType(), " from peer: ", peer, " with content: ", payload[:])
	return continueHTLC
}
