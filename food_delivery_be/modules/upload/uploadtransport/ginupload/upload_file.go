package ginupload

import (
	"learn-go/food_delivery_be/common"
	"learn-go/food_delivery_be/component"
	"learn-go/food_delivery_be/modules/upload/uploadbusiness"
	"learn-go/food_delivery_be/modules/upload/uploadstorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Upload File to S3
// 1. Get image/file from request fileHeader
// 2. Check file is real image
// 3. Save image
// 4. Save to local service
// 5. Save to cloud storage (S3)
// 6. Improve security

func Upload(appCtx component.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()

		fileHeader, err := c.FormFile("file")

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		folder := c.DefaultPostForm("folder", "img")

		file, err := fileHeader.Open()

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		defer file.Close()

		dataBytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(dataBytes); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		imgStore := uploadstorage.NewSQLStore(db)
		biz := uploadbusiness.NewUploadBiz(appCtx.GetUploadProvider(), imgStore)

		img, err := biz.Upload(c.Request.Context(), dataBytes, folder, fileHeader.Filename)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(img, nil, nil))
	}
}
