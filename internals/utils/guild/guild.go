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

func (g GuildCreate) Compare(a GuildCreate) int {
	if g.GuildId != a.GuildId {
		return -1
	}

	if g.OwnerId != a.OwnerId {
		return -1
	}

	return 0
}

var (
	ErrGuildConflict = errors.New("guild already exists")
	ErrEmptyGuildId  = errors.New("guild id not provided")
)
