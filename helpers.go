package ungo

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var NULL_ANY = EmptyOptional[any]()

func If[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func ShortPipe[T any](value T, fs []func(T) T) T {
	for _, f := range fs {
		value = f(value)
	}
	return value
}

func Memoize[K comparable, V any](fn func(K) V) func(K) V {
	cache := make(map[K]V)
	return func(k K) V {
		if v, ok := cache[k]; ok {
			return v
		}
		v := fn(k)
		cache[k] = v
		return v
	}
}

func SizedMemoize[K comparable, V any](maxSize int, fn func(K) V) func(K) V {
	cache := make(map[K]V)
	return func(k K) V {
		if v, ok := cache[k]; ok {
			return v
		}
		v := fn(k)
		cache[k] = v
		if len(cache) > maxSize {
			for k := range cache {
				delete(cache, k)
				break
			}
		}
		return v
	}
}

type Variable[T any] struct {
	val  T
	done chan struct{}
}

func NewVariable[T any]() *Variable[T] {
	return &Variable[T]{done: make(chan struct{})}
}

func (v *Variable[T]) Set(val T) {
	v.val = val
	close(v.done)
}

func (v *Variable[T]) Get() T {
	<-v.done
	return v.val
}

type Monoid[T any] struct {
	Identity T
	Combine  func(T, T) T
}

func Fold[T any](items []T, m Monoid[T]) T {
	result := m.Identity
	for _, item := range items {
		result = m.Combine(result, item)
	}
	return result
}

type FunReader[Env, Res any] func(Env) Res

func FunMap[Env, A, B any](r FunReader[Env, A], f func(A) B) FunReader[Env, B] {
	return func(e Env) B {
		return f(r(e))
	}
}

func ApplyAll[T any](slice []T, action func(T) T) []T {
	result := make([]T, len(slice))
	for i, v := range slice {
		result[i] = action(v)
	}
	return result
}

func Contract[In, Out any](
	fn func(In) Out,
	require func(In) bool,
	ensure func(Out) bool,
) func(In) Result[Out] {
	return func(in In) Result[Out] {
		if !require(in) {
			return Result[Out]{err: fmt.Errorf("requirement failed")}
		}
		out := fn(in)
		if !ensure(out) {
			return Result[Out]{err: fmt.Errorf("ensure failed")}
		}
		return Result[Out]{value: out}
	}
}

func ReinterpretCast[Src, Dest any](src Src) Dest {
	var dest Dest
	*(*unsafe.Pointer)(unsafe.Pointer(&dest)) = unsafe.Pointer(&src)
	return dest
}

type ContextVar[Env any] struct {
	val       Env
	terminate func(Env) Optional[error]
}

func NewContextVar[Env any]() *ContextVar[Env] {
	return &ContextVar[Env]{}
}

func (v *ContextVar[Env]) SetTerminator(terminate func(Env) Optional[error]) {
	v.terminate = terminate
}

func (v *ContextVar[Env]) Set(val Env) {
	v.val = val
}

func (v *ContextVar[Env]) Get() Env {
	return v.val
}

func With[Env any](open func() ContextVar[Env], action func(Env)) Optional[error] {
	ctx := open()
	action(ctx.Get())
	if ctx.terminate != nil {
		return ctx.terminate(ctx.Get())
	}
	return Optional[error]{valid: false}
}

func Retry[Res any](action func() (Res, error), retries int, delay time.Duration) Result[Res] {
	var res Res
	var err error
	for range retries {
		res, err = action()
		if err == nil {
			return Result[Res]{value: res}
		}
		if delay != 0 {
			time.Sleep(delay)
		}
	}
	return Result[Res]{err: err}
}

func Patch[T any](base *T, patch T) {
	b := reflect.ValueOf(base).Elem()
	p := reflect.ValueOf(patch)

	for i := 0; i < b.NumField(); i++ {
		pField := p.Field(i)
		if !pField.IsZero() {
			b.Field(i).Set(pField)
		}
	}
}

func Debounce[T any](interval time.Duration, fn func(T)) func(T) {
	var mu sync.Mutex
	var timer *time.Timer

	return func(arg T) {
		mu.Lock()
		defer mu.Unlock()

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(interval, func() {
			fn(arg)
		})
	}
}

func LockFun[T any](fn func(T)) func(T) {
	var mu sync.Mutex

	return func(arg T) {
		mu.Lock()
		defer mu.Unlock()

		fn(arg)
	}
}

func RlockFun[T any](fn func(T)) func(T) {
	var mu sync.RWMutex

	return func(arg T) {
		mu.Lock()
		defer mu.Unlock()

		fn(arg)
	}
}

func WithBackpressure[T any](ch chan<- T, onFull func(T)) func(T) bool {
	return func(val T) bool {
		select {
		case ch <- val:
			return true
		default:
			if onFull != nil {
				onFull(val)
			}
			return false
		}
	}
}

func Coalesce[K comparable, V any](fn func(K) (V, error)) func(K) (V, error) {
	var mu sync.Mutex
	inflight := make(map[K]*sync.WaitGroup)
	cache := make(map[K]*Result[V])

	return func(key K) (V, error) {
		mu.Lock()
		if wg, ok := inflight[key]; ok {
			mu.Unlock()
			wg.Wait() // Wait for the first caller
			return cache[key].value, cache[key].err
		}

		wg := &sync.WaitGroup{}
		wg.Add(1)
		inflight[key] = wg
		mu.Unlock()

		v, e := fn(key)

		mu.Lock()
		cache[key] = &Result[V]{v, e}
		delete(inflight, key)
		wg.Done() // Release all waiting callers
		mu.Unlock()

		return v, e
	}
}

func Pooled[T any](constructor func() T) func(func(T)) {
	pool := &sync.Pool{New: func() any { return constructor() }}

	return func(action func(T)) {
		item := pool.Get().(T)
		action(item)
		pool.Put(item)
	}
}

func Adaptive[T, R any](fn func(T) R) func(T) Result[R] {
	var avgLatency int64
	threshold := int64(time.Millisecond * 500) // Shed load if > 500ms

	return func(arg T) Result[R] {
		if atomic.LoadInt64(&avgLatency) > threshold {
			return Result[R]{err: errors.New("load shedding: latency too high")}
		}

		start := time.Now()
		res := fn(arg)
		duration := time.Since(start).Nanoseconds()

		newAvg := (atomic.LoadInt64(&avgLatency)*9 + duration) / 10
		atomic.StoreInt64(&avgLatency, newAvg)

		return Result[R]{value: res}
	}
}

func ForEach[T any](s []T, fn func(T)) {
	for _, v := range s {
		fn(v)
	}
}

func WithValve[T, R any](limit int, fn func(T) R) func(T) R {
	sem := make(chan struct{}, limit)
	return func(arg T) R {
		sem <- struct{}{}
		defer func() { <-sem }()
		return fn(arg)
	}
}

func RoundRobin[T any](input <-chan T, workers int) []chan T {
	outputs := make([]chan T, workers)
	for i := range outputs {
		outputs[i] = make(chan T)
	}

	go func() {
		defer func() {
			for _, ch := range outputs {
				close(ch)
			}
		}()
		i := 0
		for item := range input {
			outputs[i] <- item
			i = (i + 1) % workers
		}
	}()
	return outputs
}

func Sample[T any](input <-chan T, interval time.Duration) <-chan T {
	output := make(chan T)
	go func() {
		defer close(output)
		ticker := time.NewTicker(interval)
		for range ticker.C {
			item, ok := <-input
			if !ok {
				return
			}
			output <- item
		}
	}()
	return output
}

func Dedup[T comparable](input <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		var last T
		var first = true
		for item := range input {
			if first || item != last {
				out <- item
				last = item
				first = false
			}
		}
	}()
	return out
}

