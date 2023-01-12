package agent

import (
	mock_agent "chasqi-go/core/agent/mocks"
	"chasqi-go/types"
	_ "embed"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

//go:embed testdata/tree.json
var fixture []byte

type AgentTestSuite struct {
	suite.Suite
	NodeVisitorMock *mock_agent.MockNodeVisitor
	resultCh        chan types.TestResult
}

func (s *AgentTestSuite) SetupTest() {
	if fixture == nil {
		s.T().Error("fixture is null")
	}
	log.Println("fixture is not null")
	controller := gomock.NewController(s.T())
	s.NodeVisitorMock = mock_agent.NewMockNodeVisitor(controller)
	s.resultCh = make(chan types.TestResult)
}

func (s *AgentTestSuite) TestAgent_Start() {
	tree := testTree()
	subject := New(0, tree, nil, s.NodeVisitorMock)
	subject.Start()

	s.NodeVisitorMock.EXPECT().Get(
		"/api/users",
		map[string]interface{}{
			"Content-Type": "application/json",
		},
	).Times(1)
	s.NodeVisitorMock.EXPECT().Post(
		"/api/users",
		map[string]interface{}{
			"Content-Type": "application/json",
		},
		map[string]interface{}{
			"username": "JaneDoe",
			"userId":   2,
		}).Times(1)
	s.NodeVisitorMock.EXPECT().Put(
		"/api/users/2",
		map[string]interface{}{
			"Content-Type": "application/json",
		},
		map[string]interface{}{
			"username": "JaneDoe",
			"userId":   2,
		})
}

func testTree() *types.Tree {
	tree := &types.Tree{}
	json.Unmarshal(fixture, tree)
	return tree
}

func TestAgentTestSuite(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}
