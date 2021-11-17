// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"go-study/book-manage/internal/biz"
	"go-study/book-manage/internal/data"
	"go-study/book-manage/internal/service"
)

// Injectors from wire.go:

func InitBookService() *service.BookService {
	db := data.NewDB()
	bookRepo := data.NewBookRepo(db)
	bookBiz := biz.NewBookBiz(bookRepo)
	bookService := service.NewBookService(bookBiz)
	return bookService
}
