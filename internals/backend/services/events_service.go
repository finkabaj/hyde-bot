package services

import (
	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
)

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

func (e *EventsService) CreateGuild(guild *guild.GuildCreate) (*guild.Guild, error) {
	if guild, err := e.database.GetGuild(guild.GuildId); guild != nil {
		return nil, err
	}

	return nil, nil
}

func (e *EventsService) GetGuild(gId string) (*guild.Guild, error) {
	guild, err := e.database.GetGuild(gId)

	return guild, err
}
