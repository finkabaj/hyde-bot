package guild

import ()

type GuildCreate struct {
	GuildId string `json:"guildId"`
	OwnerId string `json:"ownerId"`
	//Icon    string `json:"icon"`
}

type Guild struct {
	GuildId string `json:"guildId"`
	OwnerId string `json:"ownerId"`
}
