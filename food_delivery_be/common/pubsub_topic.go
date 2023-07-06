package common

import "learn-go/food_delivery_be/pubsub"

const (
	TopicUserLikeRestaurant    pubsub.Topic = "TopicUserLikeRestaurant"
	TopicUserDislikeRestaurant pubsub.Topic = "TopicUserDislikeRestaurant"
)
