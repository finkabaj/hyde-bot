package events

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
)

func HandleSumbitModalReaction(s *discordgo.Session, event any) {
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

	rules := parseModalReactionInput(text, i.Member.User.ID, i.GuildID, emojies)

	if len(rules) == 0 {
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No valid emojies found in the input",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		if err != nil {
			logger.Error(err)
		}

		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Modal submitted successfully!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Error(err)
		return
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
		}
	}

	return result
}
