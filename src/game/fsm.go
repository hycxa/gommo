
package game

import "proto"

type PacketID proto.PacketID

const onEnter, onLeave = 1, 2

type Handler func(args interface{}) string

type Trans struct {
	Curstate  *State
	NextState *State
}

type State struct {
	Name         string
	Initstate    *State
	CurrentChild *State
	Parent      *State
	Transitions map[PacketID]Trans
	Behaviors   *Behavior
}



func NewState(name string, initstate *State) (state *State) {
	if name == "" {
		return nil
	}

	state_obj := new(State)
	state_obj.Name = name
	state_obj.CurrentChild = initstate
	state_obj.Initstate = initstate
	state_obj.Initialize()
	return state_obj
}
func (state *State) Initialize() {
	state.CurrentChild = state.Initstate
	state.Transitions = make(map[PacketID]Trans)
	if state.CurrentChild != nil {
		state.CurrentChild.SetParent(state)
		state.CurrentChild.Initialize()
	}

}
func (state *State) ProcessBehaviors(result *[]string, Id PacketID, args interface{}) {
	if state.Behaviors != nil {
		state.Behaviors.Process(result, Id, args)
	}
}
func (state *State) ProcessChildState(result *[]string, Id PacketID, args interface{}) {
	if state.CurrentChild != nil {
		state.CurrentChild.Process(result, Id, args)
	}
}
func (state *State) Process(result *[]string, Id PacketID, args interface{}) {

	if Id == onLeave {
		state.ProcessChildState(result, Id, args) 
		state.ProcessBehaviors(result, Id, nil)
	} else if Id == onEnter {
		state.ProcessBehaviors(result, Id, nil)
		state.ProcessChildState(result, Id, args)
	} else {
		state.ProcessBehaviors(result, Id, args)
		state.ProcessChildState(result, Id, args)
	}
	if state.CurrentChild != nil {
		nextState := state.FindTransition(Id)
		if nextState != nil && state.CurrentChild != nextState {
			state.SwitchState(result, nextState, Id, args)
			state.ProcessChildState(result, Id, args)
		}

	}
	return
}

func (state *State) SwitchState(result *[]string, nextState *State, Id PacketID, args interface{}) {
	if nextState == nil {
		return
	}
	if state.CurrentChild != nil {
		state.ProcessChildState(result, onLeave, args)
		state.CurrentChild.Parent = nil

	}
	state.CurrentChild = nextState
	if state.CurrentChild != nil {
		state.CurrentChild.SetParent(state)
		state.CurrentChild.Initialize()
		state.ProcessChildState(result, onEnter, args)

	}
}
func (state *State) ReTrans(currentState *State, nextState *State, Id PacketID) {

	if currentState != nil {
		tr := Trans{currentState, nextState}
		state.Transitions[Id] = tr
	}

}
func (state *State) FindTransition(Id PacketID) (nextState *State) {

	tr := state.Transitions[Id]
	if state.CurrentChild == tr.Curstate {
		return tr.NextState
	}
	return nil

}

func (state *State) GetBehavior() *Behavior {
	if state.Behaviors == nil {
		state.Behaviors = NewBehavior(state.Name)
	}
	return state.Behaviors

}

func (state *State) SetParent(Parent *State) {
	if state.Parent == nil {
		state.Parent = Parent
	}

}


func (state *State)IsDescendant(cstate *State )bool{
	if cstate !=nil && cstate == state.CurrentChild{
		return true
	}
	if state.CurrentChild !=nil {
		return state.CurrentChild.IsDescendant(cstate)
	}
	return false
}
