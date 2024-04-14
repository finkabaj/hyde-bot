package mogs

import (
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
	"github.com/stretchr/testify/mock"
)

type MockEventsService struct {
	mock.Mock
}

func NewMockEventsService() *MockEventsService {
	return &MockEventsService{}
}

func (m *MockEventsService) CreateGuild(g *guild.GuildCreate) (*guild.Guild, error) {
	args := m.Called(g)

	var a0 *guild.Guild

	if args.Get(0) != nil {
		a0 = args.Get(0).(*guild.Guild)
	}

	return a0, args.Error(1)
}

func (m *MockEventsService) GetGuild(gId string) (*guild.Guild, error) {
	args := m.Called(gId)

	var a0 *guild.Guild
	if args.Get(0) != nil {
		a0 = args.Get(0).(*guild.Guild)
	}

	return a0, args.Error(1)
}