type Tracked[T any] struct {
	ID        UUID
	CreatedAt time.Time
	Data      T
}

func WrapTracked[T any](data T) Tracked[T] {
	return Tracked[T]{
		ID:        NewUUID(),
		CreatedAt: time.Now(),
		Data:      data,
	}
}

func Tee[T any](input <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)
	go func() {
		defer close(out1)
		defer close(out2)
		for item := range input {
			// We must send to both, potentially in parallel
			// to prevent one slow reader from blocking the other
			var wg sync.WaitGroup
			wg.Add(2)
			go func(v T) { out1 <- v; wg.Done() }(item)
			go func(v T) { out2 <- v; wg.Done() }(item)
			wg.Wait()
		}
	}()
	return out1, out2
}

func WithTempFile(fn func(*os.File)) func() {
	return func() {
		f, _ := os.CreateTemp("", "worker-*")
		defer os.Remove(f.Name())
		defer f.Close()
		fn(f)
	}
}

func Prioritize[T any](high, low <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			select {
			case item, ok := <-high:
				if !ok {
					return
				}
				out <- item
			default:
				select {
				case item, ok := <-high:
					if !ok {
						return
					}
					out <- item
				case item, ok := <-low:
					if !ok {
						return
					}
					out <- item
				}
			}
		}
	}()
	return out
}

