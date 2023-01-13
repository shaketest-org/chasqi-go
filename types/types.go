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
		TreeID    string     `json:"treeID"`
		IsDone    bool       `json:"isDone"`
		StartedAt *time.Time `json:"startedAt"`
	}

	TestResult struct {
		ErrorCount   int               `json:"errorCount"`
		SuccessCount int               `json:"successCount"`
		TreeID       string            `json:"treeID"`
		Result       []*ResponseResult `json:"result"`
		StartedAt    *time.Time        `json:"startedAt"`
		FinishedAt   *time.Time        `json:"finishedAt"`
	}

	ResponseResult struct {
		StatusCode   int                 `json:"statusCode"`
		Method       string              `json:"method"`
		Url          string              `json:"url"`
		Body         map[string]any      `json:"body"`
		Headers      map[string][]string `json:"headers"`
		DurationInMs int64               `json:"durationInMs"`
		StartedAt    *time.Time          `json:"startedAt"`
		FinishedAt   *time.Time          `json:"finishedAt"`
	}
)

func (n *Node) String() string {
	return n.Method + " " + n.Path
}

func (tr *TestResult) String() string {
	return fmt.Sprintf("TestResult: { ErrorCount: %d, SuccessCount: %d, TreeID: %s, Resulst: %d, StartedAt: %v, FinishedAt: %v }", tr.ErrorCount, tr.SuccessCount, tr.TreeID, len(tr.Result), tr.StartedAt, tr.FinishedAt)
}
