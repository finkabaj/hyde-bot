package events

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

type EventHandler func(s *discordgo.Session, event interface{})

type Event struct {
	Type    string
	Handler EventHandler
	GuildID string
}

type eventManager struct {
	Events map[string]map[string]*Event // Events[type][guildID] = event
}

var em *eventManager

func NewEventManager() *eventManager {
	if em == nil {
		return &eventManager{
			Events: make(map[string]map[string]*Event),
		}
	}
	return em
}

func (em *eventManager) RegisterDefaultEvents() {
	var guildID string = ""

	if os.Getenv("ENV") == "development" {
		guildID = os.Getenv("DEV_GUILD_ID")
	}

	em.RegisterEventHandler("MessageReactionAdd", HandleDeleteReaction, guildID)
	em.RegisterEventHandler("InteractionCreate", HandleInteractionCreate, guildID)
	em.RegisterEventHandler("GuildCreate", HandleGuildCreate, "")
}

// RegisterEventHandler registers an event handler for a specific guild.
// If guildID is empty, the event handler will be registered globally.
func (em *eventManager) RegisterEventHandler(eventType string, handler EventHandler, guildID string) {
	event := &Event{
		Type:    eventType,
		Handler: handler,
		GuildID: guildID,
	}

	if _, ok := em.Events[eventType]; !ok {
		em.Events[eventType] = make(map[string]*Event)
	}

	em.Events[eventType][guildID] = event
}

// RemoveEventHandler removes an event handler for a specific guild.
// If guildID is empty, it will remove the global event handler.
func (em *eventManager) RemoveEventHandler(eventType string, handler EventHandler, guildID string) {
	if _, ok := em.Events[eventType]; ok {
		delete(em.Events[eventType], guildID)
		if len(em.Events[eventType]) == 0 {
			delete(em.Events, eventType)
		}
	}
}

// HandleEvent handles an incoming event by calling the appropriate event handlers.
func (em *eventManager) HandleEvent(s *discordgo.Session, event interface{}) {
	eventType := getEventType(event)
	guildID := getGuildID(event)

	if eventHandlers, ok := em.Events[eventType]; ok {
		if eventHandler, ok := eventHandlers[guildID]; ok {
			eventHandler.Handler(s, event)
		} else if globalEventHandler, ok := eventHandlers[""]; ok {
			globalEventHandler.Handler(s, event)
		}
	}
}

// getEventType returns the type of the event based on its underlying struct.
func getEventType(event interface{}) string {
	switch event.(type) {
	case *discordgo.InteractionCreate:
		return "InteractionCreate"
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
	case *discordgo.GuildCreate:
		return "GuildCreate"
	default:
		return ""
	}
}

// getGuildID returns the guild ID associated with the event, if applicable.
func getGuildID(event interface{}) string {
	switch e := event.(type) {
	case *discordgo.InteractionCreate:
		return e.Interaction.GuildID
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
	case *discordgo.GuildCreate:
		return e.ID
	default:
		return ""
	}
}
