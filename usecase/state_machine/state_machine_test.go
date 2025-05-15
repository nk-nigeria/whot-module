package state_machine

import (
	"fmt"
	"testing"
)

func TestStateMachine(t *testing.T) {
	t.Logf("test game state machine")

	game := NewGameStateMachine()
	fmt.Printf("State is %v\n", game.MustState())

	fmt.Printf("State is %v\n", game.MustState())
}
