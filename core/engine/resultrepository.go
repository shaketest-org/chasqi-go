package engine

import "chasqi-go/types"

type (
	ResultInserter interface {
		Insert(result *types.AgentResult) error
	}
	// TreePreparer Must be called before any tree is inserted otherwise the repository will have no configs
	// to use and compare to when a tree is inserted.
	TreePreparer interface {
		Prepare(treeID types.TreeID, config types.Config) error
	}

	TestResultGetter interface {
		Get(treeID types.TreeID) *types.TestResult
	}

	ResultClearer interface {
		Clear(treeID types.TreeID)
	}

	ResultRepository interface {
		ResultInserter
		TreePreparer
		TestResultGetter
		ResultClearer
	}
)
