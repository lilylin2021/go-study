//go:build wireinject
// +build wireinject

package main

import (
	"go-study/book-manage/internal/biz"
	"go-study/book-manage/internal/data"
	"go-study/book-manage/internal/service"

	"github.com/google/wire"
)

func InitBookService() *service.BookService {
	wire.Build(service.NewBookService, biz.NewBookBiz, data.NewBookRepo, data.NewDB)
	return &service.BookService{}
}
