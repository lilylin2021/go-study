package data

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

type BookRepo struct {
	db *gorm.DB
}

type BookPO struct {
	ID     uint64
	Name   string
	Title  string
	Author string
}

type BookDO struct {
	ID     uint64
	Name   string
	Title  string
	Author string
}

func NewBookRepo(db *gorm.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (db *BookRepo) GetBook(name string) (book *BookPO) {
	db.db.Where("name = ?", name).First(&book)
	return
}

func Convert2DO(book *BookPO) (bookDo *BookDO) {
	bytes, _ := json.Marshal(book)
	json.Unmarshal(bytes, bookDo)
	return
}
