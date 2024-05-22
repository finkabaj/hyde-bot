package rules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
)

var (
	ErrIntersectingRules = errors.New("reaction rules already exist")
	ErrRulesNotFound     = errors.New("rules not found")
)

type Rules struct {
	ReactionRules     []rule.ReactionRule `json:"reactionRules"`
	HaveReactionRules bool                `json:"haveReactionRules"`
}

type RulesDeleteDto struct {
	EmojiName string
	EmojiId   string
}

type RuleManager struct {
	rm     map[string]Rules
	client *http.Client
	lock   sync.RWMutex
}

var ruleManager *RuleManager

func NewRuleManager(client *http.Client) *RuleManager {
	if ruleManager == nil {
		ruleManager = &RuleManager{
			rm:     make(map[string]Rules),
			client: client,
		}
	}
	return ruleManager
}

func (rm *RuleManager) AddRules(guildId string, rules Rules) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	rm.rm[guildId] = rules
}

func (rm *RuleManager) DeleteReactionRules(guildID string, reactionRules []rule.ReactionRule) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	rules := rm.rm[guildID]
	for _, rule := range reactionRules {
		for i, r := range rules.ReactionRules {
			if r.EmojiName == rule.EmojiName && r.EmojiId == rule.EmojiId {
				rules.ReactionRules = append(rules.ReactionRules[:i], rules.ReactionRules[i+1:]...)
			}
		}
	}
}

func (rm *RuleManager) AddReactionRules(guildId string, reactionRules []rule.ReactionRule) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	rules := rm.rm[guildId]
	rules.ReactionRules = append(rules.ReactionRules, reactionRules...)
	rules.HaveReactionRules = true
	rm.rm[guildId] = rules
}

func (rm *RuleManager) GetRules(guildId string, locked bool) (Rules, error) {
	if !locked {
		rm.lock.RLock()
		defer rm.lock.RUnlock()
	}

	rules, ok := rm.rm[guildId]

	if !ok {
		return Rules{}, errors.New("rules not found")
	}

	return rules, nil
}

func (rm *RuleManager) GetReactionRules(guildId string, locked bool) ([]rule.ReactionRule, error) {
	if !locked {
		rm.lock.RLock()
		defer rm.lock.RUnlock()
	}
	rules, err := rm.GetRules(guildId, true)

	if err != nil {
		return nil, fmt.Errorf("error getting reaction rules: %w", err)
	}

	if !rules.HaveReactionRules {
		return nil, ErrRulesNotFound
	}

	return rules.ReactionRules, nil
}

func (rm *RuleManager) FetchReactionRules(guildId string) ([]rule.ReactionRule, error) {
	reactionRulesUrl := common.GetApiUrl(os.Getenv("API_HOST"), os.Getenv("API_PORT"), "/rules/reaction/"+guildId)
	res, err := rm.client.Get(reactionRulesUrl)

	if err != nil {
		return nil, err
	}

	body := res.Body
	defer body.Close()
	b, err := io.ReadAll(body)

	if err != nil {
		logger.Fatal(err, map[string]any{"details": "The bot cannot continue to work correctly", "at": "guild_create"})
	}

	if res.StatusCode != http.StatusOK {
		var errRes common.ErrorResponse

		if err = common.UnmarshalBodyBytes(b, &errRes); err != nil {
			logger.Error(errors.New("error while unmarshaling reaction rules error"))
			return nil, err
		}

		logger.Debug("Error response", map[string]any{"status": res.StatusCode, "error": errRes.Error, "validationErrors": errRes.ValidationErrors, "message": errRes.Message})
		return nil, errors.New(errRes.Error)
	}

	var reactionRules []rule.ReactionRule

	if err = common.UnmarshalBodyBytes(b, &reactionRules); err != nil {
		logger.Error(errors.New("error while unmarshaling reaction rules"))
		return nil, err
	}

	return reactionRules, nil
}

