package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/commands"
	"github.com/finkabaj/hyde-bot/internals/helpers"
	"github.com/finkabaj/hyde-bot/internals/logger"
)

func HandleInteractionCreate(s *discordgo.Session, event interface{}) {
	i := event.(*discordgo.InteractionCreate)

	cm := commands.NewCommandManager()

	if i.Type == discordgo.InteractionApplicationCommand {
		cmd, ok := cm.Commands[i.ApplicationCommandData().Name][i.Interaction.GuildID]

		if !ok {
			cmd, ok = cm.Commands[i.ApplicationCommandData().Name][""]
		}

		if !ok {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Command is not registered",
					Flags:   1 << 6,
				},
			})

			if err != nil {
				logger.Error(err, helpers.FillFields(i))
			}
		}

		cmd.Handler(s, i)
		logger.Info("Command executed", helpers.FillFields(i))
	}
}
