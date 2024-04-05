package db

import (
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
)

type Database interface {
	Connect(credentials *DatabaseCredentials) error
	Close()
	Status() error

	//* GUILDS *//

	CreateGuild(guild *guild.GuildCreate) (*guild.Guild, error)
	GetGuild(guildId string) (*guild.Guild, error)
}

type DatabaseCredentials struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
