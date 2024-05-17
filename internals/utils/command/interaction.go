package commandUtils

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
)

func SendDefaultResponse(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
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
		logger.Debug("Sent response", FillFields(i))
	}
}

func FillFields(i *discordgo.InteractionCreate) map[string]any {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		optionFields := make(map[string]any)

		for _, option := range i.ApplicationCommandData().Options {
			optionFields[option.Name] = option.Value
		}

		return map[string]any{
			"InteractionType":    i.Type,
			"InteractionName":    i.ApplicationCommandData().Name,
			"InteractionOptions": optionFields,
			"InteractionTarget":  i.ApplicationCommandData().TargetID,
		}
	case discordgo.InteractionModalSubmit:
		data := i.ModalSubmitData()

		modalFields := make(map[string]any)

		for _, component := range data.Components {
			switch component.(type) {
			case *discordgo.ActionsRow:
				actionsRow := component.(*discordgo.ActionsRow)
				for _, action := range actionsRow.Components {
					switch action.(type) {
					case *discordgo.TextInput:
						textInput := action.(*discordgo.TextInput)
						modalFields["TextInput"] = textInput.Value
					case *discordgo.SelectMenu:
						selectMenu := action.(*discordgo.SelectMenu)
						modalFields["SelectMenu"] = selectMenu.Options
					}
				}
			}
		}

		return map[string]any{
			"InteractionType": i.Type,
			"InteractionData": modalFields,
		}
	default:
		return map[string]any{"InteractionType": "Unknown interaction type"}
	}
}
