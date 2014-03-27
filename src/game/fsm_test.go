package game

import (
	"testing"
)



func assertequles(t *testing.T, s string, want string) {

	if s != want {
		t.Error(" Not want : ", s, want)
	}
}
func assert(t *testing.T, s bool) {
	if s == false {
		t.Error("Not True!")
	}
}
func TestRun(t *testing.T) {
	if testing.Short() {
		return
	}

	state1 := NewState("one", nil)

	statebehavior1 := state1.GetBehavior()
	statebehavior1.ReTigger(onEnter,
		func(args interface{}) string {
			return "b1.Begin"
		})

	statebehavior1.ReTigger(onLeave,
		func(args interface{}) string {
			return "b1.End"
		})

	statebehavior1.ReTigger(111,
		func(args interface{}) string {
			return "b1.1"
		})

	statebehavior1Child1 := NewBehavior("testbehavior1Child1")
	statebehavior1.AddChild(statebehavior1Child1)
	statebehavior1Child1.ReTigger(onEnter,
		func(args interface{}) string {
			return "b1.c1.Begin"
		})
	statebehavior1Child1.ReTigger(onLeave,
		func(args interface{}) string {
			return "b1.c1.End"
		})
	statebehavior1Child1.ReTigger(111,
		func(args interface{}) string {
			return "b1.c1.1"
		})
	state3 := NewState("three", nil)
	state2 := NewState("two", state3)

	statebehavior2 := state2.GetBehavior()
	statebehavior2.ReTigger(onEnter,
		func(args interface{}) string {
			return "b2.Begin"
		})

	statebehavior2.ReTigger(onLeave,
		func(args interface{}) string {
			return "b2.End"
		})

	statebehavior2.ReTigger(111,
		func(args interface{}) string {
			return "b2.1"
		})

	statebehavior2.ReTigger(25,
		func(args interface{}) string {
			return "b2.25"
		})

	statebehavior2.ReTigger(26,
		func(args interface{}) string {
			return "b2.26"
		})
	state4 := NewState("four", state1)

	result := make([]string, 1, 20)

	state4.Process(&result, onEnter, "helloworld")
	assertequles(t, result[1], "b1.Begin")
	assertequles(t, result[2], "b1.c1.Begin")

	state4.ReTrans(state1, state2, 1003)
	assert(t, !state4.IsDescendant(state4))
	assert(t, state4.IsDescendant(state1))

	result = make([]string, 1, 20)
	state4.Process(&result, 111, "")
	assertequles(t, result[1], "b1.1")
	assertequles(t, result[2], "b1.c1.1")

	result = make([]string, 1, 20)
	state4.Process(&result, 1003, "")
	assertequles(t, result[1], "b1.c1.End")
	assertequles(t, result[2], "b1.End")
	assertequles(t, result[3], "b2.Begin")
	assert(t, state4.IsDescendant(state2))

	result = make([]string, 1, 20)
	state4.Process(&result, 25, "")
	assertequles(t, result[1], "b2.25")

	result = make([]string, 1, 20)
	state4.Process(&result, 26, "")
	assertequles(t, result[1], "b2.26")

}
