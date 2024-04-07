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
		err := errors.New("Incorect type in HandleGuildCreate")
		logger.Warn(err)
		return
	}

	info := &guild.GuildCreate{
		GuildId: typedEvent.ID,
		OwnerId: typedEvent.OwnerID,
	}

	jsonInfo, err := json.Marshal(&info)

	if err != nil {
		logger.Error(err, logger.LogFields{"message": "error while marshalling guild info"})
	}

	url := "http://" + os.Getenv("API_HOST") + ":" + os.Getenv("API_PORT") + "/guild"

	bodyReader := bytes.NewReader(jsonInfo)
	res, err := http.Post(url, "application/json", bodyReader)

	if err != nil {
		logger.Error(err, logger.LogFields{"message": "error while sending post request on guild create"})
		return
	}

	body := res.Body
	defer body.Close()

	var result *guild.GuildCreate

	common.UnmarshalBody(body, &result)
}
