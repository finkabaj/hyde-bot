package mogs

import (
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
	"github.com/stretchr/testify/mock"
)

type MockReactionService struct {
	mock.Mock
}

func NewMockReactionService() *MockReactionService {
	return &MockReactionService{}
}

func (m *MockReactionService) CreateReactionRules(rules *[]rule.ReactionRule) (*[]rule.ReactionRule, error) {
	args := m.Called(rules)

	var arg0 *[]rule.ReactionRule
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*[]rule.ReactionRule)
	}

	return arg0, args.Error(1)
}

func (m *MockReactionService) GetReactionRules(gId string) (*[]rule.ReactionRule, error) {
	args := m.Called(gId)

	var arg0 *[]rule.ReactionRule
	if args.Get(0) != nil {
		arg0 = args.Get(0).(*[]rule.ReactionRule)
	}

	return arg0, args.Error(1)
}

func (m *MockReactionService) DeleteReactionRules(query *[]rule.DeleteReactionRuleQuery, gId string) error {
	args := m.Called(query, gId)
	return args.Error(0)
}
