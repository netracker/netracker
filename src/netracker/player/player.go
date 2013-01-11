package player

const (
	RUNNER = "runner"
	CORP   = "corp"
)

type Player struct {
	Role string
}

func NewPlayer(role string) *Player {
	return &Player{Role: role}
}
