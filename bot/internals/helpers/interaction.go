package helpers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/PapicBorovoi/hyde-bot/internals/logger"
)

func DefaultErrorResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:  1 << 6,
		},
	},
	)

	if err != nil {
		logger.Error(err, "Error while error responding to the interaction???")
	}
}