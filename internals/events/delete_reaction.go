package events

import (
	"slices"

	"github.com/bwmarrin/discordgo"
)

var prohibitedEmojies = []string{"🔥"}

func HandleDeleteReaction(s *discordgo.Session, event interface{}) {
	typedEvent := event.(*discordgo.MessageReactionAdd)

	contains := slices.Contains(prohibitedEmojies, typedEvent.Emoji.Name)

	if contains {
		s.MessageReactionsRemoveEmoji(typedEvent.ChannelID, typedEvent.MessageID, typedEvent.Emoji.Name)
	}
}