func Hedged[T any](delay time.Duration, task func() T) func() T {
	return func() T {
		res := make(chan T, 2)
		go func() { res <- task() }()

		select {
		case val := <-res:
			return val
		case <-time.After(delay):
			go func() { res <- task() }()
			return <-res // Returns the first one to finish
		}
	}
}

func WithDeadLetter[T any](primary, dlq chan<- T) func(T) {
	return func(item T) {
		select {
		case primary <- item:
		default:
			select {
			case dlq <- item:
			default:
			}
		}
	}
}

func Chunked[T any](size int, process func([]T)) func([]T) {
	return func(data []T) {
		for i := 0; i < len(data); i += size {
			end := i + size
			if end > len(data) {
				end = len(data)
			}
			process(data[i:end])
		}
	}
}

func WithTTL[T any](ttl time.Duration, fn func() T) func() T {
	var (
		mu     sync.Mutex
		expiry time.Time
		cached T
	)
	return func() T {
		mu.Lock()
		defer mu.Unlock()
		if time.Now().Before(expiry) {
			return cached
		}
		cached = fn()
		expiry = time.Now().Add(ttl)
		return cached
	}
}

func WithSignificance[T any](epsilon func(old, new T) bool, fn func(T)) func(T) {
	var last T
	var first = true
	var mu sync.Mutex

	return func(current T) {
		mu.Lock()
		defer mu.Unlock()

		if first || epsilon(last, current) {
			fn(current)
			last = current
			first = false
		}
	}
}

func GatedChannel[T any](input <-chan T, isHealthy func() bool) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for item := range input {
			// Spin-wait or block until healthy
			for !isHealthy() {
				time.Sleep(100 * time.Millisecond)
			}
			out <- item
		}
	}()
	return out
}

func WithQuorum[T comparable](tasks ...func() T) func() (T, error) {
	return func() (T, error) {
		counts := make(map[T]int)
		for _, t := range tasks {
			res := t()
			counts[res]++
			if counts[res] > len(tasks)/2 {
				return res, nil
			}
		}
		return *new(T), errors.New("no quorum reached")
	}
}

func Interlock(f1 func(), f2 func()) (func(), func()) {
	var mu sync.Mutex
	wrap1 := func() {
		mu.Lock()
		defer mu.Unlock()
		f1()
	}
	wrap2 := func() {
		mu.Lock()
		defer mu.Unlock()
		f2()
	}
	return wrap1, wrap2
}

func TrioMorph[A, B, C any](fn func(A) B, transform func(C) A) func(C) B {
	return func(input C) B {
		return fn(transform(input))
	}
}

func WithFuse[T any, R comparable](isSane func(R) bool, fn func(T) R) func(T) Result[R] {
	blown := false
	return func(arg T) Result[R] {
		if blown {
			return Result[R]{err: errors.New("fuse blown: insane data detected")}
		}
		res := fn(arg)
		if !isSane(res) {
			blown = true
			return Result[R]{value: res, err: errors.New("insane data: blowing fuse")}
		}
		return Result[R]{value: res}
	}
}
