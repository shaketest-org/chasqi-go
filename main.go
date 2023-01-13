package main

import (
	"chasqi-go/core/agent"
	"chasqi-go/core/engine"
	"chasqi-go/data"
	"chasqi-go/types"
	_ "embed"
	"encoding/json"
	"log"
	"sync"
)

//go:embed core/agent/testdata/tree.json
var fixture []byte

func main() {
	wg := &sync.WaitGroup{}
	tree := &types.Tree{}
	json.Unmarshal(fixture, tree)
	log.Printf("tree: %v", tree)

	e := engine.New(
		func() agent.NodeVisitor { return data.NewDefaultHttpClient() },
		make(chan struct{}),
	)
	wg.Add(1)
	go startEngine(e)
	go func() {
		e.Enqueue(tree)
	}()
	wg.Wait()
}

func startEngine(e *engine.DefaultEngine) {
	e.Start()
}
