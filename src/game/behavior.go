package game

//import "fmt"

type Behavior struct {
	Name      string
	Handlers  map[PacketID]Handler
	Chirldren []*Behavior
	Parent    *Behavior
}

func NewBehavior(name string) *Behavior {
	if name == "" {
		return nil
	}
	behavior_obj := new(Behavior)
	behavior_obj.Chirldren = make([]*Behavior, 1)
	behavior_obj.Name = name
	behavior_obj.Handlers = make(map[PacketID]Handler)
	return behavior_obj
}

func (behavior *Behavior) ReTigger(Id PacketID, handler Handler) {

	behavior.Handlers[Id] = handler
}

func (behavior *Behavior) AddChild(chirldbehavior *Behavior) {

	for _, v := range behavior.Chirldren {
		if v == chirldbehavior {
			return
		}
	}
	behavior.Chirldren = append(behavior.Chirldren, chirldbehavior)


	chirldbehavior.Parent = behavior

}
func (behavior *Behavior) Process(result *[]string, Id PacketID, args interface{}) {
	if Id == onLeave {
		for i, v := range behavior.Chirldren {
			if i == 0 {
				continue
			}
			v.Process(result, Id, args)
			v.Parent = nil
			v.ClearChildren()
		}

	}
	handler := behavior.FindHandler(Id)
	if handler != nil {

		*result = append(*result, handler(args))
	}
	if Id != onLeave {
		for i, v := range behavior.Chirldren {
			if i != 0 {

				v.Process(result, Id, args)
			}
		}
	}

}
func (behavior *Behavior) ClearChildren() {
	return
	for i, v := range behavior.Chirldren {
		if i == 0 {
			continue
		}
		v.Parent = nil
		v.ClearChildren()
	}
	behavior.Chirldren = make([]*Behavior, 1)
}

func (behavior *Behavior) FindHandler(Id PacketID) Handler {
	return behavior.Handlers[Id]

}
