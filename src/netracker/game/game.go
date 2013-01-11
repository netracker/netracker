package game

import "netracker/player"

type Game struct {
	ActivePlayer   *player.Player
	InactivePlayer *player.Player
	CorpCredits    int
	RunnerCredits  int
	Clicks         int
}

func New() *Game {
	return &Game{
		ActivePlayer:   player.NewCorp(),
		InactivePlayer: player.NewRunner(),
		CorpCredits:    5,
		RunnerCredits:  5,
		Clicks:         0,
	}
}

func (game *Game) NextTurn() {
	game.ActivePlayer, game.InactivePlayer = game.InactivePlayer, game.ActivePlayer
	game.Clicks = 0
}

func (game *Game) UseClick() {
	game.Clicks += 1
}

func (game *Game) AddCorpCredit() {
	game.CorpCredits += 1
}

func (game *Game) AddRunnerCredit() {
	game.RunnerCredits += 1
}

func (game *Game) RemoveRunnerCredit() {
	if game.RunnerCredits != 0 {
		game.RunnerCredits -= 1
	}
}

func (game *Game) RemoveCorpCredit() {
	if game.CorpCredits != 0 {
		game.CorpCredits -= 1
	}
}
