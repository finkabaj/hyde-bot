package mogs

import (
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
	"github.com/stretchr/testify/mock"
)

type MockGuildService struct {
	mock.Mock
}

func NewMockGuildService() *MockGuildService {
	return &MockGuildService{}
}

func (m *MockGuildService) CreateGuild(g guild.GuildCreate) (guild.Guild, error) {
	args := m.Called(g)

	return args.Get(0).(guild.Guild), args.Error(1)
}

func (m *MockGuildService) GetGuild(gId string) (guild.Guild, error) {
	args := m.Called(gId)

	return args.Get(0).(guild.Guild), args.Error(1)
}
