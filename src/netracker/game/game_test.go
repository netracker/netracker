package game

import (
	"github.com/bmizerany/assert"
	"netracker/player"
	"testing"
)

func TestNewGame(t *testing.T) {
	game := New()
	assert.Equal(t, game.ActivePlayer.Role, player.CORP)
	assert.Equal(t, game.InactivePlayer.Role, player.RUNNER)
	assert.Equal(t, game.CorpCredits, 5)
	assert.Equal(t, game.RunnerCredits, 5)
}

func TestNextTurn(t *testing.T) {
	game := New()
	game.NextTurn()
	assert.Equal(t, game.ActivePlayer.Role, player.RUNNER)
	assert.Equal(t, game.InactivePlayer.Role, player.CORP)
}

func TestNextTurnResetsClicks(t *testing.T) {
	game := New()
	game.Clicks = 4
	game.NextTurn()
	assert.Equal(t, game.Clicks, 0)
}

func TestUseClickIncreasesClicksForActivePlayer(t *testing.T) {
	game := New()
	assert.Equal(t, game.Clicks, 0)
	game.UseClick()
	assert.Equal(t, game.Clicks, 1)
}

func TestAddCorpCredit(t *testing.T) {
	game := New()

	assert.Equal(t, game.CorpCredits, 5)
	game.AddCorpCredit()
	assert.Equal(t, game.CorpCredits, 6)
}

func TestAddRunnerCredit(t *testing.T) {
	game := New()

	assert.Equal(t, game.RunnerCredits, 5)
	game.AddRunnerCredit()
	assert.Equal(t, game.RunnerCredits, 6)
}

func TestRemoveRunnerCredit(t *testing.T) {
	game := New()

	assert.Equal(t, game.RunnerCredits, 5)
	game.RemoveRunnerCredit()
	assert.Equal(t, game.RunnerCredits, 4)
}

func TestRemoveRunnerCreditFloorsAtZero(t *testing.T) {
	game := New()
	game.RunnerCredits = 0

	game.RemoveRunnerCredit()
	assert.Equal(t, game.RunnerCredits, 0)
}

func TestRemoveCorpCredit(t *testing.T) {
	game := New()

	assert.Equal(t, game.CorpCredits, 5)
	game.RemoveCorpCredit()
	assert.Equal(t, game.CorpCredits, 4)
}

func TestRemoveCorpCreditFloorsAtZero(t *testing.T) {
	game := New()
	game.CorpCredits = 0

	game.RemoveCorpCredit()
	assert.Equal(t, game.CorpCredits, 0)
}
