package agent

import (
	"bytes"
	mock_agent "chasqi-go/core/agent/mocks"
	"chasqi-go/types"
	_ "embed"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
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

	controller := gomock.NewController(s.T())
	s.NodeVisitorMock = mock_agent.NewMockNodeVisitor(controller)
	s.resultCh = make(chan types.TestResult)
}

func (s *AgentTestSuite) TestAgent_Start() {
	tree := testTree()
	subject := New(0, tree, nil, s.NodeVisitorMock)
	postBody := new(bytes.Buffer)
	putBody := new(bytes.Buffer)
	json.NewEncoder(postBody).Encode(
		map[string]interface{}{
			"username": "JaneDoe",
			"userId":   2,
		})
	json.NewEncoder(putBody).Encode(map[string]interface{}{
		"username": "JaneSmith",
		"userId":   2,
	})

	s.NodeVisitorMock.EXPECT().Visit(
		"GET",
		"/api/users",
		nil,
		map[string][]string{
			"Content-Type": {"application/json"},
		},
	).Times(1)
	s.NodeVisitorMock.EXPECT().Visit(
		"POST",
		"/api/users",
		postBody,
		map[string][]string{
			"Content-Type": {"application/json"},
		}).Times(1)
	s.NodeVisitorMock.EXPECT().Visit(
		"GET",
		"/api/users/2",
		nil,
		map[string][]string{
			"Content-Type": {"application/json"},
		},
	).Times(1)
	s.NodeVisitorMock.EXPECT().Visit(
		"PUT",
		"/api/users/2",
		putBody,
		map[string][]string{
			"Content-Type": {"application/json"},
		})

	subject.Start()

}

func testTree() *types.Tree {
	tree := &types.Tree{}
	json.Unmarshal(fixture, tree)
	return tree
}

func TestAgentTestSuite(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}
