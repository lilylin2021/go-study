package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func newServer() *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })

	return &http.Server{
		Handler: router,
		Addr:    ":8080",
	}
}

func main() {
	errg, ctx := errgroup.WithContext(context.Background())

	server := newServer()

	errg.Go(func() error {
		server.ListenAndServe()
	})

	errg.Go(func() error {
		<-ctx.Done

		timeOutCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		fmt.Println("start shut down http server......")
		return server.Shutdown(timeOutCtx)
	})

	errg.Go(func() error {
		quit := make(chan os.Signal, 0)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGABRT)

		select {
		case sig := <-quit:
			fmt.Println("receive signal, quit.....")
			return errors.Errorf("get os signal: %v", sig)
		}
	})

	errg.Wait()
}
