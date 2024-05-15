package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/commands"
	"github.com/finkabaj/hyde-bot/internals/logger"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

func HandleInteractionCreate(cm *commands.CommandManager) EventHandler {
	return func(s *discordgo.Session, event interface{}) {
		i := event.(*discordgo.InteractionCreate)

		if i.Type == discordgo.InteractionApplicationCommand {
			cmd, err := cm.GetCommandByName(i.ApplicationCommandData().Name, i.GuildID)

			if err != nil {
				cmd, err = cm.GetCommandByName(i.ApplicationCommandData().Name, "")

				if err != nil {
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Command is not registered",
							Flags:   1 << 6,
						},
					})

					if err != nil {
						logger.Error(err, commandUtils.FillFields(i))
					}
				}
			}

			cmd.Handler(s, i)
			logger.Info("Command executed", commandUtils.FillFields(i))
		}
	}
}
