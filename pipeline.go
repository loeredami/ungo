package ungo

import "fmt"

type Pipeline[T any] struct {
	funcs []func(T) T
}

type Pipe[T any] interface {
	Run(T) T
}

type PipeFunc[T any] func(T) T

func (p PipeFunc[T]) Run(t T) T {
	return p(t)
}

type CondPipe[T any] struct {
	cond func(T) bool
	p    Pipe[T]
}

func (cp CondPipe[T]) Run(t T) T {
	if cp.cond(t) {
		return cp.p.Run(t)
	}
	return t
}

type BranchPipe[T any] struct {
	cond func(T) bool
	p1   Pipe[T]
	p2   Pipe[T]
}

func (bp BranchPipe[T]) Run(t T) T {
	if bp.cond(t) {
		return bp.p1.Run(t)
	}
	return bp.p2.Run(t)
}

type PipeSequence[T any] struct {
	pipes []Pipe[T]
}

func (ps PipeSequence[T]) Run(t T) T {
	for _, p := range ps.pipes {
		t = p.Run(t)
	}
	return t
}

func NewPipeSequence[T any](pipes ...Pipe[T]) PipeSequence[T] {
	return PipeSequence[T]{pipes: pipes}
}

type SwitchPipe[T any] struct {
	m map[*T]Pipe[T]
}

func NewSwitchPipe[T any](m map[*T]Pipe[T]) SwitchPipe[T] {
	return SwitchPipe[T]{m: m}
}

func (sp SwitchPipe[T]) Run(t T) T {
	if p, ok := sp.m[&t]; ok {
		return p.Run(t)
	}
	return t
}

type LoopPipe[T any] struct {
	cond func(T) bool
	p    Pipe[T]
}

func NewLoopPipe[T any](cond func(T) bool, p Pipe[T]) LoopPipe[T] {
	return LoopPipe[T]{cond: cond, p: p}
}

func (lp LoopPipe[T]) Run(t T) T {
	for lp.cond(t) {
		t = lp.p.Run(t)
	}
	return t
}

type BurnPipe[T any] struct {
	f        func(T) Result[T]
	failsafe func(T) T
	silent   bool
}

func NewBurnPipe[T any](f func(T) Result[T], failsafe func(T) T, silent bool) BurnPipe[T] {
	return BurnPipe[T]{f: f, failsafe: failsafe, silent: silent}
}

func (bp BurnPipe[T]) Run(t T) T {
	r := bp.f(t)
	if r.err != nil {
		if !bp.silent {
			fmt.Println("pipe burned: ", r.err)
		}
		return bp.failsafe(t)
	}
	return r.value
}
