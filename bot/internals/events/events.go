package events

import (
	"github.com/bwmarrin/discordgo"
)

type EventHandler func(s *discordgo.Session, event interface{})

type Event struct {
	Type    string
	Handler EventHandler
	GuildID string
}

type EventManager struct {
	Events map[string][]*Event
}

func NewEventManager() *EventManager {
	return &EventManager{
		Events: make(map[string][]*Event),
	}
}

// RegisterEventHandler registers an event handler for a specific guild.
// If guildID is empty, the event handler will be registered globally.
func (em *EventManager) RegisterEventHandler(eventType string, handler EventHandler, guildID string) {
	event := &Event{
		Type:    eventType,
		Handler: handler,
		GuildID: guildID,
	}
	em.Events[eventType] = append(em.Events[eventType], event)
}

// RemoveEventHandler removes an event handler for a specific guild.
// If guildID is empty, it will remove the global event handler.
func (em *EventManager) RemoveEventHandler(eventType string, handler EventHandler, guildID string) {
	events := em.Events[eventType]
	for i, e := range events {
		if e.Type == eventType && e.GuildID == guildID {
			em.Events[eventType] = append(events[:i], events[i+1:]...)
			break
		}
	}
}

// HandleEvent handles an incoming event by calling the appropriate event handlers.
func (em *EventManager) HandleEvent(s *discordgo.Session, event interface{}) {
	eventType := getEventType(event)
	guildID := getGuildID(event)

	// Call global event handlers
	for _, e := range em.Events[eventType] {
		if e.GuildID == "" {
			e.Handler(s, event)
		}
	}

	// Call guild-specific event handlers
	for _, e := range em.Events[eventType] {
		if e.GuildID == guildID {
			e.Handler(s, event)
		}
	}
}

// getEventType returns the type of the event based on its underlying struct.
func getEventType(event interface{}) string {
	switch event.(type) {
	case *discordgo.MessageCreate:
		return "MessageCreate"
	case *discordgo.MessageUpdate:
		return "MessageUpdate"
	case *discordgo.MessageDelete:
		return "MessageDelete"
	case *discordgo.MessageReactionAdd:
		return "MessageReactionAdd"
	case *discordgo.MessageReactionRemove:
		return "MessageReactionRemove"
	default:
		return ""
	}
}

// getGuildID returns the guild ID associated with the event, if applicable.
func getGuildID(event interface{}) string {
	switch e := event.(type) {
	case *discordgo.MessageCreate:
		return e.GuildID
	case *discordgo.MessageUpdate:
		return e.GuildID
	case *discordgo.MessageDelete:
		return e.GuildID
	case *discordgo.MessageReactionAdd:
		return e.GuildID
	case *discordgo.MessageReactionRemove:
		return e.GuildID
	default:
		return ""
	}
}
