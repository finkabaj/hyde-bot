package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/commands"
)

func HandleInteractionCreate(s *discordgo.Session, event interface{}) {
	i := event.(*discordgo.InteractionCreate)

	cm := commands.NewCommandManager()

	if i.Type == discordgo.InteractionApplicationCommand {
		for _, command := range cm.Commands {
			if command.RegisteredCommand.Name == i.ApplicationCommandData().Name &&
				command.GuildID == i.Interaction.GuildID || command.GuildID == "" {
				command.Handler(s, i)
				break
			}
		}
	}
}
