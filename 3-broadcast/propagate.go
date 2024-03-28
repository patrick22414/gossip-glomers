package main

import (
	"context"
	"slices"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type PropagateBody struct {
	Message int      `json:"message"`
	Already []string `json:"already"`
	Respond bool     `json:"respond"`
	maelstrom.MessageBody
}

var PropagateOK = maelstrom.MessageBody{Type: "propagate_ok"}

func propagate(n *maelstrom.Node, ns []string, body *PropagateBody) {
	var wg sync.WaitGroup
	for _, dest := range ns {
		if slices.Contains(body.Already, dest) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			propagateNoResponse(n, dest, body)
			// propagateWithResponseAndRetry(n, dest, body)
		}()
	}
	wg.Wait()
}

// propagateNoResponse is quicker, but not fault-tolerant
func propagateNoResponse(n *maelstrom.Node, dest string, body *PropagateBody) {
	body.Respond = false
	n.Send(dest, body)
}

// propagateWithResponseAndRetry is more fault-tolerant, but slower
func propagateWithResponseAndRetry(n *maelstrom.Node, dest string, body *PropagateBody) {
	body.Respond = true
	for try := 0; try < 3; try++ {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		msg, err := n.SyncRPC(ctx, dest, body)
		if err != nil || msg.Type() != PropagateOK.Type {
			cancel()
			continue
		}

		cancel()
		break
	}
}
