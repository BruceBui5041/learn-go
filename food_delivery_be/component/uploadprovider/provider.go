package uploadprovider

import (
	"context"
	"learn-go/food_delivery_be/common"
)

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, data []byte, dst string) (*common.Image, error)
}
