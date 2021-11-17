package biz

import (
	"errors"
	"go-study/book-manage/internal/data"
)

type BookBiz struct {
	repo *data.BookRepo
}

func NewBookBiz(repo *data.BookRepo) *BookBiz {
	return &BookBiz{repo: repo}
}

func (manager *BookBiz) GetBook(name string) (*data.BookDO, error) {
	if len(name) == 0 {
		return nil, errors.New("please input valid book name")
	}

	return data.Convert2DO(manager.repo.GetBook(name)), nil
}
