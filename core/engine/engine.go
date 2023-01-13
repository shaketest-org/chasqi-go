package engine

import (
	"chasqi-go/core/agent"
	"chasqi-go/types"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	loopInterval  = 1 * time.Second
	flushInterval = 5 * 60 * time.Second
	maxResultAge  = 5 * 60 * time.Second
)

type DefaultEngine struct {
	hasStopped       bool
	visitorCreator   func() agent.NodeVisitor
	resultRepository ResultRepository
	enqueuedTrees    []*types.Tree
	statusMap        map[types.TreeID]*types.LoopStatus
	activeTrees      map[types.TreeID]*types.Tree
	resultMap        map[types.TreeID]*types.AgentResult
	resultCh         chan types.AgentResult
	exitCh           chan struct{}
	mu               *sync.Mutex
}

// Start begins the engine loop.
// Every second it will check if there are any trees to be processed.
// Every 5 minutes it will flush the result map.
func (e *DefaultEngine) Start() {
	log.Printf("Test Engine started")
	loopTimer := time.NewTicker(loopInterval)
	flushTimer := time.NewTicker(flushInterval)
coreLoop:
	for {
		select {
		case <-loopTimer.C:
			e.onTick()
		case <-flushTimer.C:
			e.onFlush()
		case result := <-e.resultCh:
			e.onResult(&result)
		case <-e.exitCh:
			break coreLoop
		}
	}
}

func (e *DefaultEngine) onTick() {
	e.mu.Lock()
	for _, t := range e.enqueuedTrees {
		log.Printf("processing tree: %v", t)
		newEnqueuedTrees := make([]*types.Tree, 0)

		for _, t2 := range e.enqueuedTrees {
			if t2 != t {
				newEnqueuedTrees = append(newEnqueuedTrees, t2)
			}
		}

		s := time.Now()
		e.statusMap[types.TreeID(t.ID)] = &types.LoopStatus{
			TreeID:    t.ID,
			IsDone:    false,
			StartedAt: &s,
		}
		e.visitTree(t)
		e.enqueuedTrees = newEnqueuedTrees
	}
	e.mu.Unlock()
}

func (e *DefaultEngine) onFlush() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for k, v := range e.resultMap {
		t := time.Now()
		if t.Sub(*v.FinishedAt) > maxResultAge {
			e.resultRepository.Clear(k)
			delete(e.resultMap, k)
			delete(e.statusMap, k)
		}
	}
}

func (e *DefaultEngine) onResult(result *types.AgentResult) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.activeTrees, types.TreeID(result.TreeID))
	now := time.Now()
	e.statusMap[types.TreeID(result.TreeID)].IsDone = true
	e.statusMap[types.TreeID(result.TreeID)].FinishedAt = &now
	e.resultMap[types.TreeID(result.TreeID)] = result

	log.Printf("Finished tree: %s", result.String())
	log.Printf("Duration was: %s", result.Duration().String())
	e.resultRepository.Insert(result)
}

func (e *DefaultEngine) ById(id string) *types.LoopStatus {
	return e.statusMap[types.TreeID(id)]
}

func (e *DefaultEngine) Enqueue(tree *types.Tree) error {
	if e.hasStopped {
		return fmt.Errorf("DefaultEngine has stopped")
	}
	log.Printf("Enqueuing tree: %s", tree.ID)
	e.mu.Lock()
	defer e.mu.Unlock()
	e.enqueuedTrees = append(e.enqueuedTrees, tree)
	return nil
}

func (e *DefaultEngine) visitTree(tree *types.Tree) {
	n := tree.Config.AgentAmount
	for i := 0; i < n; i++ {
		go func(i int) {
			a := agent.New(
				i,
				tree,
				e.resultCh,
				e.visitorCreator())
			a.Start()
		}(i)
	}
}

func (e *DefaultEngine) Get(id string) *types.TestResult {
	return e.resultRepository.Get(types.TreeID(id))
}

func (e *DefaultEngine) Cancel(id string) error {
	//TODO implement me
	panic("implement me")
}

func New(visitorCreator func() agent.NodeVisitor, repository ResultRepository, exitCh chan struct{}) *DefaultEngine {
	return &DefaultEngine{
		resultRepository: repository,
		statusMap:        make(map[types.TreeID]*types.LoopStatus),
		activeTrees:      make(map[types.TreeID]*types.Tree),
		resultMap:        make(map[types.TreeID]*types.AgentResult),
		enqueuedTrees:    make([]*types.Tree, 0),
		visitorCreator:   visitorCreator,
		resultCh:         make(chan types.AgentResult, 1000),
		exitCh:           exitCh,
		mu:               &sync.Mutex{},
	}
}
