package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/helpers"
	"github.com/finkabaj/hyde-bot/internals/logger"
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

func DeleteCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandName := i.ApplicationCommandData().Options[0].StringValue()
	content := ""

	cm := NewCommandManager()

	command, ok := cm.Commands[commandName]

	if !ok {
		content = "Command not found"
	}

	var c *Command

	if c, ok = command[i.Interaction.GuildID]; ok && !c.IsRegistered {
		if c, ok = command[""]; ok {
			content = "You can't delete default commands"
		} else {
			content = "Command not found"
		}
	}

	if content == "" {
		err := cm.DeleteCommand(s, c.RegisteredCommand, i.Interaction.GuildID)

		if err != nil {
			logger.Error(err, helpers.FillFields(i))
			content = "Error deleting the command"
		}

		logger.Info("Command deleted", helpers.FillFields(i))
		content = "Command deleted"
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   1 << 6,
		},
	},
	)

	if err != nil {
		logger.Error(err, helpers.FillFields(i))
	}
}
