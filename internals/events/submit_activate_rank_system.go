package events

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/ranks"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
)

func HandleSubmitActivateRankSystem(rankManager *ranks.RankManager, ruleManager *rules.RuleManager) EventHandler {
	return func(s *discordgo.Session, event interface{}) {
		data, i, err := commandUtils.GetDataFromModalSubmit(event)

		if err != nil {
			logger.Error(fmt.Errorf("error at HandleSumbitModalReaction: %w", err))
			return
		}

		if !strings.HasPrefix(string(data.CustomID), "activate_role_system") {
			return
		}

		ti1 := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)

		if ti1.CustomID != "role_system" {
			logger.Error(fmt.Errorf("expected role_system custom id, got %s instead", ti1.CustomID))
			return
		}

		ti2 := data.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)

		if ti2.CustomID != "xp_system" {
			logger.Error(fmt.Errorf("expected xp_system custom id, got %s instead", ti2.CustomID))
			return
		}

		roles, err := s.GuildRoles(i.GuildID)
		if err != nil {
			logger.Error(err)
			return
		}

		r, err := parseActivateRoleSystemInput(ti1.Value, ti2.Value, roles)

		if err != nil {
			commandUtils.SendDefaultResponse(s, i, err.Error())
			logger.Error(err, map[string]any{"details": "error parsing activate role system input"})
			return
		}

		if len(r) == 0 {
			commandUtils.SendDefaultResponse(s, i, "No valid roles found in the input")
			return
		}

		commandUtils.SendDefaultResponse(s, i, "Role system activated")
	}
}

func parseActivateRoleSystemInput(ids, xps string, roles []*discordgo.Role) ([]ranks.Rank, error) {
	if len(roles) == 0 {
		return nil, fmt.Errorf("no roles found in the guild")
	}
	seenRoles := make(map[string]bool)

	idsSplited := strings.Split(ids, " ")
	xpsSplited := strings.Split(xps, " ")
	var xpStep uint16
	if len(xpsSplited) == 1 {
		n, err := strconv.Atoi(xpsSplited[0])

		if err != nil {
			return nil, fmt.Errorf("xp should be a number")
		}

		xpStep = uint16(n)
	}

	if len(idsSplited) != len(xpsSplited) && len(xpsSplited) != 1 {
		return nil, fmt.Errorf("ids and xps should have the same length or xps should have length 1")
	}

	rs := make([]ranks.Rank, 0, len(ids))

	for i, id := range idsSplited {
		for _, role := range roles {
			if role.ID == id {
				if _, ok := seenRoles[id]; ok {
					return nil, fmt.Errorf("duplicate role id found: %s", id)
				}
				seenRoles[id] = true
				var xp uint16
				if len(xpsSplited) != 1 {
					n, err := strconv.Atoi(xpsSplited[i])
					if err != nil {
						return nil, fmt.Errorf("xp should be a number")
					}
					xp = uint16(n)
				} else {
					xp += xpStep
				}
				// TODO: make uuid
				rs = append(rs, ranks.Rank{
					Level: uint8(i + 1),
					XP:    xp,
					ID:    "TODO: make uuid",
					Role:  role,
				})
			}
		}
	}

	return rs, nil
}
