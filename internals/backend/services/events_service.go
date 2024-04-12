package services

import (
	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
)

type IEventsService interface {
	CreateGuild(g *guild.GuildCreate) (*guild.Guild, error)
	GetGuild(gId string) (*guild.Guild, error)
}

type EventsService struct {
	database db.Database
}

var es *EventsService

func NewEventsService(d db.Database) *EventsService {
	if es == nil {
		es = &EventsService{
			database: d,
		}
	}
	return es
}

func (e *EventsService) CreateGuild(g *guild.GuildCreate) (*guild.Guild, error) {
	if g, err := e.GetGuild(g.GuildId); err != nil {
		return nil, err
	} else if g != nil {
		return nil, guild.ErrGuildConflict
	}

	newGuild, err := e.database.CreateGuild(g)

	return newGuild, err
}

func (e *EventsService) GetGuild(gId string) (*guild.Guild, error) {
	guild, err := e.database.GetGuild(gId)

	return guild, err
}
