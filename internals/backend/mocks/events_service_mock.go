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
	ret := m.Called(g)

	var r0 *guild.Guild
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*guild.Guild)
	}
	r1 := ret.Error(1)

	return r0, r1
}

func (m *MockEventsService) GetGuild(gId string) (*guild.Guild, error) {
	args := m.Called(gId)

	if args.Get(0) == nil {
		return nil, nil
	}

	return args.Get(0).(*guild.Guild), args.Error(1)
}
