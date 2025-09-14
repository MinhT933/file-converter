package nats

import (
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func Connect(url string) *nats.Conn {
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	log.Println("✅ Connected to nats ", nc.IsConnected())
	return nc
}

func NewJetStream(nc *nats.Conn) jetstream.JetStream {
	js, _ := jetstream.New(nc)
	return js
}
