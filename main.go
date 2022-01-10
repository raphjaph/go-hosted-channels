package main

/*
USEFUL RPC COMMANDS:
lightning-datastore – Command for storing (plugin) data
lightning-listdatastore – Command for listing (plugin) data
lightning-listforwards – Command showing all htlcs and their information
lightning-listsendpays – Low-level command for querying sendpay status
lightning-sendpay – Low-level command for sending a payment via a route
lightning-sendonion – Send a payment with a custom onion packet
lightning-createonion – Low-level command to create a custom onion
*/

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/fiatjaf/lightningd-gjson-rpc/plugin"
	"github.com/raphjaph/go-hosted-channels/hcwire"

	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/syndtr/goleveldb/leveldb"
)

var continueHTLC = map[string]interface{}{"result": "continue"}
var failHTLC = map[string]interface{}{"result": "fail", "failure_message": "2002"} // TODO: hosted channel specific error codes
var resolveHTLC = map[string]interface{}{"result": "resolve", "payment_key": "0000000000000000000000000000000000000000000000000000000000000000"}

var db *leveldb.DB

func main() {

	db, err := leveldb.OpenFile("hc-database", nil)
	if err != nil {
		fmt.Println("couldn't open database: ", err)
	}
	defer db.Close()

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

			{
				Name:            "hc-pay",
				Usage:           "node_id bolt11",
				Description:     "pay an invoice through a hosted channel",
				LongDescription: "",
				Handler:         hcPay,
			},
		},

		OnInit: func(p *plugin.Plugin) {
			p.Log("hosted-channel plugin loaded")
		},
	}

	p.Run()
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

	p.Logf("got %v from %v", msg.MsgType(), peer)

	switch msg.MsgType() {
	case hcwire.MsgInvokeHostedChannel:
		// Type assertions: https://golang.org/ref/spec#Type_assertions
		invokeHC, ok := msg.(*hcwire.InvokeHostedChannel)
		if !ok {
			p.Log("unable to assert InvokeHostedChannel type")
			return continueHTLC
		}
		p.Logf("got %v from %v", invokeHC.MsgType(), peer)

		// check if secret correct

		// create a channel in database wit initial parameters

		// reply with init message
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
		_, ok := msg.(*hcwire.InitHostedChannel)
		if !ok {
			p.Log("unable to assert InitHostedChannel type")
			return continueHTLC
		}

	case hcwire.MsgLastCrossedSignedState:
		_, ok := msg.(*hcwire.LastCrossSignedState)
		if !ok {
			p.Log("unable to assert LastCrossSignedState type")
			return continueHTLC
		}

	case hcwire.MsgStateUpdate:
		_, ok := msg.(*hcwire.StateUpdate)
		if !ok {
			p.Log("unable to assert StateUpdate type")
			return continueHTLC
		}

	case hcwire.MsgStateOverride:
		_, ok := msg.(*hcwire.StateOverride)
		if !ok {
			p.Log("unable to assert StateOverrid type")
			return continueHTLC
		}

	case hcwire.MsgUpdateAddHTLC:
		addHTLC, ok := msg.(*hcwire.UpdateAddHTLC)
		if !ok {
			p.Log("unable to assert UpdateAddHTLC type")
			return continueHTLC
		}

		// TEST: create own onion blob
		// hardcoded hop l2 -> l1 (pubkey)
		//hop := map[string]interface{}{"pubkey": "039c49ccbf7341e1829152d0067b1b47b06d596a21d9a7d640616ee2f990dc8f82", "payload": "00000067000003000100000000000003e800000075000000000000000000000000000000000000000000000000"}
		// HARD CODED: id from l1
		firstHop := map[string]interface{}{"id": "039c49ccbf7341e1829152d0067b1b47b06d596a21d9a7d640616ee2f990dc8f82", "amount_msat": "1000000", "delay": "21"}
		hops := []map[string]interface{}{{"pubkey": "039c49ccbf7341e1829152d0067b1b47b06d596a21d9a7d640616ee2f990dc8f82", "payload": "00000067000003000100000000000003e800000075000000000000000000000000000000000000000000000000"}}

		hexPaymentHash := hex.EncodeToString(addHTLC.PaymentHash[:])
		//hexOnionBlob := hex.EncodeToString(addHTLC.OnionBlob[:])

		onionBlob, err := p.Client.Call("createonion", hops, hexPaymentHash)
		if err != nil {
			p.Log("couldn't create onion: ", err)
		}

		p.Log("OnionBlob: \n", onionBlob.Get("onion").String(), "\n")
		p.Log("ChanID: \n", addHTLC.ChanID, "\n")
		p.Log("PaymentHash: \n", hexPaymentHash, "\n")

		// forward the onion to the next hop
		// NOTE: sendonion adds htlc to lightningd database so it can be retrieved with listsendpays
		result, err := p.Client.Call("sendonion", addHTLC.OnionBlob, firstHop, hexPaymentHash)
		if err != nil {
			p.Log("error sending onion: ", err)
		}
		p.Log("htlc status: ", result.Get("status").Str)

	case hcwire.MsgUpdateFulfillHTLC:
		//fulfillHTLC, ok := msg.(*hcwire.UpdateFulfillHTLC)
		_, ok := msg.(*hcwire.UpdateFulfillHTLC)
		if !ok {
			p.Log("unable to assert UpdateFulfillHTLC type")
			return continueHTLC
		}

	default:
		p.Log("handeling msg type: ", msg.MsgType(), " from peer: ", peer, " with content: ", payload[:])
	}
	return continueHTLC
}

