package common

import (
	"errors"
	"log"
)

var (
	RecordNotFound = errors.New("record not found")
)

// NOTE: Hàm này để recover những panic trong go routine
func AppRecover() {
	if err := recover(); err != nil {
		log.Println("Recover error:", err)
	}
}
