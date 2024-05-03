package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/commands"
	"github.com/finkabaj/hyde-bot/internals/events"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	"github.com/joho/godotenv"
)

var s *discordgo.Session
var err error
var fs *os.File

var (
	RemoveCommands   = flag.Bool("rmcmd", false, "Remove all commands on shutdown")
	RegisterCommands = flag.Bool("regcmd", true, "Register all commands on startup")
)

func init() {
	flag.Parse()
	err = godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	s, err = discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	rm := rules.NewRuleManager()

	evtManager := events.NewEventManager(rm)

	evtManager.RegisterDefaultEvents()

	s.AddHandler(func(s *discordgo.Session, event interface{}) {
		evtManager.HandleEvent(s, event)
	})

	cmdManager := commands.NewCommandManager()

	cmdManager.RegisterDefaultCommandsToManager()

	fs, err = os.OpenFile("log/logs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		log.Fatal("Error creating new log file", err)
	}

	logger.NewLogger(fs)

	if err != nil {
		logger.Fatal(err, logger.LogFields{"message": "Error creating a new log file"})
	}

	s.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
		logger.Info("Bot is up and running!")
	})

	err = s.Open()

	if err != nil {
		logger.Fatal(err, logger.LogFields{"message": "Error opening a connection to Discord"})
	}

	if *RegisterCommands {
		err = cmdManager.RegisterDefaultCommands(s)

		if err != nil {
			logger.Fatal(err, logger.LogFields{"message": "Error registering commands"})
		}

		logger.Info("Commands registered")
	} else {
		logger.Info("Skipping command registration")
	}

	defer s.Close()
	defer fs.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	if *RemoveCommands {
		logger.Info("Removing commands")

		for _, c := range cmdManager.Commands {
			for _, command := range c {
				if !command.IsRegistered {
					continue
				}

				err = cmdManager.DeleteCommand(s, command.RegisteredCommand, command.GuildID)

				if err != nil {
					logger.Error(err, logger.LogFields{"message": "Error removing command"})
				}
				logger.Info("Removed command: " + command.ApplicationCommand.Name)
			}
		}
	}

	logger.Info("Shutting down")
}
