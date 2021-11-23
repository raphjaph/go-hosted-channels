package main

import (
	"log"

	lightning "github.com/fiatjaf/lightningd-gjson-rpc"
)

var ln *lightning.Client

func main() {

	ln = &lightning.Client{
		SparkURL:              "http://localhost:9737",
		SparkToken:            "masterkeythatcandoeverything",
		DontCheckCertificates: true,
	}

	peerInfo, err := ln.Call("listpeers")
	if err != nil {
		log.Fatal("listpeers error: " + err.Error())
	}

	log.Print(peerInfo.Get("peers.0.id"))

	myMethod, err := ln.Call("sendcustommsg", peerInfo.Get("peers.0.id").Str, "1337ffffffff")
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Print(myMethod.String())
}
