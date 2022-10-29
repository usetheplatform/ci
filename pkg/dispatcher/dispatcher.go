package dispatcher

import (
	"fmt"
	"time"

	. "github.com/usetheplatform/ci-system/pkg/common"
	"github.com/usetheplatform/ci-system/pkg/connection"
)

type Commit = string
type RunnerID = string
type DispatchedCommits = map[Commit]RunnerID
type PendingCommits = Queue[Commit]

type RunnerStatus = int

const (
	Available RunnerStatus = iota
	Busy
	Dead
)

type Runner struct {
	id         RunnerID
	connection connection.Connection
	status     RunnerStatus
}

type DispatcherOptions struct {
	CheckRunnersTimeout time.Duration
	RedistributeTimeout time.Duration
}

// @TODO: Make storage in the file system to make sure data is not lost
type Dispatcher struct {
	config            *DispatcherOptions
	dispatchedCommits DispatchedCommits // how to make it thread safe?
	pendingCommits    PendingCommits    // https://stackoverflow.com/questions/2818852/is-there-a-queue-implementation
	registeredRunners []Runner          // how to make it thread safe?
}

func NewDispatcher(config *DispatcherOptions) Dispatcher {
	// @TODO: Infer sizes based on config and make registeredRunners dynamically
	return Dispatcher{
		config:            config,
		dispatchedCommits: make(DispatchedCommits),
		pendingCommits:    *NewQueue[Commit](10),
		registeredRunners: make([]Runner, 2),
	}
}

// TODO: Redo CheckIfError not to kill the process

func (d *Dispatcher) CheckRunners() {
	// ping each registered test runner to make sure they are still responsive

	poller := NewPoller(func() {
		Info("Checking runners...")

		for _, runner := range d.registeredRunners {
			err := runner.connection.Notify("status")

			if err != nil {
				Warning("Runner %s is not responding... Previous status: %v", runner.id, runner.status)
				runner.status = Dead
			}
			// check runner
		}
	}, d.config.CheckRunnersTimeout)
}

func (d *Dispatcher) dispatch(commit Commit) {
	allocated := false
	for _, runner := range d.registeredRunners {
		if runner.status == Available {
			// @TODO: Should notify the exact runner
			err := runner.connection.Notify(fmt.Sprintf("run:%s", commit))

			if err != nil {
				Warning("Failed to allocate commit %s to runner %s", commit, runner.id)
				d.pendingCommits.Append(commit)
				Warning("Commit %s has been returned to the queue", commit)
			} else {
				Info("%s has been dispatched to %s", commit, runner.id)
				allocated = true
				d.dispatchedCommits[commit] = runner.id
			}

			break
		}
	}

	if allocated == false {
		Warning("Failed to allocate commit %s since no runner were able to pick it up", commit)
		d.pendingCommits.Append(commit)
		Warning("Commit %s has been returned to the queue", commit)
	}
}

func (d *Dispatcher) AddRunner(runner Runner) {
	d.registeredRunners = append(d.registeredRunners, runner)
}

func (d *Dispatcher) RemoveRunner(id RunnerID) {
	// @TODO: Implement me
}

func (d *Dispatcher) UpdateRunnerStatus(id RunnerID, status RunnerStatus) {
	// @TODO: Implement me
}

func (d *Dispatcher) maintain() {
	// @TODO: Once a runner has complete the job with commit,
	// we should deallocate it from memory and save it to the fs somehow...
	// Also should be able to add runners / remove runners
}

func (d *Dispatcher) Enqueue(commit Commit) {
	err := d.pendingCommits.Append(commit)
	CheckIfError(err)
	Info("Received new commit %s that has been added to the queue", commit)
}

func (d *Dispatcher) Redistribute() {
	poller := NewPoller(func() {
		Info("Running redistribute...")

		for !d.pendingCommits.Empty() {
			// @TODO: Should pop only if runner is available and ready to pick it up
			commit, err := d.pendingCommits.Pop()
			CheckIfError(err)
			d.dispatch(*commit)
		}
	}, d.config.RedistributeTimeout)
}
