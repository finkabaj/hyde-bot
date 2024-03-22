package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/PapicBorovoi/hyde-bot/internals/commands"
	"github.com/PapicBorovoi/hyde-bot/internals/events"
	"github.com/PapicBorovoi/hyde-bot/internals/logger"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var s *discordgo.Session
var err error
var fs *os.File

var ( 
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands on shutdown")
	RegisterCommands = flag.Bool("regcmd", true, "Register all commands on startup")
)

func init() {
	flag.Parse()
	err = godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}


	s, err = discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		fmt.Println("Error creating a new Discord session: ", err)
		os.Exit(1)
	}
}


func main() {
    evtManager := events.NewEventManager()

    evtManager.RegisterEventHandler(
      "MessageReactionAdd", 
      func (s *discordgo.Session, event discordgo.MessageReactionAdd)  {
        events.HandleDeleteReaction(s, event)
      },
      "",
    )


    s.AddHandler(func (s *discordgo.Session, event interface{})  {
      evtManager.HandleEvent(s, event)
    })

	cmdManager := commands.NewCommandManager()

	cmdManager.RegisterCommandToManager(commands.HelpCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commands.HelpCommandHandler(s, i, cmdManager)
	})

	cmdManager.RegisterCommandToManager(commands.DeleteCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commands.DeleteCommandHandler(s, i, cmdManager)
	})

	s.AddHandler(func (s *discordgo.Session, i* discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			for _, command := range cmdManager.Commands {
				if command.ApplicationCommand.Name == i.ApplicationCommandData().Name {
					command.Handler(s, i)
					break
				}
			}
		}
	})

	fs, err = os.OpenFile("log/logs.log" , os.O_APPEND | os.O_CREATE | os.O_RDWR, 0644)

	if err != nil {
		fmt.Println("Error creating a new log file: ", err)
		os.Exit(1)
	}

	log := logger.Init(fs)

	if err != nil {
		log.Fatal(err, "Error creating a new log file")
		os.Exit(1)
	}

	s.AddHandler(func (s *discordgo.Session, m *discordgo.Ready)  {
		log.Info("Bot is up and running!")
	})

	err = s.Open()

	if err != nil {
		log.Fatal(err, "Error opening a connection to Discord")
		os.Exit(1)
	}

	if *RegisterCommands {
		err = cmdManager.RegisterDefaultCommands(s)

		if err != nil {
			log.Fatal(err, "Error registering commands")
			os.Exit(1)
		}

		log.Info("Commands registered")
	} else {
		log.Info("Skipping command registration")
	}

	defer s.Close()
	defer fs.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	if *RemoveCommands {
		fmt.Println("\nRemoving commands")
		log.Info("Removing commands")

		for _, command := range cmdManager.RegisteredCommands {
			if (command == nil) {
				continue
			}

			err = cmdManager.DeleteCommand(s, command, "")

			if err != nil {
				log.Error(err, "Error removing command")
			}
			fmt.Println("Removed command: ", command.Name)
		}
	}

	log.Info("Shutting down")
}
