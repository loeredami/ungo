package ungo

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ExitUnknownError       = -1
	ExitSuccess            = 0
	ExitGenericFailure     = 1
	ExitStackOverflow      = 2
	ExitOutOfMemory        = 3
	ExitInvalidInput       = 4
	ExitConfigurationError = 5
)

var (
	signals = NewLazy(func() chan os.Signal {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		return sig
	})
)

func cancelWorkers() {
	if !workers.IsInitialized() {
		return
	}
	workers.Value().ForEach(func(key int, worker *Worker) {
		worker.Cancel()
	})
}

func waitForWorkers() {
	if !workers.IsInitialized() {
		return
	}
	for workers.Value().Size() > 0 {
		time.Sleep(time.Millisecond * 10)
	}
}

func Main(m func(argc int, argv []string) Optional[int], shouldWaitForWorkers Optional[bool], subProcesses Optional[[]SubProcess]) {
	go func() {
		if subProcesses.HasValue() {
			go func() {
				waiters := make([]SubProcessLink, 0, len(subProcesses.Value()))
				for _, sp := range subProcesses.Value() {
					waiters = append(waiters, StartSubProcess(sp.fn, sp.on_finish))
				}
				for _, waiter := range waiters {
					WaitForSubProcess(waiter)
				}
			}()
		}

		exitCode := m(len(os.Args), os.Args)
		if exitCode.HasValue() {
			cancelWorkers()
			os.Exit(exitCode.Value())
		}

		if shouldWaitForWorkers.HasValue() && shouldWaitForWorkers.Value() {
			waitForWorkers()
		}
		os.Exit(ExitSuccess)
	}()

	for {
		if workers.IsInitialized() {
			workers.Value().ForEach(func(key int, worker *Worker) {
				if worker.isCancelled {
					workers.Value().Delete(key)
				}
			})
		}

		if signals.IsInitialized() {
			select {
			case <-signals.Value():
				cancelWorkers()
				os.Exit(ExitSuccess)
			default:
			}
		}

		time.Sleep(time.Millisecond * 100)
	}
}

type SubProcessLink chan struct{}
type SubProcess struct {
	fn        func() Optional[int]
	on_finish Optional[func(Optional[int])]
}

func StartSubProcess(m func() Optional[int], on_finish Optional[func(Optional[int])]) SubProcessLink {
	ch := make(SubProcessLink)
	go func() {
		exitCode := m()
		if on_finish.HasValue() {
			on_finish.Value()(exitCode)
		}
		close(ch)
	}()
	return ch
}

func WaitForSubProcess(ch SubProcessLink) {
	<-ch
}
