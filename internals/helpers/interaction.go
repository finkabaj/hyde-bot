package helpers

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/sirupsen/logrus"
)

func DefaultErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   1 << 6,
		},
	},
	)

	if err != nil {
		logger.Error(err, FillFields(i))
		return
	}

	if os.Getenv("ENV") == "development" {
		logger.Debug("Sent error response", FillFields(i))
	}
}

func FillFields(i *discordgo.InteractionCreate) logrus.Fields {
	if i.Type == discordgo.InteractionApplicationCommand {
		optionFields := make(logrus.Fields)

		for _, option := range i.ApplicationCommandData().Options {
			optionFields[option.Name] = option.Value
		}

		return logrus.Fields{
			"InteractionType":    i.Type,
			"InteractionName":    i.ApplicationCommandData().Name,
			"InteractionOptions": optionFields,
			"InteractionTarget":  i.ApplicationCommandData().TargetID,
		}
	}

	return logrus.Fields{
		"InteractionType":    i.Type,
		"InteractionMessage": i.Message.Content,
	}
}
