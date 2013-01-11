package player

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestNewRunnerBuildsPlayerWithCorrectRole(t *testing.T) {
	player := NewRunner()
	assert.Equal(t, player.Role, RUNNER)
}

func TestNewCorpBuildsPlayerWithCorrectRole(t *testing.T) {
	player := NewCorp()
	assert.Equal(t, player.Role, CORP)
}
