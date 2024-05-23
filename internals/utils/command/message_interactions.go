package commandUtils

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type MessageInteractions struct {
	messageInteractions map[string]*discordgo.InteractionCreate
	lock                sync.RWMutex
}

func NewMessageInteractions() *MessageInteractions {
	return &MessageInteractions{
		messageInteractions: make(map[string]*discordgo.InteractionCreate),
	}
}

func (m *MessageInteractions) GetMessageInteraction(userID string) (*discordgo.InteractionCreate, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	i, ok := m.messageInteractions[userID]
	return i, ok
}

func (m *MessageInteractions) SetMessageID(userID string, i *discordgo.InteractionCreate) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.messageInteractions[userID] = i
}

func (m *MessageInteractions) DeleteMessageID(userID string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.messageInteractions, userID)
}
