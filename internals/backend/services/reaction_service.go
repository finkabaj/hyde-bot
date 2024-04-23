package services

import (
	"github.com/finkabaj/hyde-bot/internals/db"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
)

type IReactionService interface {
	CreateReactionRules(rules []rule.ReactionRule) ([]rule.ReactionRule, error)
	GetReactionRules(gId string) ([]rule.ReactionRule, error)
	DeleteReactionRules(query []rule.DeleteReactionRuleQuery, gId string) error
}

type ReactionService struct {
	logger       logger.ILogger
	database     db.Database
	guildService IGuildService
}

var reactionService *ReactionService

func NewReactionService(l logger.ILogger, d db.Database, g IGuildService) *ReactionService {
	if reactionService == nil {
		reactionService = &ReactionService{
			logger:       l,
			database:     d,
			guildService: g,
		}
	}
	return reactionService
}

func (rs *ReactionService) CreateReactionRules(rules []rule.ReactionRule) ([]rule.ReactionRule, error) {
	if len(rules) < 1 {
		return []rule.ReactionRule{}, common.ErrBadRequest
	}

	gId := rules[0].GuildId

	_, err := rs.guildService.GetGuild(gId)

	if err != nil {
		return []rule.ReactionRule{}, err
	}

	foundRules, _ := rs.GetReactionRules(gId)

	for _, v := range rules {
		if v.GuildId != gId {
			return []rule.ReactionRule{}, common.ErrBadRequest
		}

		if v.EmojiId != "" && common.ContainsFieldValue(foundRules, "EmojiId", v.EmojiId) {
			return []rule.ReactionRule{}, rule.ErrRuleReactionConflict
		} else if v.EmojiName != "" && common.ContainsFieldValue(foundRules, "EmojiName", v.EmojiName) {
			return []rule.ReactionRule{}, rule.ErrRuleReactionConflict
		}

		if v.EmojiId == "" && v.EmojiName == "" {
			return []rule.ReactionRule{}, common.ErrBadRequest
		}

		if v.EmojiName != "" && v.EmojiId != "" {
			return []rule.ReactionRule{}, rule.ErrRuleReactionIncompatible
		}

		actionsLen := len(v.Actions)

		if actionsLen == 0 || len(common.RemoveDuplicates(v.Actions)) != actionsLen {
			return []rule.ReactionRule{}, common.ErrBadRequest
		}
	}

	createdRules, err := rs.database.CreateReactionRules(rules)

	if err != nil {
		return []rule.ReactionRule{}, err
	}

	return createdRules, nil
}

func (rs *ReactionService) GetReactionRules(gId string) ([]rule.ReactionRule, error) {
	_, err := rs.guildService.GetGuild(gId)

	if err != nil {
		return []rule.ReactionRule{}, err
	}

	rRules, err := rs.database.ReadReactionRules(gId)

	if err != nil {
		return []rule.ReactionRule{}, err
	}

	return rRules, err
}

func (rs *ReactionService) DeleteReactionRules(query []rule.DeleteReactionRuleQuery, gId string) error {
	if len(query) < 1 {
		return common.ErrBadRequest
	}

	_, err := rs.guildService.GetGuild(gId)

	if err != nil {
		return err
	}

	for _, r := range query {
		if r.EmojiId == "" && r.EmojiName == "" {
			return common.ErrBadRequest
		}

		if r.EmojiId != "" && r.EmojiName != "" {
			return rule.ErrRuleReactionIncompatible
		}
	}

	err = rs.database.DeleteReactionRules(query, gId)

	if err != nil {
		return err
	}

	return nil
}
