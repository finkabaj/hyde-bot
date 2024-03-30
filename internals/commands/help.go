package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/helpers"
	"github.com/finkabaj/hyde-bot/internals/logger"
)

var dmHelpPermission = true
var memberHelpPermission int64 = discordgo.PermissionAllText

var HelpCommand = &discordgo.ApplicationCommand{
	Name:                     "help",
	Description:              "Shows a list of available commands",
	Type:                     discordgo.ChatApplicationCommand,
	DMPermission:             &dmHelpPermission,
	DefaultMemberPermissions: &memberHelpPermission,
}

func HelpCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands := "Available commands:\n"

	cmdManager := NewCommandManager()

	for _, command := range cmdManager.Commands {
		if command[i.Interaction.GuildID] != nil && command[i.Interaction.GuildID].IsRegistered {
			commands += "/" + command[i.Interaction.GuildID].ApplicationCommand.Name + "\n"
		} else if command[""] != nil && command[""].IsRegistered {
			commands += "/" + command[""].ApplicationCommand.Name + "\n"
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: commands,
			Flags:   1 << 6,
		},
	},
	)

	if err != nil {
		logger.Error(err, helpers.FillFields(i))
	}
}
