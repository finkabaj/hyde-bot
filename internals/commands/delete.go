package commands

import (
	"github.com/PapicBorovoi/hyde-bot/internals/logger"
	"github.com/bwmarrin/discordgo"
)

var dmDeletePermission = false
var memberDeletePermission int64 = discordgo.PermissionAdministrator

var DeleteCommand = &discordgo.ApplicationCommand{
	Name:         "delete",
	Description:  "Deletes a command from the bot globally",
	Type:         discordgo.ChatApplicationCommand,
	DMPermission: &dmDeletePermission,
	DefaultMemberPermissions: &memberDeletePermission,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: 			discordgo.ApplicationCommandOptionString,
			Name: 			"command",
			Description: 	"The command to delete",
			Required: 		true,
		},
	},
}

func DeleteCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cmdManager *CommandManager) {
	commandName := i.ApplicationCommandData().Options[0].StringValue()
	content := ""

	for _, command := range cmdManager.RegisteredCommands {
		if command.Name == commandName {
			err := cmdManager.DeleteCommand(s, command, "")
			if err != nil {
				logger.Error(err, "Error deleting the command on" + i.GuildID + " guild" + i.Member.User.ID + " user" + commandName + " command")
				content = "Error deleting the command"
				break
			}
			logger.Info("Command " + commandName + " deleted on" + i.GuildID + " guild" + i.Member.User.ID + " user")
			content = "Command deleted"
			break
		}
	}

	if content == "" {
		content = "Command not found"
	}


	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:  1 << 6,
		},
	},
	)

	if err != nil {
		logger.Error(err, "Error responding to the help command")
	}
}