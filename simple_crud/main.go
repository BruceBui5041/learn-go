package main

import (
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/**
	omitempty: if value is nil, 0, "" then the key will be removed
**/
type Note struct {
	Id      int    `json:"id,omitempty" gorm:"column:id"`
	Title   string `json:"title" gorm:"column:title"`
	Content string `json:"content" gorm:"column:content"`
}

func (Note) TableName() string {
	return "notes"
}

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	dsn := viper.GetString("DBConnectionStr")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	newNote := Note{Title: "Test Note", Content: "Test content"}
	db.Create(newNote)
}
