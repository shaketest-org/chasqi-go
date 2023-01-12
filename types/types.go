package types

import "time"

type (
	Node struct {
		Method  string          `json:"method"`
		Path    string          `json:"path"`
		Body    *map[string]any `json:"body"`
		Headers *map[string]any `json:"header"`
		Next    *Edge           `json:"next"`
		Before  *Edge           `json:"before"`
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

	ResultPair struct {
		Node       *Node      `json:"node"`
		FinishedAt *time.Time `json:"finishedAt"`
	}

	TestResult struct {
		Success      bool          `json:"success"`
		ErrorCount   int           `json:"errorCount"`
		SuccessCount int           `json:"successCount"`
		TreeID       string        `json:"treeID"`
		Result       []*ResultPair `json:"result"`
		StartedAt    *time.Time    `json:"startedAt"`
		FinishedAt   *time.Time    `json:"finishedAt"`
	}

	ResponseResult struct {
		StatusCode   int   `json:"statusCode"`
		DurationInMs int64 `json:"durationInMs"`
	}
)

func (n *Node) String() string {
	return n.Method + " " + n.Path
}
