package guild

import "errors"

type GuildCreate struct {
	GuildId string `json:"guildId" validate:"required,len=19"`
	OwnerId string `json:"ownerId" validate:"required,min=17,max=18"`
	//Icon    string `json:"icon"`
}

type Guild struct {
	GuildId string `json:"guildId"`
	OwnerId string `json:"ownerId"`
}

var (
	ErrGuildNotFound = errors.New("guild not found")
	ErrGuildConflict = errors.New("guild already exists")
	EmptyGuildId     = errors.New("guild id not provided")
)
