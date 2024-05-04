package rules

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
)

type Rules struct {
	ReactionRules     []rule.ReactionRule `json:"reactionRules"`
	HaveReactionRules bool                `json:"haveReactionRules"`
}

type RuleManager map[string]Rules

var ruleManager RuleManager

func NewRuleManager() *RuleManager {
	if ruleManager == nil {
		ruleManager = make(RuleManager)
	}
	return &ruleManager
}

func (rm *RuleManager) AddRules(guildId string, rules Rules) {
	(*rm)[guildId] = rules
}

func (rm *RuleManager) GetRules(guildId string) (Rules, error) {
	rules, ok := (*rm)[guildId]

	if !ok {
		return Rules{}, errors.New("rules not found")
	}

	return rules, nil
}

func (rm *RuleManager) GetReactionRules(guildId string) ([]rule.ReactionRule, error) {
	rules, err := rm.GetRules(guildId)

	if err != nil {
		return nil, fmt.Errorf("error getting reaction rules: %w", err)
	}

	if !rules.HaveReactionRules {
		return nil, errors.New("reaction rules not found")
	}

	return rules.ReactionRules, nil
}

func (rm *RuleManager) FetchReactionRules(guildId string) ([]rule.ReactionRule, error) {
	reactionRulesUrl := common.GetApiUrl(os.Getenv("API_HOST"), os.Getenv("API_PORT"), "/rules/reaction/"+guildId)
	res, err := http.Get(reactionRulesUrl)

	if err != nil {
		return nil, err
	}

	body := res.Body
	defer body.Close()
	b, err := io.ReadAll(body)

	if err != nil {
		logger.Fatal(err, logger.LogFields{"MESSAGE": "The bot cannot continue to work correctly", "AT": "guild_create"})
	}

	var reactionRules []rule.ReactionRule

	err = common.UnmarshalBodyBytes(b, &reactionRules)

	if err != nil {
		var errRes common.ErrorResponse

		if err = common.UnmarshalBodyBytes(b, &errRes); err != nil {
			logger.Error(errors.New("error while unmarshaling reaction rules error"))
			return nil, err
		}

		err = errors.New(errRes.Error)

		return nil, err
	}

	return reactionRules, nil
}
