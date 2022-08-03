package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// omitempty: if value is nil, 0, "" then the key will be removed
type Note struct {
	Id      int    `json:"id,omitempty" gorm:"column:id"`
	Title   string `json:"title" gorm:"column:title"`
	Content string `json:"content" gorm:"column:content"`
}

// *string: so now we can update title = ""
type NoteUpdate struct {
	Title *string `json:"title" gorm:"column:title"`
}

// TableName overrides the table name used by Note to `notes`
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

	db.AutoMigrate(&Note{})
	/**
		Insert
	**/
	// newNote := Note{Title: "Test Note 3", Content: "Test content 3"}

	// if err := db.Create(&newNote); err != nil {
	// 	fmt.Println(err)
	// }

	var notes []Note
	db.Where("content IS NOT NULL").Find(&notes)

	var note Note
	if err := db.Where("id = ?", 3).First(&note); err == nil {
		log.Println(err)
	}

	newTitle := ""
	db.Updates(&NoteUpdate{Title: &newTitle})

	db.Table(Note{}.TableName()).Where("id = ?", 5).Delete(nil)

	fmt.Println(notes)
}
