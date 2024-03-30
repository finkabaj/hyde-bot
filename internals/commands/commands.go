package commands

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(s *discordgo.Session, i *discordgo.InteractionCreate)
	GuildID            string // GuildID is the ID of the guild where the command is registered. If empty, the command is registered globally.
	IsRegistered       bool   // IsRegistered is a flag that indicates if the command is registered or not. It is set to true when the command is registered.
	RegisteredCommand  *discordgo.ApplicationCommand
}

type CommandManager struct {
	Commands []*Command
}

var cmdManagerInstance *CommandManager

func NewCommandManager() *CommandManager {
	if cmdManagerInstance == nil {
		cmdManagerInstance = &CommandManager{
			Commands: make([]*Command, 0),
		}
	}
	return cmdManagerInstance
}

func (cm *CommandManager) RegisterDefaultCommandsToManager() {
	var guildID string

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
	cm.Commands = append(cm.Commands, command)
}

// RegisterDefaultCommands registers all commands in the CommandManager.
func (cm *CommandManager) RegisterDefaultCommands(s *discordgo.Session) error {
	for _, cmd := range cm.Commands {
		registeredCmd, err := s.ApplicationCommandCreate(s.State.User.ID, cmd.GuildID, cmd.ApplicationCommand)
		if err != nil {
			return err
		}
		cmd.IsRegistered = true

		cmd.RegisteredCommand = registeredCmd
	}
	return nil
}

// RegisterCommand registers a command on a specific guild by its ID. If guildID is empty, it will register the command globally.
func (cm *CommandManager) RegisterCommand(s *discordgo.Session, command *discordgo.ApplicationCommand, guildID string) error {
	cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, command)
	if err != nil {
		return err
	}

	for _, c := range cm.Commands {
		if c.ApplicationCommand.Name == cmd.Name {
			c.RegisteredCommand = cmd
			c.IsRegistered = true
		}
	}

	return nil
}

// DeleteCommand deletes a command on a specific guild by its ID. If guildID is empty, it will delete the command globally.
func (cm *CommandManager) DeleteCommand(s *discordgo.Session, command *discordgo.ApplicationCommand, guildID string) error {
	err := s.ApplicationCommandDelete(s.State.User.ID, guildID, command.ID)
	if err != nil {
		return err
	}
	return nil
}
