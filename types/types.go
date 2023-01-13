package types

import (
	"fmt"
	"time"
)

type (
	Node struct {
		Method  string              `json:"method"`
		Path    string              `json:"path"`
		Body    map[string]any      `json:"body"`
		Headers map[string][]string `json:"headers"`
		Next    *Edge               `json:"next"`
		Before  *Edge               `json:"before"`
	}

	Edge struct {
		Delay float32 `json:"delay"`
		// If after is null then it's the end of the journey
		After *Node `json:"after"`
		// If before is null then it's the root
		Before *Node `json:"before"`
	}

	Config struct {
		AgentAmount    int     `json:"agentAmount"`
		MaxFailureRate float32 `json:"maxFailureRate"`
	}

	TreeID string

	Tree struct {
		Config *Config `json:"config"`
		Root   *Node   `json:"root"`
		ID     string  `json:"id"`
	}

	LoopStatus struct {
		IsDone     bool       `json:"isDone"`
		TreeID     string     `json:"treeID"`
		StartedAt  *time.Time `json:"startedAt"`
		FinishedAt *time.Time `json:"finishedAt"`
	}

	ResponseResult struct {
		StatusCode int                 `json:"statusCode"`
		Method     string              `json:"method"`
		Url        string              `json:"url"`
		Body       map[string]any      `json:"body"`
		Headers    map[string][]string `json:"headers"`
		Duration   time.Duration       `json:"durationInMs"`
		StartedAt  time.Time           `json:"startedAt"`
		FinishedAt time.Time           `json:"finishedAt"`
	}

	AgentResult struct {
		ErrorCount   int               `json:"errorCount"`
		SuccessCount int               `json:"successCount"`
		AgentID      int               `json:"agentID"`
		TreeID       string            `json:"treeID"`
		Result       []*ResponseResult `json:"result"`
		FinishedAt   time.Time         `json:"finishedAt"`
	}

	TestResult struct {
		Successful        bool           `json:"successful"`
		TotalErrorCount   int            `json:"totalErrorCount"`
		TotalSuccessCount int            `json:"totalSuccessCount"`
		TreeID            string         `json:"treeID"`
		AgentResults      []*AgentResult `json:"agentResults"`
		TotalDuration     time.Duration  `json:"totalDuration"`
	}
)

func (n *Node) String() string {
	return n.Method + " " + n.Path
}

func (tr *AgentResult) String() string {
	return fmt.Sprintf("AgentResult: { ErrorCount: %d, SuccessCount: %d, TreeID: %s, Resulst: %d}", tr.ErrorCount, tr.SuccessCount, tr.TreeID, len(tr.Result))
}

func (tr *AgentResult) Duration() time.Duration {
	d := time.Duration(0)
	for _, r := range tr.Result {
		d += r.Duration
	}
	return d
}
