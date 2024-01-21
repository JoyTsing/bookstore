package main

import (
	_ "bookstore/internal/store"
	"bookstore/server"
	"bookstore/store/factory"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s, err := factory.New("mem")
	if err != nil {
		panic(err)
	}
	srv := server.NewBookStoreServer(":8080", s)
	errCh, err := srv.ListenAndServe()
	if err != nil {
		log.Println("web server start failed:", err)
		return
	}
	log.Println("web server start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-errCh:
		log.Println("web server run failed", err)
	case <-c:
		log.Println("bookstore program is exiting...")
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		err = srv.Shutdown(ctx)
	}

	if err != nil {
		log.Println("bookstore program exit error:", err)
		return
	}
	log.Println("bootstore program exit.")
}
