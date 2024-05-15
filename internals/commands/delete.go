package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

var dmDeletePermission = false
var memberDeletePermission int64 = discordgo.PermissionAdministrator

var DeleteCommand = &discordgo.ApplicationCommand{
	Name:                     "delete",
	Description:              "Deletes a command from this guild",
	Type:                     discordgo.ChatApplicationCommand,
	DMPermission:             &dmDeletePermission,
	DefaultMemberPermissions: &memberDeletePermission,
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "command",
			Description: "The command to delete",
			Required:    true,
		},
	},
}

func DeleteCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, cm *CommandManager) {
	commandName := i.ApplicationCommandData().Options[0].StringValue()
	content := ""

	command, err := cm.GetCommandByName(commandName, i.GuildID)

	if err != nil || !command.IsRegistered {
		command, err = cm.GetCommandByName(commandName, "")

		if err != nil {
			content = "Command not found"
		} else {
			content = "You can't delete default commands"
		}
	}

	if content == "" {
		err := cm.DeleteCommand(s, command.RegisteredCommand, i.GuildID, false)

		if err != nil {
			logger.Error(err, commandUtils.FillFields(i))
			content = "Error deleting the command"
		}

		logger.Info("Command deleted", commandUtils.FillFields(i))
		content = "Command deleted"
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   1 << 6,
		},
	},
	)

	if err != nil {
		logger.Error(err, commandUtils.FillFields(i))
	}
}
