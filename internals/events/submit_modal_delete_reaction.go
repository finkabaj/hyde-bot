package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/commands"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

func HandleSubmitDeleteReactionModal(rm *rules.RuleManager, messageInteractions *commandUtils.MessageInteractions) EventHandler {
	return func(s *discordgo.Session, event any) {
		i, ok := event.(*discordgo.InteractionCreate)

		if !ok {
			return
		}

		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		data := i.MessageComponentData()

		if data.CustomID != "delete_reaction_rules" {
			return
		}

		fmt.Printf("%+v\n", data)

		logger.Info("Delete reaction rules modal submitted", map[string]any{"user": i.Member.User.ID})

		commands.DeleteMessageReactionDeleteSelect(s, i, messageInteractions)
	}
}
