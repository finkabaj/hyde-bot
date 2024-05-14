package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
)

var dmCreateReactionRulePermission = false
var createReactionRulePermission int64 = discordgo.PermissionAdministrator

var CreateReactionRuleCommand = &discordgo.ApplicationCommand{
	Name:                     "васеризация",
	Description:              "Запрет вассеру пердеть реакциями",
	Type:                     discordgo.ChatApplicationCommand,
	DMPermission:             &dmCreateReactionRulePermission,
	DefaultMemberPermissions: &createReactionRulePermission,
}

func CreateReactionRuleHandler(s *discordgo.Session, i *discordgo.InteractionCreate, rm *rules.RuleManager) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "emoji_ban" + i.Member.User.ID,
			Title:    "пизда вассеридзе",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "emoji_ban",
							Label:       "reactions",
							Style:       discordgo.TextInputShort,
							Placeholder: "provide reactions to ban by name or id, separated by space",
							Required:    true,
							MaxLength:   300,
							MinLength:   1,
						},
					},
				},
			},
		},
	})

	if err != nil {
		logger.Error(err)
		return
	}
}
