package services

import (
	"reflect"

	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
)

type IGuildService interface {
	CreateGuild(g guild.GuildCreate) (guild.Guild, error)
	GetGuild(gId string) (guild.Guild, error)
}

type GuildService struct {
	database db.Database
}

var es *GuildService

func NewGuildService(d db.Database) *GuildService {
	if es == nil {
		es = &GuildService{
			database: d,
		}
	}
	return es
}

func (e *GuildService) CreateGuild(g guild.GuildCreate) (guild.Guild, error) {
	if g, err := e.GetGuild(g.GuildId); err != nil {
		return guild.Guild{}, err
	} else if !reflect.DeepEqual(g, guild.Guild{}) {
		return guild.Guild{}, guild.ErrGuildConflict
	}

	newGuild, err := e.database.CreateGuild(g)

	return newGuild, err
}

func (e *GuildService) GetGuild(gId string) (guild.Guild, error) {
	guild, err := e.database.ReadGuild(gId)

	return guild, err
}
