package ranks

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Rank struct {
	Level uint8           `json:"level" validate:"required,number"`
	XP    uint16          `json:"xp" validate:"required,number"`
	ID    string          `json:"id" validate:"required,uuid4"`
	Role  *discordgo.Role `json:"role"`
}

type Ranks struct {
	GuildID string `json:"guild_id" validate:"required"`
	OwnerID string `json:"owner_id" validate:"required"`
	Ranks   []Rank `json:"ranks" validate:"required,dive"`
}

type RankManager struct {
	lock  sync.RWMutex
	ranks map[string]Ranks
}

var rankManager *RankManager

func NewRankManager() *RankManager {
	if rankManager == nil {
		rankManager = &RankManager{
			ranks: make(map[string]Ranks),
		}
	}

	return rankManager
}

func (rm *RankManager) PostRanks(guildID, ownerID string, ranks []Rank) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	rm.ranks[guildID] = Ranks{
		GuildID: guildID,
		OwnerID: ownerID,
		Ranks:   ranks,
	}
}
