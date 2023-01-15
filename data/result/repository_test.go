package result

import (
	"chasqi-go/types"
	"testing"
	"time"
)

func TestRepository_Insert(t *testing.T) {
	treeID := "tree-id"
	t.Run("should return an error if Prepare was not called", func(t *testing.T) {
		subject := NewManager()

		err := subject.Insert(&types.AgentResult{
			TreeID: treeID,
		})
		if err == nil {
			t.Error("expected error but got nil")
		}
		result := subject.Get(types.TreeID(treeID))
		if result != nil {
			t.Errorf("expected nil but got %v", result)
		}
	})
	t.Run("should return not an error if Prepare was not called", func(t *testing.T) {
		subject := NewManager()

		subject.Prepare(types.TreeID(treeID), types.Config{})
		err := subject.Insert(&types.AgentResult{
			TreeID: treeID,
		})
		if err != nil {
			t.Error("unexpected error ", err)
		}
		result := subject.Get(types.TreeID(treeID))
		if result != nil {
			t.Error("expected no result but got one")
		}
	})
	t.Run("should return a result if max agentCount equals to agentResults", func(t *testing.T) {
		subject := NewManager()
		f := time.Now()

		subject.Prepare(types.TreeID(treeID),
			types.Config{
				AgentAmount: 1,
			})
		err := subject.Insert(&types.AgentResult{
			TreeID:     treeID,
			FinishedAt: f,
		})
		if err != nil {
			t.Error("unexpected error ", err)
		}
		result := subject.Get(types.TreeID(treeID))

		if result == nil {
			t.Error("expected a result but got none")
		}
	})
}
