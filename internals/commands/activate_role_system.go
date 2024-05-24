package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
)

var dmActivateRankSystemPermission = false
var activateRankSystemPermission int64 = discordgo.PermissionManageRoles

var ActivateRankSystemCommand = &discordgo.ApplicationCommand{
	Name: "activate-rank-system",
	Description: `Activate the role system for the server. 
  You will be able to rank up using this system`,
	Type:                     discordgo.ChatApplicationCommand,
	DMPermission:             &dmActivateRankSystemPermission,
	DefaultMemberPermissions: &activateRankSystemPermission,
}

func ActivateRankSystemHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "activate_role_system" + i.Member.User.ID,
			Title:    "Activate role system",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "role_system",
							Label:       "roles",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Ids of roles separated by space in order of rank\nNOTE if role system exists it will be overwritten",
							Required:    true,
							MaxLength:   600,
							MinLength:   1,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "xp_system",
							Label:       "xp",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "provide one or multiple xps for each rank separated by space",
							Required:    true,
							MaxLength:   600,
							MinLength:   1,
						},
					},
				},
			},
		},
	})

	if err != nil {
		logger.Error(err, map[string]any{"details": "Error responding to interaction when activating role system"})
		return
	}
}
