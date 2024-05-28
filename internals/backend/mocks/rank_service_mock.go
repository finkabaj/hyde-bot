package mogs

import (
	"github.com/finkabaj/hyde-bot/internals/ranks"
	"github.com/stretchr/testify/mock"
)

type MockRankService struct {
	mock.Mock
}

func (m *MockRankService) GetRanks(gId string) (ranks.Ranks, error) {
	args := m.Called(gId)
	return args.Get(0).(ranks.Ranks), args.Error(1)
}

func (m *MockRankService) CreateRanks(r ranks.Ranks) (ranks.Ranks, error) {
	args := m.Called(r)
	return args.Get(0).(ranks.Ranks), args.Error(1)
}

func (m *MockRankService) DeleteRank(gId string, rId string) error {
	args := m.Called(gId, rId)
	return args.Error(0)
}

func (m *MockRankService) DeleteRanks(gId string) error {
	args := m.Called(gId)
	return args.Error(0)
}
