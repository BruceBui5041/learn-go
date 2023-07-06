package asyncjob

import (
	"context"
	"log"
	"sync"
)

type group struct {
	isConcurrent bool
	jobs         []Job
	wg           *sync.WaitGroup
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup),
	}

	return g
}

func (g *group) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs))

	errChan := make(chan error, len(g.jobs))

	for _, job := range g.jobs {

		if g.isConcurrent {
			// Do this instead
			go func(aj Job) {
				errChan <- g.runJob(ctx, aj)
				g.wg.Done()
			}(job)

			continue
		}

		job := job
		errChan <- g.runJob(ctx, job)
		g.wg.Done()
	}

	var err error

	for i := 1; i < len(g.jobs); i++ {
		if v := <-errChan; v != nil {
			err = v
		}
	}

	g.wg.Wait()
	return err
}

// retry
func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			log.Print(err)
			if j.State() == StateRetryFailed {
				return err
			}

			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}

	return nil
}
