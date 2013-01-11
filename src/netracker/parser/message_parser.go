package parser

import (
	"netracker/game"
)

type MessageParser struct {
	game *game.Game
}

func (parser *MessageParser) parse(message string) {
	switch message {
	case "nextturn":
		parser.game.NextTurn()
	case "click":
		parser.game.UseClick()
	case "addrunnercredit":
		parser.game.AddRunnerCredit()
	case "addcorpcredit":
		parser.game.AddCorpCredit()
	case "removecorpcredit":
		parser.game.RemoveCorpCredit()
	case "removerunnercredit":
		parser.game.RemoveRunnerCredit()
	}
}
