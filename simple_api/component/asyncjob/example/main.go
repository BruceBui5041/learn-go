package main

import (
	"context"
	"learn-go/simple_api/component/asyncjob"
	"log"
	"time"
)

func main() {
	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		log.Println("I am job 1")
		time.Sleep(time.Second)
		// return errors.New("err of job 1")
		return nil
	})

	job2 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 2)
		log.Println("I am job 2")
		// return errors.New("err of job 2")
		return nil
	})

	job3 := asyncjob.NewJob(func(ctx context.Context) error {
		time.Sleep(time.Second * 3)
		log.Println("I am job 3")
		return nil
	})

	// job2.SetRetryDurations([]time.Duration{time.Second * 2})

	group := asyncjob.NewGroup(true, job1, job2, job3)
	err := group.Run(context.Background())

	log.Println("Group result:", err)

	// if err := job1.Execute(context.Background()); err != nil {
	// 	log.Println("Job 1 err:", err)

	// 	for {
	// 		if err := job1.Retry(context.Background()); err == nil {
	// 			break
	// 		}

	// 		log.Println("Job 1 retry err:", err)

	// 		if job1.State() == asyncjob.StateRetryFailed {
	// 			break
	// 		}
	// 	}
	// }
}
