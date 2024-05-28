package services

import (
	"testing"

	mogs "github.com/finkabaj/hyde-bot/internals/backend/mocks"
	"github.com/finkabaj/hyde-bot/internals/utils/common"
	"github.com/finkabaj/hyde-bot/internals/utils/guild"
	"github.com/finkabaj/hyde-bot/internals/utils/rule"
	"github.com/stretchr/testify/assert"
)

var mockDb = mogs.NewDbMock()
var mockGuildService = mogs.NewMockGuildService()
var mockReactionService = NewReactionService(mogs.NewMockLogger(), mockDb, mockGuildService)

const rac = rule.ReactActionCount

func TestGetReactionRules(t *testing.T) {
	t.Run("Positive", testGetReactionRulesPositive)
	t.Run("NotFound", testGetReactionRulesNotFound)
	t.Run("Internal", testGetReactionRulesInternal)
}

func TestCreateReactionRules(t *testing.T) {
	t.Run("Positive", testCreateReactionRulesPositive)
	t.Run("MinLen", testCreateReactionRulesMinLen)
	t.Run("GuildNotFound", testCreateReactionRulesGuildNotFound)
	t.Run("GuildNotEqual", testCreateReactionRulesNotEqualGID)
	t.Run("EmojiIdConflict", testCreateReactionRulesEmojiIdConflict)
	t.Run("EmojiNameConflict", testCreateReactionRulesEmojiNameConflict)
	t.Run("EmptyEmojiIdAndName", testCreateReactionRulesEmptyEmojiIdAndName)
	t.Run("DuplicateActions", testCreateReactionRulesDuplicateActions)
	t.Run("DbReturnError", testCreateReactionRulesDbReturnError)
}

func TestDeleteReactionRules(t *testing.T) {
	t.Run("Positive", testDeleteReactionRulesPositive)
	t.Run("MinLen", testDeleteReactionRulesMinLen)
	t.Run("NotFound", testDeleteReactionRulesNotFound)
	t.Run("HaveNoEmojiIdAndName", testDeleteReactionRulesHaveNoEmojiIdAndName)
	t.Run("DbReturnError", testDeleteReactionRulesDbReturnError)
}

func testGetReactionRulesPositive(t *testing.T) {
	gId := "123131421"
	guild := guild.Guild{
		GuildId: gId,
		OwnerId: "12312",
	}
	expectedResult := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "asdsa",
			EmojiId:    "1231",
			Actions:    [rac]rule.ReactAction{1, 2},
		},
		{
			GuildId:    gId,
			RuleAuthor: "fsd",
			EmojiName:  "das",
			Actions:    [rac]rule.ReactAction{1},
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild, nil)
	mockDb.On("ReadReactionRules", gId).Return(expectedResult, nil)

	actualResult, err := mockReactionService.GetReactionRules(gId)

	assert.Equal(t, expectedResult, actualResult)
	assert.Nil(t, err)

	mockDb.AssertExpectations(t)
	mockGuildService.AssertExpectations(t)
}

func testGetReactionRulesNotFound(t *testing.T) {
	gId := "123132"

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, common.ErrNotFound)

	actualResult, err := mockReactionService.GetReactionRules(gId)

	assert.Equal(t, []rule.ReactionRule{}, actualResult)
	assert.Equal(t, common.ErrNotFound, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "ReadReactionRules")
}

