package service

import (
	"context"
	"go-study/book-manage/api"
	"go-study/book-manage/internal/biz"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookService struct {
	api.UnimplementedBookManageServer
	biz *biz.BookBiz
}

func NewBookService(biz *biz.BookBiz) *BookService {
	return &BookService{biz: biz}
}

func (svc *BookService) GetBook(ctx context.Context, request *api.QueryBookRequest) (*api.QueryBookResponse, error) {
	book, err := svc.biz.GetBook(request.Name)
	if err != nil {
		return nil, status.Error(codes.Code(5), "system error")
	}

	return &api.QueryBookResponse{
		Book: &api.Book{
			Bid:    book.ID,
			Name:   book.Name,
			Title:  book.Title,
			Author: book.Author,
		},
	}, nil
}
