package events

import (
	"net/http"
	"os"
	"reflect"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/commands"
	"github.com/finkabaj/hyde-bot/internals/ranks"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

type EventHandler func(s *discordgo.Session, event interface{})

type Event struct {
	Type    string
	Handler EventHandler
	GuildID string
}

type EventManager struct {
	ruleManager         *rules.RuleManager
	rankManager         *ranks.RankManager
	cm                  *commands.CommandManager
	messageInteractions *commandUtils.MessageInteractions
	Events              map[string]map[string]*Event // Events[type][guildID] = event
	client              *http.Client
}

var em *EventManager

func NewEventManager(ruleManager *rules.RuleManager, cm *commands.CommandManager,
	client *http.Client, messageInteractions *commandUtils.MessageInteractions, rankManager *ranks.RankManager) *EventManager {
	if em == nil {
		return &EventManager{
			ruleManager:         ruleManager,
			rankManager:         rankManager,
			cm:                  cm,
			messageInteractions: messageInteractions,
			Events:              make(map[string]map[string]*Event),
			client:              client,
		}
	}
	return em
}

func (em *EventManager) RegisterDefaultEvents() {
	var guildID string

	if os.Getenv("ENV") == "development" {
		guildID = os.Getenv("DEV_GUILD_ID")
	}

	em.RegisterEventHandler("MessageReactionAdd", HandleDeleteReaction(em.ruleManager), guildID)
	em.RegisterEventHandler("InteractionCreate", HandleInteractionCreate(em.cm), guildID)
	em.RegisterEventHandler("GuildCreate", HandleGuildCreate(em.ruleManager, em.client), "")
	em.RegisterEventHandler("ModalSubmitReaction", HandleSubmitModalReaction(em.ruleManager), guildID)
	em.RegisterEventHandler("MessageSubmitDeleteReactions", HandleSubmitDeleteReactionModal(em.ruleManager, em.messageInteractions), guildID)
	em.RegisterEventHandler("ModalSubmitActivateRankSystem", HandleSubmitActivateRankSystem(em.rankManager, em.ruleManager), guildID)
}

// RegisterEventHandler registers an event handler for a specific guild.
// If guildID is empty, the event handler will be registered globally.
func (em *EventManager) RegisterEventHandler(eventType string, handler EventHandler, guildID string) {
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
func (em *EventManager) RemoveEventHandler(eventType string, handler EventHandler, guildID string) {
	if _, ok := em.Events[eventType]; ok {
		delete(em.Events[eventType], guildID)
		if len(em.Events[eventType]) == 0 {
			delete(em.Events, eventType)
		}
	}
}

// HandleEvent handles an incoming event by calling the appropriate event handlers.
func (em *EventManager) HandleEvent(s *discordgo.Session, event interface{}) {
	eventType := getEventType(event)
	guildID := getGuildID(event)

	if eventHandlers, ok := em.Events[eventType]; ok {
		if eventHandler, ok := eventHandlers[guildID]; ok {
			go eventHandler.Handler(s, event)
		} else if globalEventHandler, ok := eventHandlers[""]; ok {
			go globalEventHandler.Handler(s, event)
		}
	}
}

// getEventType returns the type of the event based on its underlying struct.
func getEventType(event interface{}) string {
	t := reflect.TypeOf(event)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if i, ok := event.(*discordgo.InteractionCreate); ok {
		switch i.Type {
		case discordgo.InteractionModalSubmit:
			data := i.ModalSubmitData()

			switch data.CustomID {
			case "emoji_ban" + i.Member.User.ID:
				return "ModalSubmitReaction"
			case "activate_role_system" + i.Member.User.ID:
				return "ModalSubmitActivateRankSystem"
			}
		case discordgo.InteractionMessageComponent:
			return "MessageSubmitDeleteReactions"
		}
	}

	return t.Name()
}

// getGuildID returns the guild ID associated with the event, if applicable.
func getGuildID(event any) string {
	v := reflect.ValueOf(event)

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		field := v.FieldByName("GuildID")
		if field.IsValid() {
			return field.String()
		}
	}

	if e, ok := event.(*discordgo.GuildCreate); ok {
		return e.ID
	}

	return ""
}
