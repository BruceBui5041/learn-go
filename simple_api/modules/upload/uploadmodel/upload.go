package uploadmodel

import (
	"fmt"
	"learn-go/simple_api/common"
	"strings"
)

func ErrFileIsNotImage(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("File %s is not an image", strings.ToLower("upload")),
		fmt.Sprintf("ErrFileIsNotImage%s", "upload"),
	)
}

func ErrCannotSaveFile(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("File %s is not an image", strings.ToLower("upload")),
		fmt.Sprintf("ErrCannotSaveFile%s", "upload"),
	)
}
