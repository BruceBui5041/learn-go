package subscriber

import (
	"context"
	"learn-go/simple_api/component"
)

func SetUp(appCtx component.AppContext) {
	IncreaseLikeCountAfterUserLikedRestaurant(appCtx, context.Background())
}
