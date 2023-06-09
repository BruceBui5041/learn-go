package component

import (
	"learn-go/simple_api/component/uploadprovider"

	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetUploadProvider() uploadprovider.UploadProvider
	SecretKey() string
}

type appCtx struct {
	db             *gorm.DB
	uploadProvider uploadprovider.UploadProvider
	jwtSecretKey   string
}

func NewAppContext(db *gorm.DB, uploadProvider uploadprovider.UploadProvider, jwtSecretKey string) *appCtx {
	return &appCtx{db: db, uploadProvider: uploadProvider, jwtSecretKey: jwtSecretKey}
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
