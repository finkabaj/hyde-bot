package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
)

func HandleGuildCreate(rm *rules.RuleManager) func(s *discordgo.Session, event any) {
	return func(s *discordgo.Session, event any) {
		typedEvent, ok := event.(*discordgo.GuildCreate)

		if !ok {
			logger.Warn(errors.New("incorect type in HandleGuildCreate"))
			return
		}

		info := &guild.GuildCreate{
			GuildId: typedEvent.ID,
			OwnerId: typedEvent.OwnerID,
		}

		jsonInfo, err := json.Marshal(&info)

		if err != nil {
			logger.Error(err, logger.LogFields{"message": "error while marshalling guild info"})
			return
		}

		url := common.GetApiUrl(os.Getenv("API_HOST"), os.Getenv("API_PORT"), "/guild")

		bodyReader := bytes.NewReader(jsonInfo)
		res, err := http.Post(url, "application/json", bodyReader)

		if err != nil {
			logger.Error(err, logger.LogFields{"message": "error while sending post request on guild create"})
			return
		}

		body := res.Body
		defer body.Close()
		b, err := io.ReadAll(body)

		if err != nil {
			logger.Fatal(err, logger.LogFields{"message": "the bot cannot continue to work correctly", "at": "guild_create"})
		}

		var result guild.Guild

		if err := common.UnmarshalBodyBytes(b, &result); err != nil {
			logger.Error(errors.New("error while unmarshaling guild create"))
		}

		var errRes common.ErrorResponse

		if err := common.UnmarshalBodyBytes(b, &errRes); err != nil {
			logger.Error(errors.New("error while unmarshaling guild create error"))
		}

		if errRes.Error == guild.ErrGuildConflict.Error() {
			logger.Info("Guild already exists", logger.LogFields{"guildId": info.GuildId})
		} else if errRes.Error != "" {
			logger.Error(errors.New(errRes.Error), logger.ToLogFields(errRes.ValidationErrors))
			return
		}

		if result.GuildId == info.GuildId || errRes.Error == guild.ErrGuildConflict.Error() {
			rules, err := fetchRules(info.GuildId, rm)

			if err != nil {
				logger.Error(err, logger.LogFields{"message": "error on fetching rules", "at": "guild_create", "guildId": info.GuildId})
				return
			}

			rm.AddRules(info.GuildId, rules)

			logger.Info("Guild created", logger.LogFields{"guildId": info.GuildId})
		}
	}
}

func fetchRules(guildId string, rm *rules.RuleManager) (rules.Rules, error) {
	rRules, err := rm.FetchReactionRules(guildId)

	if err != nil {
		return rules.Rules{}, err
	}

	if len(rRules) == 0 {
		return rules.Rules{
			ReactionRules:     nil,
			HaveReactionRules: false,
		}, nil
	}

	return rules.Rules{
		ReactionRules:     rRules,
		HaveReactionRules: true,
	}, nil
}
