package events

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/enescakir/emoji"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
)

func HandleSumbitModalReaction(rm *rules.RuleManager) EventHandler {
	return func(s *discordgo.Session, event any) {
		i, ok := event.(*discordgo.InteractionCreate)

		if !ok {
			logger.Error(errors.New("failed to cast event to *discordgo.InteractionCreate"))
			return
		}

		if i.Type != discordgo.InteractionModalSubmit {
			logger.Error(errors.New("incorect type in HandleSumbitModalReaction"))
		}

		data := i.ModalSubmitData()

		if !strings.HasPrefix(string(data.CustomID), "emoji_ban") {
			return
		}

		text := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		emojies, err := s.GuildEmojis(i.GuildID)
		if err != nil {
			logger.Error(err)
			return
		}

		for _, emoji := range emojies {
			fmt.Println(emoji.Name, emoji.ID)
		}

		r := parseModalReactionInput(text, i.Member.User.ID, i.GuildID, emojies)

		if len(r) == 0 {
			commandUtils.SendDefaultResponse(s, i, "No valid emojies found in the input")

			return
		}

		rRules, err := rm.PostReactionRules(i.GuildID, r)

		fmt.Printf("rRules: %+v\n", rRules)

		if err != nil {
			if errors.Is(err, rules.ErrIntersectingRules) {
				commandUtils.SendDefaultResponse(s, i, "Reaction rules already exist")

				return
			}

			commandUtils.SendDefaultResponse(s, i, "Failed to post reaction rules")

			logger.Error(err)

			return
		}

		rm.AddReactionRules(i.GuildID, rRules)

		commandUtils.SendDefaultResponse(s, i, "Modal submitted successfully!")
	}
}

func parseModalReactionInput(text string, ruleAuthor string, guildId string, emojies []*discordgo.Emoji) []rule.ReactionRule {
	if len(emojies) == 0 {
		return nil
	}

	textSplited := strings.Split(text, " ")

	if len(textSplited) == 0 {
		return nil
	}

	result := make([]rule.ReactionRule, 0, len(textSplited))

	for _, v := range textSplited {
		if common.ContainsFieldValue(emojies, "ID", v) {
			result = append(result, rule.ReactionRule{
				RuleAuthor: ruleAuthor,
				GuildId:    guildId,
				EmojiId:    v,
				Actions:    []rule.ReactAction{rule.Delete},
			})
		} else if common.ContainsFieldValue(emojies, "Name", v) {
			result = append(result, rule.ReactionRule{
				RuleAuthor: ruleAuthor,
				GuildId:    guildId,
				EmojiName:  v,
				Actions:    []rule.ReactAction{rule.Delete},
			})
		} else if emoji.Exist(v) {
			result = append(result, rule.ReactionRule{
				RuleAuthor: ruleAuthor,
				GuildId:    guildId,
				EmojiName:  v,
				Actions:    []rule.ReactAction{rule.Delete},
			})
		}
	}

	return result
}
