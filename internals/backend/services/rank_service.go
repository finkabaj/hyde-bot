package services

import "github.com/finkabaj/hyde-bot/internals/ranks"

type IRankService interface {
	CreateRanks(r ranks.Ranks) (ranks.Ranks, error)
	GetRanks(gId string) (ranks.Ranks, error)
	DeleteRank(gId string, rId string) error
	DeleteRanks(gId string) error
}
