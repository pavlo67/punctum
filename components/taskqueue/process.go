package taskqueue

import (
	"time"

	"github.com/pavlo67/workshop/components/runner"

	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/joiner"
	"github.com/pavlo67/workshop/common/logger"
	"github.com/pavlo67/workshop/common/selectors"

	"github.com/pavlo67/workshop/components/tasks"
)

const timeToWait = time.Millisecond * 1000
const numToOmitSilently = 10

// TODO!!! use Process() in single-thread way only

func Process(tasksOp tasks.Operator, joinerOp joiner.Operator, l logger.Operator) {
	numOmitted := 0

	for {
		items, err := SelectTasksToProcess(tasksOp)
		if err != nil {
			l.Error("on SelectTasksToProcess(): ", err)
			time.Sleep(timeToWait)
			continue
		}
		if len(items) < 1 {
			numOmitted++
			if numOmitted >= numToOmitSilently {
				l.Infof("on SelectTasksToProcess(): no tasks to process, %d times", numOmitted)
				numOmitted = 0
			}

			time.Sleep(timeToWait)
			continue
		}

		numOmitted = 0

		for _, item := range items {
			workerOp, ok := joinerOp.Interface(item.ActorKey).(runner.Actor)
			if !ok {
				l.Errorf("no worker.Actor for task (%#v)", item)
				time.Sleep(timeToWait)
				continue
			}

			err = tasksOp.Start(item.ID, nil)
			if err != nil {
				l.Errorf("on tasksOp.Start(%s, nil): %s", item.ID, err)
				time.Sleep(timeToWait)
				continue
			}

			_, err := workerOp.Init(item.Task.Params)
			if err != nil {
				l.Errorf("on workerOp.Init(item.Task.Params) for item (%#v): %s", item, err)
			}

			// TODO!!! use goroutines
			posterior, err := workerOp.Run()
			if err != nil {
				l.Errorf("on workerOp.Run() for task (%#v): %s", item, err)
			}

			result := tasks.Result{
				// Timing: will be set automatically
				ErrStr:    err,
				Posterior: posterior,
			}
			err = tasksOp.Finish(item.ID, result, nil)
			if err != nil {
				l.Errorf("on tasksOp.Finish(%s, %#v, nil): %s", item.ID, result, err)
			}
		}
	}
}

func SelectTasksToProcess(tasksOp tasks.Operator) ([]tasks.Item, error) {
	return tasksOp.List(selectors.Binary(selectors.Eq, "status", selectors.Value{""}), &crud.GetOptions{Limit0: 0, Limit1: 1})
	// return tasksOp.ListTags(selectors.In("status", ""), &crud.GetOptions{Limit0: 0, Limit1: 1})
}
