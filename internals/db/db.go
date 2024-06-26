package db

import (
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
)

type Database interface {
	Connect(credentials DatabaseCredentials) error
	Close()
	Status() error

	//* GUILDS *//

	CreateGuild(guild guild.GuildCreate) (guild.Guild, error)
	ReadGuild(guildId string) (guild.Guild, error)

	// * RULES * //

	/// ** REACTIONS ** ///

	CreateReactionRules(rules []rule.ReactionRule) ([]rule.ReactionRule, error)
	DeleteReactionRules(rules []rule.DeleteReactionRuleQuery, gId string) error
	ReadReactionRules(gId string) ([]rule.ReactionRule, error)
}

type DatabaseCredentials struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}
