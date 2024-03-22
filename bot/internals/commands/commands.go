package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type CommandManager struct {
	Commands []*Command
	RegisteredCommands []*discordgo.ApplicationCommand
}

func NewCommandManager() *CommandManager {
	return &CommandManager{
		Commands: make([]*Command, 0),
		RegisteredCommands: make([]*discordgo.ApplicationCommand, 0),
	}
}

// RegisterCommandToManager registers a command in the CommandManager.
func (cm *CommandManager) RegisterCommandToManager(cmd *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	command := &Command{
		ApplicationCommand: cmd,
		Handler:            handler,
	}
	cm.Commands = append(cm.Commands, command)
}

// RegisterDefaultCommands registers all commands in the CommandManager on the bot globally.
func (cm *CommandManager) RegisterDefaultCommands(s *discordgo.Session) error {
	for _, cmd := range cm.Commands {
		registeredCmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd.ApplicationCommand)
		if err != nil {
			return err
		}
		cm.RegisteredCommands = append(cm.RegisteredCommands, registeredCmd)
	}
	return nil
}

// RegisterCommand registers a command on a specific guild by its ID. If guildID is empty, it will register the command globally.
func (cm *CommandManager) RegisterCommand(s *discordgo.Session, command *discordgo.ApplicationCommand, guildID string) error {
	cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, command)
	if err != nil {
		return err
	}
	cm.RegisteredCommands = append(cm.RegisteredCommands, cmd)
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

