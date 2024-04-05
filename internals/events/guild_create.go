package events

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
)

func HandleGuildCreate(s *discordgo.Session, event interface{}) {
	typedEvent, ok := event.(*discordgo.GuildCreate)

	if !ok {
		err := errors.New("Incorect type in HandleGuildCreate")
		logger.Warn(err)
		return
	}

	url := "http://" + os.Getenv("API_HOST") + ":" + os.Getenv("API_PORT") + "/guild"
	jsonBody := []byte(fmt.Sprintf(`{"guildId": %s}`, typedEvent.ID))
	bodyReader := bytes.NewReader(jsonBody)

	res, err := http.Post(url, "application-json", bodyReader)

	if err != nil {
		logger.Error(err, logger.LogFields{"message": "error while sending post request on guild create"})
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(res.Body)

		if err != nil {
			logger.Error(err, logger.LogFields{"message": "error while reading response body"})
		}

		fmt.Println(string(body))
	}
}
