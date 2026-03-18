package ungo

import "slices"

var workers Lazy[*FastMap[int, *Worker]] = Lazy[*FastMap[int, *Worker]]{
	initializer: func() *FastMap[int, *Worker] {
		return &FastMap[int, *Worker]{}
	},
}

type Worker struct {
	id          int
	result      Result[any]
	isRunning   bool
	isCancelled bool
	fn          func() Result[any]
}

func (w *Worker) GetResult() Result[any] {
	return w.result
}

func (w *Worker) SetResult(result Result[any]) {
	w.result = result
}

func GetWorker(id int) Optional[*Worker] {
	worker, ok := workers.Value().Get(id)
	if !ok {
		return EmptyOptional[*Worker]()
	}
	return MakeOptional(worker)
}

func MakeWorker(fn func() Result[any]) *Worker {
	var id int
	IDs := workers.Value().Keys()
	for slices.Contains(IDs, id) {
		id++
	}
	worker := &Worker{
		id: id,
		fn: fn,
	}
	workers.Value().Set(worker.id, worker)
	return worker
}

func (w *Worker) Run() {
	if w.isRunning {
		return
	}
	w.isRunning = true
	go func() {
		w.result = w.fn()
		workers.Value().Delete(w.id)
		w.isRunning = false
	}()
}

func (w *Worker) Cancel() {
	w.isCancelled = true
}
