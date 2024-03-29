package pblocal

import (
	"context"
	"learn-go/food_delivery_be/appsocketio"
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/pubsub"
	"log"
	"sync"
)

// A pb run locally (in-mem)
// It has a queue (buffer channel) at it's core and many group of subscribers.
// Because we want to send a message with a specific topic for many subscribers in a group can handle.

type localPubSub struct {
	messageQueue chan *pubsub.Message
	mapChannel   map[pubsub.Topic][]chan *pubsub.Message
	locker       *sync.RWMutex
}

func NewPubSub(realtimeEngine appsocketio.RealtimeEngine) *localPubSub {
	pb := &localPubSub{
		messageQueue: make(chan *pubsub.Message, 10000),
		mapChannel:   make(map[pubsub.Topic][]chan *pubsub.Message),
		locker:       new(sync.RWMutex),
	}

	pb.run()

	return pb
}

func (pb *localPubSub) Publish(ctx context.Context, channel pubsub.Topic, data *pubsub.Message) error {
	data.SetChannel(channel)

	// Sủ dụng routine ở đây để đảm bảo không bị lock khi chờ channel nạp dữ liệu
	go func() {
		defer common.AppRecover()
		pb.messageQueue <- data
		log.Println("New channel published: ", data.String())
	}()

	return nil
}

func (ps *localPubSub) Subscribe(ctx context.Context, channel pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	c := make(chan *pubsub.Message)

	// Lock để đảm bảo nếu nhiều routine goi Subscribe cùng một lúc sẽ không crash
	ps.locker.Lock()
	if val, ok := ps.mapChannel[channel]; ok {
		val = append(ps.mapChannel[channel], c)
		ps.mapChannel[channel] = val
	} else {
		ps.mapChannel[channel] = []chan *pubsub.Message{c}
	}
	ps.locker.Unlock()

	return c, func() {
		log.Println("Unsubscribe")

		if chans, ok := ps.mapChannel[channel]; ok {
			for i := range chans {
				if chans[i] == c {
					chans = append(chans[:i], chans[i+1:]...)

					ps.locker.Lock()
					ps.mapChannel[channel] = chans
					ps.locker.Unlock()
					break
				}
			}
		}
	}

}

func (pb *localPubSub) run() error {
	log.Println("Pubsub started")

	go func() {
		for {
			msg := <-pb.messageQueue
			log.Println("Message dequeue:", msg)

			if subs, ok := pb.mapChannel[msg.Channel()]; ok {
				for i := range subs {
					go func(c chan *pubsub.Message) {
						c <- msg
					}(subs[i])
				}
			}
			//else {
			//	ps.messageQueue <- mess
			//}
		}
	}()

	return nil
}
