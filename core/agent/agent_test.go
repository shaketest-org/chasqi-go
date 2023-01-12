package agent

import (
	"bytes"
	mock_agent "chasqi-go/core/agent/mocks"
	"chasqi-go/types"
	_ "embed"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"sync"
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

func (s *AgentTestSuite) TestAgent_StartShouldCallExpectedRoutes() {
	tree := testTree()
	resultCh := make(chan types.TestResult, 0)
	subject := New(0, tree, resultCh, s.NodeVisitorMock)
	postBody := new(bytes.Buffer)
	putBody := new(bytes.Buffer)
	successfulResponse := &types.ResponseResult{
		StatusCode: 200,
	}
	var result *types.TestResult
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case r := <-resultCh:
				result = &r
				wg.Done()
			}
		}
	}()
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
	).Return(successfulResponse, nil).Times(1)
	s.NodeVisitorMock.EXPECT().Visit(
		"POST",
		"/api/users",
		postBody,
		map[string][]string{
			"Content-Type": {"application/json"},
		}).Return(successfulResponse, nil).Times(1)
	s.NodeVisitorMock.EXPECT().Visit(
		"GET",
		"/api/users/2",
		nil,
		map[string][]string{
			"Content-Type": {"application/json"},
		},
	).Return(successfulResponse, nil).Times(1)
	s.NodeVisitorMock.EXPECT().Visit(
		"PUT",
		"/api/users/2",
		putBody,
		map[string][]string{
			"Content-Type": {"application/json"},
		}).Return(successfulResponse, nil).Times(1)

	subject.Start()
	wg.Wait()
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), 4, result.SuccessCount)
	assert.Equal(s.T(), 0, result.ErrorCount)
	assert.Equal(s.T(), tree.ID, result.TreeID)
}

func testTree() *types.Tree {
	tree := &types.Tree{}
	json.Unmarshal(fixture, tree)
	return tree
}

func TestAgentTestSuite(t *testing.T) {
	suite.Run(t, new(AgentTestSuite))
}
