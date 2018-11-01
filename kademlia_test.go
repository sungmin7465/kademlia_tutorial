package kademlia

import "testing"

func TestPing(t *testing.T) {
  me := Contact{NewRandomNodeID(), "127.0.0.1:8989"};
  k := NewKademlia(&me, "test");
  k.Serve();

  someone := Contact{NewRandomNodeID(), "127.0.0.1:8989"};
  if err := k.Call(
      &someone,
      "KademliaCore.Ping",
      &PingRequest{RPCHeader{&someone, k.NetworkId}},
      &PingResponse{}); err != nil {
    t.Error(err);
  }
}

func TestFindNode(t *testing.T) {
  me := Contact{NewRandomNodeID(), "127.0.0.1:8989"};
  k := NewKademlia(&me, "test");
  kc := KademliaCore{k};
  
  var contacts [100]Contact;
  for i := 0; i < len(contacts); i++ {
    contacts[i] = Contact{NewRandomNodeID(), "127.0.0.1:8989"};
    if err := kc.Ping(&PingRequest{RPCHeader{&contacts[i], k.NetworkId}},
                      &PingResponse{}); err != nil {
      t.Error(err);
    }
  }

  args := FindNodeRequest{RPCHeader{&contacts[0], k.NetworkId}, contacts[0].id};
  response := FindNodeResponse{};
  if err := kc.FindNode(&args, &response); err != nil {
    t.Error(err);
  }
  
  if len(response.contacts) != BucketSize {
    t.Fail();
  }
}
