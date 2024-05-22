package events

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/commands"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

func HandleSubmitDeleteReactionModal(rm *rules.RuleManager, messageInteractions *commandUtils.MessageInteractions) EventHandler {
	return func(s *discordgo.Session, event any) {
		i, ok := event.(*discordgo.InteractionCreate)

		if !ok {
			return
		}

		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		data := i.MessageComponentData()

		if data.CustomID != "delete_reaction_rules" {
			return
		}

		deleteRulesDto := make([]rules.RulesDeleteDto, len(data.Values))
		for _, v := range data.Values {
			vSplit := strings.Split(v, ":")
			emojiName := vSplit[0]
			emojiID := vSplit[1]

			if emojiID == "NULL" {
				emojiID = ""
			}

			deleteRulesDto = append(deleteRulesDto, rules.RulesDeleteDto{
				EmojiName: emojiName,
				EmojiId:   emojiID,
			})
		}

		err := rm.DeleteReactionRulesApi(i.GuildID, deleteRulesDto)

		if err != nil {
			logger.Error(err, map[string]any{"details": "failed to delete reaction rules"})

			i, ok := messageInteractions.GetMessageInteraction(i.Member.User.ID)

			if !ok {
				return
			}

			content := "Failed to delete reaction rules"
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content:    &content,
				Components: &[]discordgo.MessageComponent{},
			})

			if err != nil {
				logger.Error(err, map[string]any{"details": "failed to edit interaction response"})
				return
			}

			messageInteractions.DeleteMessageID(i.Member.User.ID)

			return
		}

		err = commands.DeleteMessageReactionDeleteSelect(s, i, messageInteractions)

		if err != nil {
			logger.Error(err, map[string]any{"details": "failed to delete message"})
		}
	}
}
