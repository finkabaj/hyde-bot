package mogs

import (
	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/ranks"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
	"github.com/stretchr/testify/mock"
)

type DbMock struct {
	mock.Mock
}

func NewDbMock() *DbMock {
	return &DbMock{}
}

func (m *DbMock) Connect(credentials db.DatabaseCredentials) error {
	args := m.Called(credentials)
	return args.Error(0)
}

func (m *DbMock) Close() {
	m.Called()
}

func (m *DbMock) Status() error {
	args := m.Called()
	return args.Error(0)
}

func (m *DbMock) CreateGuild(g guild.GuildCreate) (guild.Guild, error) {
	args := m.Called(g)
	return args.Get(0).(guild.Guild), args.Error(1)
}

func (m *DbMock) ReadGuild(guildId string) (guild.Guild, error) {
	args := m.Called(guildId)
	return args.Get(0).(guild.Guild), args.Error(1)
}

func (m *DbMock) CreateReactionRules(rules []rule.ReactionRule) ([]rule.ReactionRule, error) {
	args := m.Called(rules)
	return args.Get(0).([]rule.ReactionRule), args.Error(1)
}

func (m *DbMock) DeleteReactionRules(rules []rule.DeleteReactionRuleQuery, gId string) error {
	args := m.Called(rules, gId)
	return args.Error(0)
}

func (m *DbMock) ReadReactionRules(gId string) ([]rule.ReactionRule, error) {
	args := m.Called(gId)
	return args.Get(0).([]rule.ReactionRule), args.Error(1)
}

func (m *DbMock) CreateRanks(r []ranks.Rank) ([]ranks.Rank, error) {
	args := m.Called(r)
	return args.Get(0).([]ranks.Rank), args.Error(1)
}

func (m *DbMock) ReadRanks(gId string) ([]ranks.Rank, error) {
	args := m.Called(gId)
	return args.Get(0).([]ranks.Rank), args.Error(1)
}

func (m *DbMock) DeleteRank(gId string, rId string) error {
	args := m.Called(gId, rId)
	return args.Error(0)
}

func (m *DbMock) DeleteRanks(gId string) error {
	args := m.Called(gId)
	return args.Error(0)
}
