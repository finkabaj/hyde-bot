package rule

import "errors"

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

var (
	ErrRuleReactionConflict     = errors.New("rule reaction conflict")
	ErrRuleReactionIncompatible = errors.New("rule reaction incompatible")
)
