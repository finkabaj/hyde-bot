package guild

import "errors"

type GuildCreate struct {
	GuildId string `json:"guildId,inline" validate:"required"`
	OwnerId string `json:"ownerId,inline" validate:"required"`
}

type Guild struct {
	GuildId string `json:"guildId"`
	OwnerId string `json:"ownerId"`
}

var (
	ErrGuildConflict = errors.New("guild already exists")
	EmptyGuildId     = errors.New("guild id not provided")
)
