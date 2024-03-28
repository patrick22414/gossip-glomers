package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"slices"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastBody struct {
	Message int `json:"message"`
	maelstrom.MessageBody
}

var BroadcastOK = maelstrom.MessageBody{Type: "broadcast_ok"}

type TopologyBody struct {
	Topology map[string][]string `json:"topology"`
	maelstrom.MessageBody
}

var TopologyOK = maelstrom.MessageBody{Type: "topology_ok"}

type ReadBody maelstrom.MessageBody

type ReadResponse struct {
	Messages []int `json:"messages"`
	maelstrom.MessageBody
}

func main() {
	n := maelstrom.NewNode()
	s := NewSet()
	var ns []string // neighbouring nodes

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body BroadcastBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		v := body.Message
		if !s.Has(v) {
			s.Add(v)
			body := &PropagateBody{
				Message: v,
				Already: []string{n.ID()},
				MessageBody: maelstrom.MessageBody{
					Type: "propagate",
				},
			}
			go propagate(n, ns, body)
		}

		return n.Reply(msg, BroadcastOK)
	})

	n.Handle("propagate", func(msg maelstrom.Message) error {
		var body PropagateBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		v := body.Message
		if !s.Has(v) {
			s.Add(v)
			body.Already = append(body.Already, n.ID())
			go propagate(n, ns, &body)
		}

		if body.Respond {
			return n.Reply(msg, PropagateOK)
		}
		return nil
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var body ReadBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		resp := ReadResponse{
			Messages: s.List(),
			MessageBody: maelstrom.MessageBody{
				Type: "read_ok",
			},
		}

		return n.Reply(msg, resp)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body TopologyBody
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		ns = body.Topology[n.ID()]

		// Add some random nodes as neighbours
		randoms := takeRand(n.NodeIDs(), 5)
		for _, nid := range randoms {
			if !slices.Contains(ns, nid) {
				ns = append(ns, nid)
			}
		}

		return n.Reply(msg, TopologyOK)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func takeRand[S ~[]E, E any](s S, n int) S {
	if n >= len(s) {
		return s
	}

	perm := rand.Perm(len(s))
	perm = perm[:n]
	s2 := make(S, 0, n)
	for p := range perm {
		s2 = append(s2, s[p])
	}
	return s2
}
