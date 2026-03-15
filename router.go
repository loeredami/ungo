package ungo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type PathStats struct {
	latency int64
	hits    int64
}

type Router[In, Out any] struct {
	mu    sync.RWMutex
	paths []func(In) Out
	stats []PathStats
}

func NewRouter[In, Out any](candidates ...func(In) Out) *Router[In, Out] {
	return &Router[In, Out]{
		paths: candidates,
		stats: make([]PathStats, len(candidates)),
	}
}

func (r *Router[In, Out]) Transform(input In) Out {
	bestIdx := r.selectBestPath()

	start := time.Now()
	output := r.paths[bestIdx](input)
	duration := time.Since(start).Nanoseconds()

	r.updateStats(bestIdx, duration)
	return output
}

func (r *Router[In, Out]) selectBestPath() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	bestIdx := 0
	minLatency := atomic.LoadInt64(&r.stats[0].latency)

	for i := 1; i < len(r.stats); i++ {
		lat := atomic.LoadInt64(&r.stats[i].latency)
		if lat == 0 {
			return i
		}
		if lat < minLatency {
			minLatency = lat
			bestIdx = i
		}
	}
	return bestIdx
}

func (r *Router[In, Out]) updateStats(idx int, duration int64) {
	oldLat := atomic.LoadInt64(&r.stats[idx].latency)
	var newLat int64
	if oldLat == 0 {
		newLat = duration
	} else {
		newLat = (oldLat*9 + duration) / 10
	}

	atomic.StoreInt64(&r.stats[idx].latency, newLat)
	atomic.AddInt64(&r.stats[idx].hits, 1)
}

func (r *Router[In, Out]) GetDiagnostics() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reports := make([]string, len(r.stats))
	for i, s := range r.stats {
		reports[i] = fmt.Sprintf("Path %d: %dns avg (%d hits)",
			i, atomic.LoadInt64(&s.latency), atomic.LoadInt64(&s.hits))
	}
	return reports
}
