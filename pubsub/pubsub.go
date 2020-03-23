package pubsub

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

type G struct {
	sync.Mutex
	CredentialFile string
	Credential     Credential
	opt            option.ClientOption
	bucketHandle   *storage.BucketHandle
	err            error
}

type Credential struct {
	Type                        string `json:"type"`
	Project_id                  string `json:"project_id"`
	Private_key_id              string `json:"private_key_id"`
	Private_key                 string `json:"private_key"`
	Client_email                string `json:"client_email"`
	Client_id                   string `json:"client_id"`
	Auth_uri                    string `json:"auth_uri"`
	Token_uri                   string `json:"token_uri"`
	Auth_provider_x509_cert_url string `json:"auth_provider_x509_cert_url"`
	Client_x509_cert_url        string `json:"client_x509_cert_url"`
}

func FindFile() ([]byte, string) {

	directories := []string{"../credentials",
		"/credentials",
		"/etc/credentials", "./credentials", "../../credentials"}
	file := "pubsub.json"

	for _, v := range directories {
		path := fmt.Sprintf("%s/%s", v, file)
		data, err := ioutil.ReadFile(path)
		if err == nil {
			return data, path
		}

	}
	errors.New("Cannot find pubsub.json")
	return nil, ""
}

func NewG() *G {

	data, file := FindFile()
	g := G{CredentialFile: file}
	credential := Credential{}
	json.Unmarshal([]byte(data), &credential)
	g.Credential = credential
	g.opt = option.WithCredentialsFile(g.CredentialFile)
	return &g
}

func (g *G) CreateTopic(topicName string) (*pubsub.Topic, error) {
	g.Lock()
	defer g.Unlock()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, g.Credential.Project_id, g.opt)

	var topic *pubsub.Topic
	topic = client.Topic(topicName)
	// Create the topic if it doesn't exist.
	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Printf("topic: %v\n", err)
		return nil, err
	}
	if !exists {
		log.Printf("Topic %v doesn't exist - creating it", topicName)
		_, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return topic, nil
}

func (g *G) CreateSub(subName string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	g.Lock()
	defer g.Unlock()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, g.Credential.Project_id, g.opt)

	var subscription *pubsub.Subscription
	config := pubsub.SubscriptionConfig{
		Topic:               topic,
		PushConfig:          pubsub.PushConfig{},
		AckDeadline:         30000000000, // 30 seconds
		RetainAckedMessages: false,
		RetentionDuration:   0,
		ExpirationPolicy:    nil,
		Labels:              nil,
		DeadLetterPolicy:    nil,
	}
	subscription, err = client.CreateSubscription(ctx, subName, config)

	if err != nil {
		log.Printf("topic: %v\n", err)
		return subscription, err
	}

	return subscription, nil
}

func (g *G) Publish(w io.Writer, topicID, msg string) (string, error) {
	g.Lock()
	defer g.Unlock()
	// projectID := "my-project-id"
	// topicID := "my-topic"
	// msg := "Hello World"
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, g.Credential.Project_id, g.opt)
	if err != nil {
		return "", fmt.Errorf("pubsub.NewClient: %v", err)
	}

	t := client.Topic(topicID)
	defer t.Stop()
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		return "", fmt.Errorf("Get: %v", err)
	}
	fmt.Fprintf(w, "Published a message; msg ID: %v\n", id)

	return id, nil
}

func (g *G) PullMsgs(w io.Writer, subID string) ([]byte, error) {
	g.Lock()
	defer g.Unlock()

	var message []byte
	// projectID := "my-project-id"
	// subID := "my-sub"
	// topic of type https://godoc.org/cloud.google.com/go/pubsub#Topic
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, g.Credential.Project_id, g.opt)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}

	// Consume 1 messages.
	var mu sync.Mutex
	received := 0
	sub := client.Subscription(subID)
	cctx, cancel := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Fprintf(w, "Got message: %q\n", string(msg.Data))
		msg.Ack()
		mu.Lock()
		defer mu.Unlock()
		received++
		message = msg.Data
		if received == 1 {
			cancel()
		}
	})
	if err != nil {
		return nil, fmt.Errorf("Receive: %v", err)
	}
	return message, nil
}
