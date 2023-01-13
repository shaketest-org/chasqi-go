package result

import (
	"chasqi-go/types"
	"fmt"
	"log"
	"sync"
	"time"
)

type Manager struct {
	treeConfigMap  map[types.TreeID]*types.Config
	resultMap      map[types.TreeID]*types.TestResult
	agentResultMap map[types.TreeID][]*types.AgentResult
	mu             *sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		treeConfigMap:  make(map[types.TreeID]*types.Config),
		resultMap:      make(map[types.TreeID]*types.TestResult),
		agentResultMap: make(map[types.TreeID][]*types.AgentResult),
		mu:             &sync.Mutex{},
	}
}

func (r *Manager) Get(treeID types.TreeID) *types.TestResult {
	return r.resultMap[treeID]
}

func (r *Manager) Insert(result *types.AgentResult) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	treeID := types.TreeID(result.TreeID)
	if r.treeConfigMap[treeID] == nil {
		return fmt.Errorf("config for tree with id %s does not exist", result.TreeID)
	}

	arMap := r.agentResultMap[treeID]
	arMap = append(arMap, result)
	r.agentResultMap[treeID] = arMap

	if r.checkIfAllAgentsHaveFinished(treeID) {
		log.Printf("all agents have finished for tree %s", treeID)
		errCount, okCount, duration := GetTotalErrorCount(arMap)
		r.resultMap[treeID] = &types.TestResult{
			TreeID:            result.TreeID,
			AgentResults:      arMap,
			Successful:        true,
			TotalErrorCount:   errCount,
			TotalSuccessCount: okCount,
			TotalDuration:     time.Duration(duration * float64(time.Second)),
		}
		delete(r.agentResultMap, treeID)
	}
	return nil
}

func (r *Manager) Clear(treeID types.TreeID) {
	delete(r.treeConfigMap, treeID)
	delete(r.resultMap, treeID)
	delete(r.agentResultMap, treeID)
}

func (r *Manager) Prepare(treeID types.TreeID, config types.Config) error {
	if r.treeConfigMap[treeID] != nil {
		return fmt.Errorf("tree with id %s already exists", treeID)
	}

	r.treeConfigMap[treeID] = &config
	return nil
}
func (r *Manager) checkIfAllAgentsHaveFinished(treeID types.TreeID) bool {
	config := r.treeConfigMap[treeID]
	agentResults := r.agentResultMap[treeID]

	return len(agentResults) == config.AgentAmount
}
