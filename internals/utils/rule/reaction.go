package rule

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
)

type ReactAction int

const (
	Delete ReactAction = iota
	Warn
	Ban
	Kick
)

type ReactionRule struct {
	EmojiName  string      `json:"emojiName,omitempty"`
	EmojiId    string      `json:"emojiId,omitempty"`
	GuildId    string      `json:"guildId" validate:"required,len=19"`
	RuleAuthor string      `json:"ruleAuthor" validate:"required,min=17,max=18"`
	Action     ReactAction `json:"action" validate:"number"`
}

type DeleteReactionRuleQuery struct {
	EmojiId   string `json:"emojiId,omitempty"`
	EmojiName string `json:"emojiName,omitempty"`
}

var (
	ErrRuleReactionConflict     = errors.New("rule reaction conflict")
	ErrRuleReactionIncompatible = errors.New("rule reaction incompatible")
)

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
