package ranks

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Rank struct {
	Level uint8
	XP    uint16
	Role  *discordgo.Role
}

type Ranks struct {
	guildID string
	ownerID string
	ranks   []Rank
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
		guildID: guildID,
		ownerID: ownerID,
		ranks:   ranks,
	}
}
