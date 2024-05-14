package events

import (
	"errors"
	"fmt"
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

		reactionRules, err := rm.GetReactionRules(typedEvent.GuildID)

		fmt.Printf("reactionRules: %+v\n", reactionRules)

		if err != nil && !errors.Is(err, rules.ErrRulesNotFound) {
			logger.Debug("Failed to get reaction rules:" + err.Error())
			return
		} else if errors.Is(err, rules.ErrRulesNotFound) {
			return
		}

		fmt.Println(typedEvent.Emoji.ID)
		fmt.Println(typedEvent.Emoji.Name)

		fmt.Println(emoji.Parse(typedEvent.Emoji.Name))

		if slices.ContainsFunc(reactionRules, func(rule rule.ReactionRule) bool {
			return rule.EmojiId != "" && typedEvent.Emoji.ID != "" && rule.EmojiId == typedEvent.Emoji.ID
		}) {
			s.MessageReactionsRemoveEmoji(typedEvent.ChannelID, typedEvent.MessageID, typedEvent.Emoji.ID)
			return
		}

		if slices.ContainsFunc(reactionRules, func(rule rule.ReactionRule) bool {
			return rule.EmojiName != "" && typedEvent.Emoji.Name != "" && emoji.Parse(rule.EmojiName) == typedEvent.Emoji.Name
		}) {
			s.MessageReactionRemove(typedEvent.ChannelID, typedEvent.MessageID, emoji.Parse(typedEvent.Emoji.Name), typedEvent.UserID)
			return
		}
	}
}
