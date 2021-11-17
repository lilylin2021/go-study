package main

import (
	"geekbang-go-camp/homework/fourth/api"

	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/core/elog"
	"github.com/gotomicro/ego/server/egrpc"
)

func main() {
	if err := ego.New().Serve(func() *egrpc.Component {
		srv := egrpc.Load("server.http").Build()
		api.RegisterUserServiceServer(srv, InitBookService())
		return srv
	}()).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
