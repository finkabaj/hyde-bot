package commands

import (
	"errors"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

type Command struct {
	ApplicationCommand *discordgo.ApplicationCommand
	Handler            func(s *discordgo.Session, i *discordgo.InteractionCreate)
	GuildID            string // GuildID is the ID of the guild where the command is registered. If empty, the command is registered globally.
	IsRegistered       bool   // IsRegistered is a flag that indicates if the command is registered or not. It is set to true when the command is registered.
	RegisteredCommand  *discordgo.ApplicationCommand
}

type CommandManager struct {
	commands            map[string]map[string]*Command // commands[name][guildID] = command
	messageInteractions *commandUtils.MessageInteractions
	rm                  *rules.RuleManager
	lock                sync.RWMutex
}

var cmdManagerInstance *CommandManager

func NewCommandManager(rm *rules.RuleManager, messageInteractions *commandUtils.MessageInteractions) *CommandManager {
	if cmdManagerInstance == nil {
		cmdManagerInstance = &CommandManager{
			commands:            make(map[string]map[string]*Command),
			messageInteractions: messageInteractions,
			rm:                  rm,
		}
	}
	return cmdManagerInstance
}

func (cm *CommandManager) RegisterDefaultCommandsToManager() {
	guildID := ""

	if os.Getenv("ENV") == "development" {
		guildID = os.Getenv("DEV_GUILD_ID")
	}

	cm.RegisterCommandToManager(HelpCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		HelpCommandHandler(s, i, cm)
	}, "")

	cm.RegisterCommandToManager(CreateReactionRuleCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		CreateReactionRuleHandler(s, i, cm.rm)
	}, guildID)

	cm.RegisterCommandToManager(DeleteReactionRuleCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		DeleteReactionRulesCommandHandler(s, i, cm.rm, cm.messageInteractions)
	}, guildID)

	if os.Getenv("ENV") == "development" {
		cm.RegisterCommandToManager(DeleteCommand, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			DeleteCommandHandler(s, i, cm)
		}, guildID)
	}
}

// RegisterCommandToManager registers a command in the CommandManager. If guildID is an empty string, the command will be registered globally.
func (cm *CommandManager) RegisterCommandToManager(cmd *discordgo.ApplicationCommand, handler func(s *discordgo.Session, i *discordgo.InteractionCreate), guildID string) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	command := &Command{
		ApplicationCommand: cmd,
		Handler:            handler,
		GuildID:            guildID,
	}

	if _, ok := cm.commands[cmd.Name]; !ok {
		cm.commands[cmd.Name] = make(map[string]*Command)
	}

	if _, ok := cm.commands[cmd.Name][guildID]; ok {
		logger.Warn(errors.New("command already exists"), map[string]any{"command": cmd.Name, "guildID": guildID})
		return
	}

	cm.commands[cmd.Name][guildID] = command
}

// RegisterDefaultCommands registers all commands in the CommandManager.
func (cm *CommandManager) RegisterDefaultCommands(s *discordgo.Session) (err error) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	for _, cmd := range cm.commands {
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
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, command)

	if err != nil {
		return
	}

	c, ok := cm.commands[cmd.Name][guildID]

	if !ok {
		err = errors.New("command not found")
		logger.Warn(err, map[string]any{"command": cmd.Name, "guildID": guildID})
		return
	}

	c.IsRegistered = true
	c.RegisteredCommand = cmd

	return
}

// DeleteCommand deletes a command on a specific guild by its ID. If guildID is empty, it will delete the command globally.
func (cm *CommandManager) DeleteCommand(s *discordgo.Session, command *discordgo.ApplicationCommand, guildID string, isLocked bool) (err error) {
	if !isLocked {
		cm.lock.Lock()
		defer cm.lock.Unlock()
	}

	err = s.ApplicationCommandDelete(s.State.User.ID, guildID, command.ID)

	if err != nil {
		return
	}

	c, ok := cm.commands[command.Name][guildID]

	if !ok {
		err = errors.New("command not found")
		logger.Warn(err, map[string]any{"command": command.Name, "guildID": guildID})
		return
	}

	c.IsRegistered = false
	c.RegisteredCommand = nil

	return
}

// GetCommandByName returns a command by its name and guildID.
func (cm *CommandManager) GetCommandByName(name string, guildID string) (*Command, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	command, ok := cm.commands[name][guildID]

	if !ok {
		return nil, errors.New("command not found")
	}

	return command, nil
}

// GetGuildCommands returns all commands registered on a specific guild by its ID. If withGlobal is true, it will return global commands as well.
func (cm *CommandManager) GetGuildCommands(guildID string, withGlobal ...bool) (commands []*Command) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	for _, command := range cm.commands {
		if c, ok := command[guildID]; ok {
			commands = append(commands, c)
			continue
		}
		if len(withGlobal) > 0 && withGlobal[0] {
			if c, ok := command[""]; ok {
				commands = append(commands, c)
			}
		}
	}

	return
}

func (cm *CommandManager) AddCommandToGuild(s *discordgo.Session, command *discordgo.ApplicationCommand, guildID string) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.RegisterCommandToManager(command, func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		HelpCommandHandler(s, i, cm)
	}, guildID)

	return cm.RegisterCommand(s, command, guildID)
}

func (cm *CommandManager) DeleteAllCommands(s *discordgo.Session) error {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	for _, commands := range cm.commands {
		for _, c := range commands {
			if !c.IsRegistered {
				continue
			}

			err := cm.DeleteCommand(s, c.RegisteredCommand, c.GuildID, true)

			if err != nil {
				logger.Error(err, map[string]any{"details": "Error removing command"})
				return err
			}
			logger.Info("Removed command: " + c.ApplicationCommand.Name)
		}
	}

	return nil
}
