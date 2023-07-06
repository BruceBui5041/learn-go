package subscriber

import (
	"context"
	"learn-go/food_delivery_be/component"
)

func SetUp(appCtx component.AppContext) {
	IncreaseLikeCountAfterUserLikedRestaurant(appCtx, context.Background())
}
