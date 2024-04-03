package events

import (
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
)

var prohibitedEmojies = []string{"🔥"}

func HandleDeleteReaction(s *discordgo.Session, event interface{}) {
	typedEvent, ok := event.(*discordgo.MessageReactionAdd)

	if !ok {
		logger.Debug("Failed to cast event to *discordgo.MessageReactionAdd")
		return
	}

	contains := slices.Contains(prohibitedEmojies, typedEvent.Emoji.Name)

	if contains {
		s.MessageReactionsRemoveEmoji(typedEvent.ChannelID, typedEvent.MessageID, typedEvent.Emoji.Name)
	}
}
