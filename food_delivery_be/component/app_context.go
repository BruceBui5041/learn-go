package component

import (
	"learn-go/food_delivery_be/component/uploadprovider"
	"learn-go/food_delivery_be/pubsub"

	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetUploadProvider() uploadprovider.UploadProvider
	SecretKey() string
	GetPubSub() pubsub.Pubsub
}

type appCtx struct {
	db             *gorm.DB
	uploadProvider uploadprovider.UploadProvider
	jwtSecretKey   string
	pb             pubsub.Pubsub
}

func NewAppContext(db *gorm.DB, uploadProvider uploadprovider.UploadProvider, jwtSecretKey string, pb pubsub.Pubsub) *appCtx {
	return &appCtx{db: db, uploadProvider: uploadProvider, jwtSecretKey: jwtSecretKey, pb: pb}
}

func (ctx *appCtx) SecretKey() string {
	return ctx.jwtSecretKey
}

func (ctx *appCtx) GetMainDBConnection() *gorm.DB {
	return ctx.db
}

func (ctx *appCtx) GetUploadProvider() uploadprovider.UploadProvider {
	return ctx.uploadProvider
}

func (ctx *appCtx) GetPubSub() pubsub.Pubsub {
	return ctx.pb
}
