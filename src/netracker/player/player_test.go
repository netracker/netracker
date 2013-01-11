package player

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewPlayerBuildsAPlayerWithARole(t *testing.T) {
	player := NewPlayer(RUNNER)
	assert.Equal(t, player.Role, RUNNER)
}
