package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
)

func HandleGuildCreate(s *discordgo.Session, event interface{}) {
	typedEvent, ok := event.(*discordgo.GuildCreate)

	if !ok {
		logger.Warn(errors.New("Incorect type in HandleGuildCreate"))
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

	var result interface{}

	common.UnmarshalBody(body, &result)

	if _, ok := result.(guild.Guild); !ok {
		if err, ok := result.(common.ErrorResponse); ok {
			logger.Error(errors.New(err.Message), logger.ToLogFields(err.ValidationErrors))
		}

		logger.Error(errors.New("Error while unmarshalling guild response"))
	}
}
