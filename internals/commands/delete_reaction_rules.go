package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

var dmDeleteReactionRulesPermision = false

var deleteReactionRulesPermission int64 = discordgo.PermissionAdministrator

var DeleteReactionRuleCommand = &discordgo.ApplicationCommand{
	Name:                     "delete-reaction-rules",
	Description:              "Delete reaction rules for the server",
	Type:                     discordgo.ChatApplicationCommand,
	DMPermission:             &dmDeleteReactionRulesPermision,
	DefaultMemberPermissions: &deleteReactionRulesPermission,
}

func DeleteMessageReactionDeleteSelect(s *discordgo.Session, i *discordgo.InteractionCreate,
	messageInteractions *commandUtils.MessageInteractions) error {
	i, ok := messageInteractions.GetMessageInteraction(i.Member.User.ID)

	if !ok {
		return nil
	}

	err := s.InteractionResponseDelete(i.Interaction)

	if err != nil {
		return fmt.Errorf("failed to delete message in deleteMessageReationRules: %w", err)
	}

	messageInteractions.DeleteMessageID(i.Member.User.ID)

	return nil
}

// in my opinion, interacion delete and edit works strange in this handler but who cares

func DeleteReactionRulesCommandHandler(s *discordgo.Session,
	i *discordgo.InteractionCreate, rm *rules.RuleManager, messageInteractions *commandUtils.MessageInteractions) {
	err := DeleteMessageReactionDeleteSelect(s, i, messageInteractions)

	if err != nil {
		logger.Error(err, map[string]any{"details": "failed to delete message"})
	}

	opts, err := createSelectMenuOptions(rm, i.GuildID)

	if err != nil {
		logger.Error(err, map[string]any{"details": "failed to create select menu options"})
		commandUtils.SendDefaultResponse(s, i, "Failed to create select menu options")
		return
	}

	minValue := 1

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Select the reaction rules you want to delete",
			Flags:   1 << 6,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "delete_reaction_rules",
							Placeholder: "Select tags to search on StackOverflow",
							MinValues:   &minValue,
							MaxValues:   len(opts),
							Options:     opts,
						},
					},
				},
			},
		},
	})

	if err != nil {
		logger.Error(err, map[string]any{"details": "failed to respond to delete reaction rules command"})
		return
	}

	messageInteractions.SetMessageID(i.Member.User.ID, i)

	go func() {
		<-time.After(30 * time.Second)

		if i, ok := messageInteractions.GetMessageInteraction(i.Member.User.ID); ok {
			content := "You took too long to respond, please try again."
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content:    &content,
				Components: &[]discordgo.MessageComponent{},
			})

			if err != nil {
				//if user spams the command, the bot will try to edit the message that was already deleted
				//so this error is expected
				//don't try to add messageInteractions.DeleteMessageID(i.Member.User.ID) here
				//this will cause last message to not be deleted
				logger.Error(err, map[string]any{"details": "failed to follow up message in delete reaction rules command"})
			}
		}
	}()
}

func createSelectMenuOptions(rm *rules.RuleManager, guildId string) ([]discordgo.SelectMenuOption, error) {
	options := make([]discordgo.SelectMenuOption, 0)

	rRules, err := rm.GetReactionRules(guildId, false)

	if err != nil {
		return nil, fmt.Errorf("failed to get reaction rules in createSelectMenuOptions: %w", err)
	}

	for _, rule := range rRules {
		if rule.IsCustom {
			options = append(options, discordgo.SelectMenuOption{
				Label: "server emoji",
				Value: fmt.Sprintf("%s:%s", rule.EmojiName, rule.EmojiId),
				Emoji: discordgo.ComponentEmoji{
					ID:   rule.EmojiId,
					Name: rule.EmojiName,
				},
			})

		} else {
			options = append(options, discordgo.SelectMenuOption{
				Label: "ordinary emoji",
				Value: fmt.Sprintf("%s:NULL", rule.EmojiName),
				Emoji: discordgo.ComponentEmoji{
					Name: rule.EmojiName,
				},
			})
		}
	}

	return options, nil
}
