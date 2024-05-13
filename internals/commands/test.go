package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
)

var dmTestPermission = false
var testPermission int64 = discordgo.PermissionAdministrator

var testCommand = &discordgo.ApplicationCommand{
	Name:                     "test",
	Description:              "Test command",
	Type:                     discordgo.ChatApplicationCommand,
	DMPermission:             &dmTestPermission,
	DefaultMemberPermissions: &testPermission,
}

func TestHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
