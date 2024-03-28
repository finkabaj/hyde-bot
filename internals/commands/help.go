package commands

import (
	"github.com/PapicBorovoi/hyde-bot/internals/logger"
	"github.com/bwmarrin/discordgo"
)

var dmHelpPermission = true
var memberHelpPermission int64 = discordgo.PermissionAllText

var HelpCommand = &discordgo.ApplicationCommand{
	Name:         "help",
	Description:  "Shows a list of available commands",
	Type:         discordgo.ChatApplicationCommand,
	DMPermission: &dmHelpPermission,
	DefaultMemberPermissions: &memberHelpPermission,
}

func HelpCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cmdManager *CommandManager) {
	commands := "Available commands:\n"

	for _, command := range cmdManager.Commands {
		commands += "/" + command.ApplicationCommand.Name + " - " + command.ApplicationCommand.Description + "\n"
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: commands,
			Flags:  1 << 6,
		},
	},
	)

	if err != nil {
		logger.Error(err, "Error responding to the help command")
	}
}