func hcPay(p *plugin.Plugin, params plugin.Params) (interface{}, int, error) {
	/*
		- [] decode invoice
		- [] create route with getroute -> hard code for now: l1 --invoice--> l3 --hcpay--> l2 --sendonion--> l3
			- l1 id: 039c49ccbf7341e1829152d0067b1b47b06d596a21d9a7d640616ee2f990dc8f82
			- l2 id: 02ba62ac6b9819d140695c4676eaac7330574af9f08abbe3928d765a396ed915a9
			- l3 id: 03361b104205147575ccaf4c7f111dfd533b8db7cd590c1087051edbc450248467
		- [] createonion hops payment_hash(assoc_data)
		- [] create update_add_htlc for hc peer (with onion blob for next hop)
		- [] send to hc peer with sendcustommsg
		- [] wait for htlc_fulfill

		PEER VIEW:
		- [] receive htlc
		- [] make htlc with next hop
		- [] wait for htlc_fulfill
		- []
	*/

	invoice, err := p.Client.Call("decode", params.Get("bolt11").String())
	if err != nil {
		return nil, 1, err
	}
	//payee := invoice.Get("payee").String()
	paymentHash := invoice.Get("payment_hash").String()
	//amount := invoice.Get("amount_msat")
	//result, err := p.Client.Call("getroute", payee, amount, 10)
	//route := result.Get("route")
	// parse route so that payload of n is parameters of n+1

	// create (last) onion payload for l2 -> l1; see here
	//paymentSecret := invoice.Get("payment_secret").String()

	// hardcoded hop l2 -> l1 (pubkey)
	firstHop := map[string]interface{}{"pubkey": "039c49ccbf7341e1829152d0067b1b47b06d596a21d9a7d640616ee2f990dc8f82", "payload": "00000067000003000100000000000003e800000075000000000000000000000000000000000000000000000000"}
	hops := []map[string]interface{}{firstHop}

	onionBlob, err := p.Client.Call("createonion", hops, paymentHash)
	if err != nil {
		return nil, 1, err
	}
	tmp, _ := hex.DecodeString(onionBlob.Get("onion").String())
	var onionBlobBytes [1366]byte
	copy(onionBlobBytes[:], tmp)

	// channel id between l1-----l2
	//channelid, _ := hex.DecodeString("4f642a22b2205f87c6e4a939b73c7c45368c944979365b89ccc8a4470a6b6fc0")
	// l1----l3
	//channelid, _ := hex.DecodeString("f76489c99659f713890649f721dde11f569ea1ffeff80febe53656e54c67ea80"
	channelid, _ := hex.DecodeString("6e0000010001")

	var cid [32]byte
	copy(cid[:], channelid)

	tmp, _ = hex.DecodeString(paymentHash)
	var paymentHashBytes [32]byte
	copy(paymentHashBytes[:], tmp)

	addHTLC := hcwire.UpdateAddHTLC{
		lnwire.UpdateAddHTLC{
			ChanID:      cid,
			ID:          1,
			Amount:      lnwire.MilliSatoshi(1000011),
			PaymentHash: paymentHashBytes,
			Expiry:      144,
			OnionBlob:   onionBlobBytes,
		},
	}

	buf := new(bytes.Buffer)
	if _, err := hcwire.WriteMessage(buf, &addHTLC, 1); err != nil {
		return nil, 1, err
	}
	payload := hex.EncodeToString(buf.Bytes())

	p.Log("hc-pay sending message: \n", addHTLC)

	p.Client.Call("sendcustommsg", params.Get("node_id").String(), payload)

	// wait for payment preimage from hc peer

	return nil, 0, nil
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
		HOST VIEW:
		- [] parse HTLC message
		- [] return return HTLC to lightningd if short channel ID (scid) not hosted from hosted channels
		- [] forward HTLC to hosted channel peer
		- [] wait for htlc_fulfill from CLIENT
		- [] resolve HTLC with non-hosted-channel peer with payment key/preimage from htlc_fufill
	*/

	p.Log("HTLC:\n", params.Get("@pretty").String())

	return continueHTLC
}
