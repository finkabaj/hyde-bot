package commands

import (
	"errors"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(s *discordgo.Session, i *discordgo.InteractionCreate)
	GuildID            string // GuildID is the ID of the guild where the command is registered. If empty, the command is registered globally.
	IsRegistered       bool   // IsRegistered is a flag that indicates if the command is registered or not. It is set to true when the command is registered.
	RegisteredCommand  *discordgo.ApplicationCommand
}

type CommandManager struct {
	Commands map[string]map[string]*Command // commands[name][guildID] = command
}

var cmdManagerInstance *CommandManager

func NewCommandManager() *CommandManager {
	if cmdManagerInstance == nil {
		cmdManagerInstance = &CommandManager{
			Commands: make(map[string]map[string]*Command),
		}
	}
	return cmdManagerInstance
}

func (cm *CommandManager) RegisterDefaultCommandsToManager() {
	var guildID string = ""

	if os.Getenv("ENV") == "development" {
		guildID = os.Getenv("DEV_GUILD_ID")
	}

	cm.RegisterCommandToManager(HelpCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		HelpCommandHandler(s, i)
	}, guildID)

	cm.RegisterCommandToManager(DeleteCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		DeleteCommandHandler(s, i)
	}, guildID)
}

// RegisterCommandToManager registers a command in the CommandManager. If guildID is an empty string, the command will be registered globally.
func (cm *CommandManager) RegisterCommandToManager(cmd *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate), guildID string) {
	command := &Command{
		ApplicationCommand: cmd,
		Handler:            handler,
		GuildID:            guildID,
	}

	if _, ok := cm.Commands[cmd.Name]; !ok {
		cm.Commands[cmd.Name] = make(map[string]*Command)
	}

	if _, ok := cm.Commands[cmd.Name][guildID]; ok {
		logger.Warn(errors.New("Command already exists"), logger.LogFields{"command": cmd.Name, "guildID": guildID})
		return
	}

	cm.Commands[cmd.Name][guildID] = command
}

// RegisterDefaultCommands registers all commands in the CommandManager.
func (cm *CommandManager) RegisterDefaultCommands(s *discordgo.Session) (err error) {

	for _, cmd := range cm.Commands {
		for _, c := range cmd {
			c.RegisteredCommand, err = s.ApplicationCommandCreate(s.State.User.ID, c.GuildID, c.ApplicationCommand)

			if err != nil {
				return
			}

			c.IsRegistered = true
		}
	}
	return
}

// RegisterCommand registers a command on a specific guild by its ID. If guildID is empty, it will register the command globally.
func (cm *CommandManager) RegisterCommand(s *discordgo.Session, command *discordgo.ApplicationCommand, guildID string) (err error) {
	cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, command)

	if err != nil {
		return
	}

	c, ok := cm.Commands[cmd.Name][guildID]

	if !ok {
		err = errors.New("Command not found")
		logger.Warn(err, logger.LogFields{"command": cmd.Name, "guildID": guildID})
		return
	}

	c.IsRegistered = true
	c.RegisteredCommand = cmd

	return
}

// DeleteCommand deletes a command on a specific guild by its ID. If guildID is empty, it will delete the command globally.
func (cm *CommandManager) DeleteCommand(s *discordgo.Session, command *discordgo.ApplicationCommand, guildID string) (err error) {
	err = s.ApplicationCommandDelete(s.State.User.ID, guildID, command.ID)

	if err != nil {
		return
	}

	c, ok := cm.Commands[command.Name][guildID]

	if !ok {
		err = errors.New("Command not found")
		logger.Warn(err, logger.LogFields{"command": command.Name, "guildID": guildID})
		return
	}

	c.IsRegistered = false
	c.RegisteredCommand = nil

	return
}
