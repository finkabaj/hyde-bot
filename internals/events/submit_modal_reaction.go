package events

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/enescakir/emoji"
	"github.com/finkabaj/hyde-bot/internals/logger"
	"github.com/finkabaj/hyde-bot/internals/rules"
	commandUtils "github.com/finkabaj/hyde-bot/internals/utils/command"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
)

func HandleSumbitModalReaction(rm *rules.RuleManager) EventHandler {
	return func(s *discordgo.Session, event any) {
		data, i, err := commandUtils.GetDataFromModalSubmit(event)

		if err != nil {
			logger.Error(fmt.Errorf("error at HandleSumbitModalReaction: %w", err))
			return
		}

		if !strings.HasPrefix(string(data.CustomID), "emoji_ban") {
			return
		}

		text := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

		emojies, err := s.GuildEmojis(i.GuildID)
		if err != nil {
			logger.Error(err)
			return
		}

		r := parseModalReactionInput(text, i.Member.User.ID, i.GuildID, emojies)

		if len(r) == 0 {
			commandUtils.SendDefaultResponse(s, i, "No valid emojies found in the input")

			return
		}

		rRules, err := rm.PostReactionRules(i.GuildID, r)

		if err != nil {
			if errors.Is(err, rules.ErrIntersectingRules) {
				commandUtils.SendDefaultResponse(s, i, "Reaction rules already exist")

				return
			}

			commandUtils.SendDefaultResponse(s, i, "Failed to post reaction rules")

			logger.Error(err)

			return
		}

		rm.AddReactionRules(i.GuildID, rRules)

		commandUtils.SendDefaultResponse(s, i, "Modal submitted successfully!")
	}
}

func parseModalReactionInput(text string, ruleAuthor string, guildId string, emojies []*discordgo.Emoji) []rule.ReactionRule {
	if len(emojies) == 0 {
		return nil
	}

	textSplited := strings.Split(text, " ")

	if len(textSplited) == 0 {
		return nil
	}

	result := make([]rule.ReactionRule, 0, len(textSplited))

	const rac = rule.ReactActionCount

	for _, v := range textSplited {
		if emoji.Exist(v) {
			result = append(result, rule.ReactionRule{
				GuildId:    guildId,
				RuleAuthor: ruleAuthor,
				EmojiName:  emoji.Parse(v),
				IsCustom:   false,
				Actions:    [rac]rule.ReactAction{rule.Delete},
			})
		} else if i := slices.IndexFunc(emojies, func(e *discordgo.Emoji) bool {
			return v != "" && e.ID == v
		}); i != -1 {
			result = append(result, rule.ReactionRule{
				GuildId:    guildId,
				RuleAuthor: ruleAuthor,
				EmojiId:    v,
				EmojiName:  emojies[i].Name,
				IsCustom:   true,
				Actions:    [rac]rule.ReactAction{rule.Delete},
			})
		} else if strings.HasPrefix(v, ":") && strings.HasSuffix(v, ":") {
			eId := ""
			eParsed := strings.Trim(v, ":")
			i := slices.IndexFunc(emojies, func(e *discordgo.Emoji) bool {
				return e.Name == eParsed
			})
			if i != -1 {
				eId = emojies[i].ID
			} else {
				return nil
			}

			result = append(result, rule.ReactionRule{
				GuildId:    guildId,
				RuleAuthor: ruleAuthor,
				EmojiName:  eParsed,
				EmojiId:    eId,
				IsCustom:   true,
				Actions:    [rac]rule.ReactAction{rule.Delete},
			})
		} else {
			var emojiSequence string
			b := []byte(v)
			for i := 0; i < len(b); i++ {
				r, size := utf8.DecodeRune(b[i:])

				// Check if the rune is not an emoji
				if (r < 0x1F600 || r > 0x1F64F) && (r < 0x1F300 || r > 0x1F5FF) &&
					(r < 0x1F680 || r > 0x1F6FF) && (r < 0x2600 || r > 0x26FF) &&
					(r < 0x2700 || r > 0x27BF) && (r < 0xFE00 || r > 0xFE0F) &&
					(r < 0x1F900 || r > 0x1F9FF) && (r < 0x1F1E6 || r > 0x1F1FF) && r != 0x200D {
					if emojiSequence != "" {
						result = append(result, rule.ReactionRule{
							GuildId:    guildId,
							RuleAuthor: ruleAuthor,
							EmojiName:  emojiSequence,
							IsCustom:   false,
							Actions:    [rac]rule.ReactAction{rule.Delete},
						})
						emojiSequence = ""
					}
				} else {
					emojiSequence += string(r)
				}
				i += size - 1
			}
			if emojiSequence != "" {
				result = append(result, rule.ReactionRule{
					GuildId:    guildId,
					RuleAuthor: ruleAuthor,
					EmojiName:  v,
					IsCustom:   false,
					Actions:    [rac]rule.ReactAction{rule.Delete},
				})
			}
		}
	}

	return result
}
