package events

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

func HandleSubmitDeleteReactionModal(rm *rules.RuleManager) EventHandler {
	return func(s *discordgo.Session, event any) {
		data, i, err := commandUtils.GetDataFromModalSubmit(event)

		if err != nil {
			logger.Error(fmt.Errorf("error at HandleSubmitDeleteReactionModal: %w", err))
		}

		if !strings.HasPrefix(data.CustomID, "emoji_unban") {
			return
		}

		options := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.SelectMenu)

		for _, option := range options.Options {
			fmt.Println(option.Emoji.Name)
			fmt.Println(option.Value)
		}

		commandUtils.SendDefaultResponse(s, i, "Success")
	}
}
