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
		if command.GuildID != "" && command.GuildID != i.Interaction.GuildID || !command.IsRegistered {
			continue
		}
		commands += "/" + command.RegisteredCommand.Name + " - " + command.RegisteredCommand.Description + "\n"
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
