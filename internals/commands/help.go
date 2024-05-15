package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
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

func HelpCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cmdManager *CommandManager) {
	commands := "Available commands:\n"

	cmds := cmdManager.GetGuildCommands(i.GuildID)

	if len(cmds) == 0 {
		commands = "No commands available"
	}

	for _, command := range cmds {
		if command.IsRegistered {
			commands += "/" + command.ApplicationCommand.Name + "\n"
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
		logger.Error(err, commandUtils.FillFields(i))
	}
}
