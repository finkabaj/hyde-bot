package events

import (
	"errors"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/enescakir/emoji"
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

		reactionRules, err := rm.GetReactionRules(typedEvent.GuildID, false)

		if err != nil && !errors.Is(err, rules.ErrRulesNotFound) {
			logger.Debug("Failed to get reaction rules:" + err.Error())
			return
		} else if errors.Is(err, rules.ErrRulesNotFound) {
			return
		}

		if slices.ContainsFunc(reactionRules, func(rule rule.ReactionRule) bool {
			return rule.EmojiId != "" && typedEvent.Emoji.ID != "" && rule.EmojiId == typedEvent.Emoji.ID
		}) {
			e := ""
			if typedEvent.Emoji.ID != "" {
				e = typedEvent.Emoji.Name + ":" + typedEvent.Emoji.ID
			} else {
				e = typedEvent.Emoji.Name
			}

			err = s.MessageReactionsRemoveEmoji(typedEvent.ChannelID, typedEvent.MessageID, e)

			if err != nil {
				logger.Error(err, map[string]any{"emojiName": typedEvent.Emoji.Name, "emojiId": typedEvent})
			}

			return
		}

		if slices.ContainsFunc(reactionRules, func(rule rule.ReactionRule) bool {
			return rule.EmojiName != "" && typedEvent.Emoji.Name != "" && emoji.Parse(rule.EmojiName) == typedEvent.Emoji.Name
		}) {
			e := ""
			if typedEvent.Emoji.ID != "" {
				e = typedEvent.Emoji.Name + ":" + typedEvent.Emoji.ID
			} else {
				e = typedEvent.Emoji.Name
			}

			err = s.MessageReactionsRemoveEmoji(typedEvent.ChannelID, typedEvent.MessageID, e)

			if err != nil {
				logger.Error(err, map[string]any{"emojiName": typedEvent.Emoji.Name, "emojiId": typedEvent.Emoji.ID})
			}
		}
	}
}
