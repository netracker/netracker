package parser

import (
	"github.com/bmizerany/assert"
	"netracker/game"
	"netracker/player"
	"testing"
)

func TestParseNextTurn(t *testing.T) {
	game := game.New()
	parser := MessageParser{game: game}
	parser.parse("nextturn")
	assert.Equal(t, game.ActivePlayer.Role, player.RUNNER)
}

func TestParseClick(t *testing.T) {
	game := game.New()
	parser := MessageParser{game: game}
	parser.parse("click")
	assert.Equal(t, game.Clicks, 1)
}

func TestParseAddRunnerCredit(t *testing.T) {
	game := game.New()
	parser := MessageParser{game: game}
	parser.parse("addrunnercredit")
	assert.Equal(t, game.RunnerCredits, 6)
}

func TestParseAddCorpCredit(t *testing.T) {
	game := game.New()
	parser := MessageParser{game: game}
	parser.parse("addcorpcredit")
	assert.Equal(t, game.CorpCredits, 6)
}

func TestParseRemoveCorpCredit(t *testing.T) {
	game := game.New()
	parser := MessageParser{game: game}
	parser.parse("removecorpcredit")
	assert.Equal(t, game.CorpCredits, 4)
}
