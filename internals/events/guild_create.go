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
	b, err := io.ReadAll(body)

	if err != nil {
		logger.Fatal(err, logger.LogFields{"MESSAGE": "The bot cannot continue to work correctly", "AT": "guild_create"})
	}

	var result guild.Guild

	if err := common.UnmarshalBodyBytes(b, &result); err != nil {
		logger.Error(errors.New("Error while unmarshaling guild create"))
	}

	if result.GuildId != "" && result.OwnerId != "" {
		return
	}

	var errRes common.ErrorResponse

	if err := common.UnmarshalBodyBytes(b, &errRes); err != nil {
		logger.Error(errors.New("Error while unmarshaling guild create error"))
	}

	logger.Error(errors.New(errRes.Error), logger.ToLogFields(errRes.ValidationErrors))
}
