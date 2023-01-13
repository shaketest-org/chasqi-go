package loop

import "chasqi-go/types"

type StatusGetter interface {
	ById(id string) *types.LoopStatus
}

type TreeEnqueuer interface {
	Enqueue(tree *types.Tree) error
}

type Canceler interface {
	Cancel(id string) error
}

type ResultGetter interface {
	Get(id string) *types.TestResult
}

type Engine interface {
	StatusGetter
	TreeEnqueuer
	ResultGetter
	Canceler
}
