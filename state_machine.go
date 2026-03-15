package ungo

type StateMachine[S comparable, E any] struct {
	current     S
	transitions map[S]map[string]S
}

func (f *StateMachine[S, E]) Trigger(event string) {
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

func (f *StateMachine[S, E]) Current() S {
	return f.current
}

func NewStateMachine[S comparable, E any](initial S, transitions map[S]map[string]S) *StateMachine[S, E] {
	return &StateMachine[S, E]{
		current:     initial,
		transitions: transitions,
	}
}