func (rm *RuleManager) PostReactionRules(guildId string, reactionRules []rule.ReactionRule) ([]rule.ReactionRule, error) {
	existingRules, err := rm.GetReactionRules(guildId, false)

	if err != nil && !errors.Is(err, ErrRulesNotFound) {
		return nil, fmt.Errorf("error posting reaction rules: %w", err)
	}

	if common.HaveIntersection(existingRules, reactionRules) {
		return nil, ErrIntersectingRules
	}

	b, err := json.Marshal(reactionRules)

	if err != nil {
		return nil, fmt.Errorf("error marshaling reaction rules: %w", err)
	}

	bb := bytes.NewReader(b)

	rRulesApiUrl := common.GetApiUrl(os.Getenv("API_HOST"), os.Getenv("API_PORT"), "/rules/reaction")
	res, err := rm.client.Post(rRulesApiUrl, "application/json", bb)

	if err != nil {
		return nil, fmt.Errorf("error posting reaction rules: %w", err)
	}

	body := res.Body
	defer body.Close()
	b, err = io.ReadAll(body)

	if err != nil {
		logger.Fatal(err, map[string]any{"details": "The bot cannot continue to work correctly", "at": "guild_create"})
	}

	if res.StatusCode != http.StatusCreated {
		var errRes common.ErrorResponse

		if err = common.UnmarshalBodyBytes(b, &errRes); err != nil {
			logger.Error(errors.New("error posting unmarshaling reaction rules error"))
			return nil, err
		}

		logger.Debug("Error response", map[string]any{"status": res.StatusCode, "error": errRes.Error, "validationErrors": errRes.ValidationErrors, "message": errRes.Message})

		return nil, errors.New(errRes.Error)
	}

	var rRules []rule.ReactionRule

	if err = common.UnmarshalBodyBytes(b, &rRules); err != nil {
		logger.Error(errors.New("error posting unmarshaling reaction rules"))
		return nil, err
	}

	return rRules, nil
}

func (rm *RuleManager) DeleteReactionRulesApi(guildId string, deleteDto []RulesDeleteDto) error {
	rRules, err := rm.GetReactionRules(guildId, false)

	if err != nil {
		return fmt.Errorf("error deleting reaction rules: %w", err)
	}

	query := make([]rule.DeleteReactionRuleQuery, len(deleteDto))

	for _, dto := range deleteDto {
		for _, rRule := range rRules {
			if rRule.EmojiName == dto.EmojiName && rRule.EmojiId == dto.EmojiId {
				query = append(query, rule.DeleteReactionRuleQuery{
					EmojiName: dto.EmojiName,
					EmojiId:   dto.EmojiId,
				})
			}
		}
	}

	if len(query) != len(deleteDto) {
		return errors.New("some rules are not found")
	}

	queryString := rule.EncodeDeleteReactQuery(query)

	rRulesApiUrl := common.GetApiUrl(os.Getenv("API_HOST"), os.Getenv("API_PORT"), "/rules/reaction/"+guildId+"?"+queryString)

	req, err := http.NewRequest(http.MethodDelete, rRulesApiUrl, nil)

	if err != nil {
		return fmt.Errorf("error deleting reaction rules: %w", err)
	}

	res, err := rm.client.Do(req)

	if err != nil {
		return fmt.Errorf("error deleting reaction rules: %w", err)
	}

	body := res.Body
	defer body.Close()
	b, err := io.ReadAll(body)

	if err != nil {
		logger.Error(err, map[string]any{"details": "error while reading response body"})
		return fmt.Errorf("error deleting reaction rules: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		var errRes common.ErrorResponse

		if err = common.UnmarshalBodyBytes(b, &errRes); err != nil {
			logger.Error(errors.New("error posting unmarshaling reaction rules error"))
			return fmt.Errorf("error deleting reaction rules: %w", err)
		}

		logger.Debug("Error response", map[string]any{"status": res.StatusCode, "error": errRes.Error, "validationErrors": errRes.ValidationErrors, "message": errRes.Message})

		return errors.New(errRes.Error)
	}

	return nil
}
