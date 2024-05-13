package events

import (
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
)

func HandleDeleteReaction(rm *rules.RuleManager) EventHandler {
	return func(s *discordgo.Session, event any) {
		typedEvent, ok := event.(*discordgo.MessageReactionAdd)

		if !ok {
			logger.Debug("Failed to cast event to *discordgo.MessageReactionAdd")
			return
		}

		reactionRules, err := rm.GetReactionRules(typedEvent.GuildID)

		if err != nil {
			logger.Debug("Failed to get reaction rules:" + err.Error())
			return
		}

		if slices.ContainsFunc(reactionRules, func(rule rule.ReactionRule) bool {
			return typedEvent.Emoji.ID != "" && rule.EmojiId == typedEvent.Emoji.ID
		}) {
			s.MessageReactionsRemoveEmoji(typedEvent.ChannelID, typedEvent.MessageID, typedEvent.Emoji.ID)
			return
		}

		if slices.ContainsFunc(reactionRules, func(rule rule.ReactionRule) bool {
			return typedEvent.Emoji.Name != "" && rule.EmojiName == typedEvent.Emoji.Name
		}) {
			s.MessageReactionsRemoveEmoji(typedEvent.ChannelID, typedEvent.MessageID, typedEvent.Emoji.Name)
			return
		}
	}
}
