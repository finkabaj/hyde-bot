package services

import "github.com/finkabaj/hyde-bot/internals/utils/rule"

type IReactionService interface {
	CreateReactionRules(rules *[]rule.ReactionRule) (*[]rule.ReactionRule, error)
	GetReactionRules(gId string) (*[]rule.ReactionRule, error)
	DeleteReactionRules(rules *[]rule.ReactionRule) error
}
