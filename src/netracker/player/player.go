package player

const (
	RUNNER = "runner"
	CORP   = "corp"
)

type Player struct {
	Role string
}

func NewRunner() *Player {
	return &Player{Role: RUNNER}
}

func NewCorp() *Player {
	return &Player{Role: CORP}
}
