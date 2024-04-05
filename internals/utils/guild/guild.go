package guild

import ()

type GuildCreate struct {
	GuildId string `json:"guildId"`
	OwnerId string `json:"ownerId"`
}

type Guild struct {
	GuildId string
	OwnerId string
}
