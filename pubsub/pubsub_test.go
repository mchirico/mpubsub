package pubsub

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_creds(t *testing.T) {
	g := NewG()
	if len(g.Credential.Project_id) > 3 {
		t.Logf("We have project: %s\n", g.Credential.Project_id)
	} else {
		t.Fatalf("Can't read credential file: ../credentials/credential.json")
	}
}

func TestFindFile(t *testing.T) {
	_, s := FindFile()
	if strings.Contains(s, ".json") {
		t.Logf("found: %s\n", s)
	} else {
		t.Fatalf("Cannot find .json")
	}
}

func TestG_Publish(t *testing.T) {
	g := NewG()
	var buf bytes.Buffer
	id, err := g.Publish(&buf, "test", "test")
	if err != nil {
		t.Fatalf("error: %v\n", err)
	}
	fmt.Printf("id: %v\n", id)
}

func TestG_CreateTopic(t *testing.T) {
	g := NewG()
	_, err := g.CreateTopic("test")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}

func TestG_CreateSubForCloudFunctions(t *testing.T) {
	g := NewG()
	topic, _ := g.CreateTopic("gocloud")
	_, err := g.CreateSub("sub-gocloud", topic)
	if err != nil {
		if strings.Contains(err.Error(), "code = AlreadyExists desc = Resource ") {
			t.Logf("This is okay... it should exist")
		} else {
			t.Fatal("Sub error")
		}
	}
	var buf bytes.Buffer
	g.Publish(&buf, "gocloud", "test")
	topic.Stop()
}


func TestG_CreateSub(t *testing.T) {
	g := NewG()
	topic, _ := g.CreateTopic("test")
	_, err := g.CreateSub("sub-test", topic)
	if err != nil {
		if strings.Contains(err.Error(), "code = AlreadyExists desc = Resource ") {
			t.Logf("This is okay... it should exist")
		} else {
			t.Fatal("Sub error")
		}
	}
	topic.Stop()
}

func CreateMsg() {
	g := NewG()
	var buf bytes.Buffer
	id, err := g.Publish(&buf, "test", "test")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	fmt.Printf("id: %v\n", id)
}

func TestG_PullMsgs(t *testing.T) {

	CreateMsg()
	g := NewG()
	var buf bytes.Buffer
	msg, err := g.PullMsgs(&buf, "sub-test")
	if err != nil {
		t.Fatalf("No message")
	}
	fmt.Printf("msg: %s\n", msg)

}
