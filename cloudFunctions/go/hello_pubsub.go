package helloworld

import (
	"context"
	"log"
)

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

//  consumes a Pub/Sub message.
func GoPubSub(ctx context.Context, m PubSubMessage) error {
	name := string(m.Data)
	if name == "" {
		name = "World"
	}
	log.Printf("Hello, %s!", name)
	return nil
}

// [END functions_helloworld_background]
// [END functions_helloworld_pubsub]
