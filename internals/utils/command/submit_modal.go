package commandUtils

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func GetDataFromModalSubmit(event any) (*discordgo.ModalSubmitInteractionData, *discordgo.InteractionCreate, error) {
	i, ok := event.(*discordgo.InteractionCreate)

	if !ok {
		return nil, nil, errors.New("failed to cast event to *discordgo.InteractionCreate")
	}

	if i.Type != discordgo.InteractionModalSubmit {
		return nil, nil, fmt.Errorf("incorect type, want: %s get: %s", discordgo.InteractionModalSubmit, i.Type)
	}

	data := i.ModalSubmitData()

	return &data, i, nil
}
