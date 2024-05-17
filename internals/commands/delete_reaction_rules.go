package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

var dmDeleteReactionRulesPermision = false
var deleteReactionRulesPermission int64 = discordgo.PermissionAdministrator

var DeleteReactionRuleCommand = &discordgo.ApplicationCommand{
	Name:                     "девасеризация",
	Description:              "Delete reaction rules",
	Type:                     discordgo.ChatApplicationCommand,
	DMPermission:             &dmDeleteReactionRulesPermision,
	DefaultMemberPermissions: &deleteReactionRulesPermission,
}

func DeleteReactionRulesCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, rm *rules.RuleManager) {
	opts, err := createSelectMenuOptions(rm, i.GuildID)

	if err != nil {
		logger.Error(err, map[string]any{"details": "failed to create select menu options"})
		commandUtils.SendDefaultResponse(s, i, "Failed to create select menu options")
		return
	}

	minValue := 1

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "emoji_unban" + i.Member.User.ID,
			Title:    "Не пизда",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "emoji_unban_select" + i.Member.User.ID,
							Placeholder: "select emoji to delete rules",
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
	}
}

func createSelectMenuOptions(rm *rules.RuleManager, guildId string) ([]discordgo.SelectMenuOption, error) {
	options := make([]discordgo.SelectMenuOption, 0)

	rRules, err := rm.GetReactionRules(guildId, false)

	if err != nil {
		return nil, fmt.Errorf("failed to get reaction rules in createSelectMenuOptions: %w", err)
	}

	for _, rule := range rRules {
		if rule.EmojiId != "" {
			options = append(options, discordgo.SelectMenuOption{
				Label: rule.EmojiId,
				Value: rule.EmojiId,
			})

		} else {
			options = append(options, discordgo.SelectMenuOption{
				Label: rule.EmojiName,
				Value: rule.EmojiName,
			})
		}
	}

	return options, nil
}
