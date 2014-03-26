// state.go
package game

import ()

type State struct {
	Current *State
}

func NewState() *State {
	return &State{}
}
