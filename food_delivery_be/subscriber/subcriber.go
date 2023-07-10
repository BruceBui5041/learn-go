package subscriber

import (
	"context"
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/component/asyncjob"
	"learn-go/food_delivery_be/pubsub"
	"log"
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, message *pubsub.Message) error
}

type consumerEngine struct {
	appCtx component.AppContext
}

func NewEngine(appCtx component.AppContext) *consumerEngine {
	return &consumerEngine{appCtx: appCtx}
}

func (engine *consumerEngine) Start() error {
	engine.startSubscribeTopic(
		common.TopicUserLikeRestaurant,
		true,
		RunIncreaseLikeCountAfterUserLikedRestaurant(engine.appCtx),
	)

	engine.startSubscribeTopic(
		common.TopicUserDislikeRestaurant,
		true,
		RunDecreaseLikeCountAfterUserUnlikedRestaurant(engine.appCtx),
	)
	return nil
}

func (engine *consumerEngine) startSubscribeTopic(topic pubsub.Topic, isConcurrent bool, consumerJobs ...consumerJob) error {
	c, _ := engine.appCtx.GetPubSub().Subscribe(context.Background(), topic)

	for _, item := range consumerJobs {
		log.Println("Setup consumer for: ", item.Title)
	}

	getJobHandler := func(job *consumerJob, message *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) error {
			log.Println("running job for", job.Title, ". Value: ", message.Data())
			return job.Hld(ctx, message)
		}
	}

	go func() {
		for {
			msg := <-c

			jobHandlerArray := make([]asyncjob.Job, len(consumerJobs))

			for i := range consumerJobs {
				jobHandlerArray[i] = asyncjob.NewJob(getJobHandler(&consumerJobs[i], msg))
			}

			groud := asyncjob.NewGroup(isConcurrent, jobHandlerArray...)

			if err := groud.Run(context.Background()); err != nil {
				log.Println(err)
			}

		}
	}()

	return nil
}
