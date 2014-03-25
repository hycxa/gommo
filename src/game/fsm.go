// state.go
package gogame

import ()

type State struct {
	Current *State
}

func NewState() *State {
	return &State{}
}
