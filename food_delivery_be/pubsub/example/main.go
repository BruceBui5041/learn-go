package main

import (
	"context"
	"learn-go/food_delivery_be/pubsub"
	"learn-go/food_delivery_be/pubsub/pblocal"
	"log"
	"time"
)

func main() {
	var localPb pubsub.Pubsub = pblocal.NewPubSub()
	topic := pubsub.Topic("OrderCreated")

	con1, close1 := localPb.Subscribe(context.Background(), topic)
	con2, _ := localPb.Subscribe(context.Background(), topic)

	localPb.Publish(context.Background(), topic, pubsub.NewMessage(1))
	localPb.Publish(context.Background(), topic, pubsub.NewMessage(2))

	go func() {
		for {
			log.Println("Con1:", (<-con1).Data())
			time.Sleep(time.Millisecond * 400)
		}
	}()

	go func() {
		for {
			log.Println("Con2:", (<-con2).Data())
			time.Sleep(time.Millisecond * 400)
		}
	}()

	time.Sleep(time.Second * 3)
	close1()
	// close2()

	localPb.Publish(context.Background(), topic, pubsub.NewMessage(3))

	time.Sleep(time.Second * 2)
}
