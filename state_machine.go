package ungo

type StateMachine[S comparable] struct {
	current     S
	transitions map[S]map[string]S
}

func (f *StateMachine[S]) Trigger(event string) {
	cur, ok := f.transitions[f.current]
	if !ok {
		return
	}
	next, ok := cur[event]
	if !ok {
		return
	}
	f.current = next
}

func (f *StateMachine[S]) Current() S {
	return f.current
}

func NewStateMachine[S comparable](initial S, transitions map[S]map[string]S) *StateMachine[S] {
	return &StateMachine[S]{
		current:     initial,
		transitions: transitions,
	}
}
