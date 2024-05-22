package rule

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type ReactAction int

const (
	Delete ReactAction = iota + 1
	Warn
	Ban
	Kick
)

const ReactActionCount = 4

type ReactionRule struct {
	EmojiName  string                        `json:"emojiName,omitempty" validate:"required"`
	EmojiId    string                        `json:"emojiId,omitempty" validate:"omitempty"`
	IsCustom   bool                          `json:"isCustom" validate:"boolean"`
	GuildId    string                        `json:"guildId" validate:"required"`
	RuleAuthor string                        `json:"ruleAuthor" validate:"required"`
	Actions    [ReactActionCount]ReactAction `json:"actions" validate:"dive"`
}

type DeleteReactionRuleQuery struct {
	EmojiId   string `json:"emojiId,omitempty"`
	EmojiName string `json:"emojiName"`
}

var (
	ErrRuleReactionConflict     = errors.New("rule reaction conflict")
	ErrRuleReactionIncompatible = errors.New("rule reaction incompatible")
)

func (a ReactionRule) Compare(b ReactionRule) int {
	if a.EmojiName != b.EmojiName {
		return -1
	}

	if a.EmojiId != b.EmojiId {
		return -1
	}

	if a.IsCustom != b.IsCustom {
		return -1
	}

	if a.GuildId != b.GuildId {
		return -1
	}

	if a.RuleAuthor != b.RuleAuthor {
		return -1
	}

	if len(a.Actions) != len(b.Actions) {
		return -1
	}

	for i, action := range a.Actions {
		if action != b.Actions[i] {
			return -1
		}
	}

	return 0
}

func (a *ReactAction) UnmarshalJSON(data []byte) error {
	var intValue int
	err := json.Unmarshal(data, &intValue)
	if err != nil {
		return err
	}
	*a = ReactAction(intValue)
	return nil
}

func EncodeDeleteReactQuery(queries []DeleteReactionRuleQuery) string {
	values := url.Values{}

	for i, query := range queries {

		if query.EmojiId != "" {
			values.Add("rules["+strconv.Itoa(i)+"][emojiId]", query.EmojiId)
		}

		if query.EmojiName != "" {
			values.Add("rules["+strconv.Itoa(i)+"][emojiName]", query.EmojiName)
		}
	}

	return values.Encode()
}

func DecodeDeleteReactQuery(query string) []DeleteReactionRuleQuery {
	values, _ := url.ParseQuery(query)
	rules := []DeleteReactionRuleQuery{}

	for i := 0; ; i++ {
		emojiId := values.Get("rules[" + strconv.Itoa(i) + "][emojiId]")
		emojiName := values.Get("rules[" + strconv.Itoa(i) + "][emojiName]")

		if emojiId == "" && emojiName == "" {
			break
		}

		rule := DeleteReactionRuleQuery{
			EmojiId:   emojiId,
			EmojiName: emojiName,
		}

		rules = append(rules, rule)
	}

	return rules
}
