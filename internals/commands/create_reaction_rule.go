package commands

import (
	//"context"
	//"time"

	"github.com/bwmarrin/discordgo"
	//"github.com/finkabaj/hyde-bot/internals/logger"
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
	//	msg, err := s.ChannelMessageSendEmbed(i.ChannelID, &discordgo.MessageEmbed{
	//		Title:       "Vaserization",
	//		Description: "React to this message with the reaction you want to add to the reaction rules",
	//		Color:       0x538ac4,
	//		Fields: []*discordgo.MessageEmbedField{
	//			{
	//				Name:   "Approve",
	//				Value:  "✅",
	//				Inline: true,
	//			},
	//		},
	//	})
	//
	//	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	//		Type: discordgo.InteractionResponseChannelMessageWithSource,
	//		Data: &discordgo.InteractionResponseData{
	//			Content: "React to this message with the reaction you want to add to the reaction rules",
	//			Flags:   1 << 6,
	//			Components: []discordgo.MessageComponent{
	//				discordgo.ActionsRow{
	//					Components: []discordgo.MessageComponent{
	//						discordgo.Button{
	//							Emoji: discordgo.ComponentEmoji{
	//								Name: "✅",
	//							},
	//							Style:    discordgo.SuccessButton,
	//							Label:    "Approve",
	//							CustomID: "rr_approve",
	//						},
	//						discordgo.Button{
	//							Emoji: discordgo.ComponentEmoji{
	//								Name: "❌",
	//							},
	//							Style:    discordgo.DangerButton,
	//							Label:    "Deny",
	//							CustomID: "rr_deny",
	//						},
	//					},
	//				},
	//			},
	//		},
	//	})
	//
	//	if err != nil {
	//		logger.Error(err, logger.LogFields{"guildId": i.GuildID})
	//		return
	//	}
	//
	// sleepCtx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	//
	//	defer func() {
	//		if sleepCtx.Err() == context.DeadlineExceeded {
	//			err = s.ChannelMessageDelete(i.Interaction.ChannelID, i.Interaction.Message.ID)
	//
	//			if err != nil {
	//				logger.Error(err, logger.LogFields{"guildId": i.GuildID})
	//			}
	//		}
	//		cancel()
	//	}()
}