func testGetReactionRulesInternal(t *testing.T) {
	gId := "123"
	guild := guild.Guild{
		GuildId: gId,
		OwnerId: "sda",
	}

	mockGuildService.On("GetGuild", gId).Return(guild, nil)
	mockDb.On("ReadReactionRules", gId).Return([]rule.ReactionRule{}, common.ErrInternal)

	actualResult, err := mockReactionService.GetReactionRules(gId)

	assert.Equal(t, []rule.ReactionRule{}, actualResult)
	assert.Equal(t, common.ErrInternal, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

func testCreateReactionRulesPositive(t *testing.T) {
	gId := "1ass"
	expectedResult := []rule.ReactionRule{
		{
			EmojiName:  "🚌",
			RuleAuthor: "me",
			GuildId:    gId,
			Actions:    [rac]rule.ReactAction{0, 1},
		},
		{
			EmojiId:    "1337",
			RuleAuthor: "not me",
			GuildId:    gId,
			Actions:    [rac]rule.ReactAction{2},
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("ReadReactionRules", gId).Return([]rule.ReactionRule{}, nil)
	mockDb.On("CreateReactionRules", expectedResult).Return(expectedResult, nil)

	actualResult, err := mockReactionService.CreateReactionRules(expectedResult)

	assert.Equal(t, expectedResult, actualResult)
	assert.Nil(t, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

func testCreateReactionRulesMinLen(t *testing.T) {
	actualResponse, err := mockReactionService.CreateReactionRules([]rule.ReactionRule{})

	assert.Equal(t, []rule.ReactionRule{}, actualResponse)
	assert.Equal(t, common.ErrBadRequest, err)

	mockGuildService.AssertNotCalled(t, "GetGuild")
	mockDb.AssertNotCalled(t, "ReadReactionRules")
	mockDb.AssertNotCalled(t, "CreateReactionRules")
}

func testCreateReactionRulesGuildNotFound(t *testing.T) {
	gId := "bust dat nut"

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, common.ErrNotFound)

	actualResponse, err := mockReactionService.CreateReactionRules([]rule.ReactionRule{{
		GuildId:    gId,
		RuleAuthor: "me)",
		EmojiId:    "131",
		Actions:    [rac]rule.ReactAction{0},
	}})

	assert.Equal(t, []rule.ReactionRule{}, actualResponse)
	assert.Equal(t, common.ErrNotFound, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "ReadReactionRules")
	mockDb.AssertNotCalled(t, "CreateReactionRules")
}

func testCreateReactionRulesNotEqualGID(t *testing.T) {
	gId := "uk"
	rules := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "fsdf",
			EmojiId:    "fsfsgf",
			Actions:    [rac]rule.ReactAction{0},
		},
		{
			GuildId:    gId,
			RuleAuthor: "sdfds",
			EmojiName:  "1",
			Actions:    [rac]rule.ReactAction{0, 1},
		},
		{
			GuildId:    "fdsf",
			RuleAuthor: "fs",
			EmojiId:    "ffsd",
			Actions:    [rac]rule.ReactAction{0, 2},
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("ReadReactionRules", gId).Return([]rule.ReactionRule{}, nil)

	actualResponse, err := mockReactionService.CreateReactionRules(rules)

	assert.Equal(t, []rule.ReactionRule{}, actualResponse)
	assert.Equal(t, common.ErrBadRequest, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "CreateReactionRules")
}

func testCreateReactionRulesEmojiIdConflict(t *testing.T) {
	gId := "ua"
	rules := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "fsdf",
			EmojiId:    "fsfsgf",
			Actions:    [rac]rule.ReactAction{0},
		},
		{
			GuildId:    gId,
			RuleAuthor: "sdfds",
			EmojiId:    "vsd",
			Actions:    [rac]rule.ReactAction{0, 1},
		},
	}
	foundRules := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "fsdf",
			EmojiId:    "sdvsvs",
			Actions:    [rac]rule.ReactAction{0},
		},
		{
			GuildId:    gId,
			RuleAuthor: "sdfds",
			EmojiId:    "vsd",
			Actions:    [rac]rule.ReactAction{0, 1},
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("ReadReactionRules", gId).Return(foundRules, nil)

	actualResponse, err := mockReactionService.CreateReactionRules(rules)

	assert.Equal(t, []rule.ReactionRule{}, actualResponse)
	assert.Equal(t, rule.ErrRuleReactionConflict, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "CreateReactionRules")
}

func testCreateReactionRulesEmojiNameConflict(t *testing.T) {
	gId := "jkk"
	rules := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "fsdf",
			EmojiName:  "🚌",
			Actions:    [rac]rule.ReactAction{0},
		},
		{
			GuildId:    gId,
			RuleAuthor: "sdfds",
			EmojiName:  "131",
			Actions:    [rac]rule.ReactAction{0, 1},
		},
	}
	foundRules := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "fsdf",
			EmojiName:  "🚌",
			Actions:    [rac]rule.ReactAction{0},
		},
		{
			GuildId:    gId,
			RuleAuthor: "sdfds",
			EmojiName:  "1",
			Actions:    [rac]rule.ReactAction{0, 1},
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("ReadReactionRules", gId).Return(foundRules, nil)

	actualResponse, err := mockReactionService.CreateReactionRules(rules)

	assert.Equal(t, []rule.ReactionRule{}, actualResponse)
	assert.Equal(t, rule.ErrRuleReactionConflict, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "CreateReactionRules")
}

func testCreateReactionRulesEmptyEmojiIdAndName(t *testing.T) {
	gId := "goo"
	rules := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "fsdf",
			Actions:    [rac]rule.ReactAction{0},
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("ReadReactionRules", gId).Return([]rule.ReactionRule{}, nil)

	actualResponse, err := mockReactionService.CreateReactionRules(rules)

	assert.Equal(t, []rule.ReactionRule{}, actualResponse)
	assert.Equal(t, common.ErrBadRequest, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "CreateReactionRules")
}

func testCreateReactionRulesDuplicateActions(t *testing.T) {
	gId := "beepboop"
	rules := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "fsdf",
			EmojiId:    "123",
			Actions:    [rac]rule.ReactAction{1, 1, 0, 0},
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("ReadReactionRules", gId).Return([]rule.ReactionRule{}, nil)

	actualResponse, err := mockReactionService.CreateReactionRules(rules)

	assert.Equal(t, []rule.ReactionRule{}, actualResponse)
	assert.Equal(t, common.ErrBadRequest, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "CreateReactionRules")
}

func testCreateReactionRulesDbReturnError(t *testing.T) {
	gId := "beepboop"
	rules := []rule.ReactionRule{
		{
			GuildId:    gId,
			RuleAuthor: "fsdf",
			EmojiId:    "123",
			Actions:    [rac]rule.ReactAction{0},
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("ReadReactionRules", gId).Return([]rule.ReactionRule{}, nil)
	mockDb.On("CreateReactionRules", rules).Return([]rule.ReactionRule{}, common.ErrInternal)

	actualResponse, err := mockReactionService.CreateReactionRules(rules)

	assert.Equal(t, []rule.ReactionRule{}, actualResponse)
	assert.Equal(t, common.ErrInternal, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

func testDeleteReactionRulesPositive(t *testing.T) {
	gId := "del"
	rules := []rule.DeleteReactionRuleQuery{
		{
			EmojiId: "1",
		},
		{
			EmojiName: "a",
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("DeleteReactionRules", rules, gId).Return(nil)

	err := mockReactionService.DeleteReactionRules(rules, gId)

	assert.Nil(t, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

func testDeleteReactionRulesMinLen(t *testing.T) {
	err := mockReactionService.DeleteReactionRules([]rule.DeleteReactionRuleQuery{}, "1")

	assert.Equal(t, common.ErrBadRequest, err)

	mockGuildService.AssertNotCalled(t, "GetGuild")
	mockDb.AssertNotCalled(t, "DeleteReactionRules")
}

func testDeleteReactionRulesNotFound(t *testing.T) {
	gId := "dllss"
	rules := []rule.DeleteReactionRuleQuery{
		{
			EmojiId: "1",
		},
		{
			EmojiName: "a",
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, common.ErrNotFound)

	err := mockReactionService.DeleteReactionRules(rules, gId)

	assert.Equal(t, common.ErrNotFound, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "DeleteReactionRules")
}

func testDeleteReactionRulesHaveNoEmojiIdAndName(t *testing.T) {
	gId := "del"
	rules := []rule.DeleteReactionRuleQuery{
		{
			EmojiId:   "",
			EmojiName: "",
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)

	err := mockReactionService.DeleteReactionRules(rules, gId)

	assert.Equal(t, common.ErrBadRequest, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertNotCalled(t, "DeleteReactionRules")
}

func testDeleteReactionRulesDbReturnError(t *testing.T) {
	gId := "dl"
	rules := []rule.DeleteReactionRuleQuery{
		{
			EmojiId: "1",
		},
		{
			EmojiName: "a",
		},
	}

	mockGuildService.On("GetGuild", gId).Return(guild.Guild{}, nil)
	mockDb.On("DeleteReactionRules", rules, gId).Return(common.ErrInternal)

	err := mockReactionService.DeleteReactionRules(rules, gId)

	assert.Equal(t, common.ErrInternal, err)

	mockGuildService.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}
