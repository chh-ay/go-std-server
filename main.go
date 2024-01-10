package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type serverKey int

const (
	keyServerAddr serverKey = iota
)

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Println(ctx.Value(keyServerAddr), "Request from:", r.URL.Path)
	w.Write([]byte("Hello World"))
}

func getPong(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Println(ctx.Value(keyServerAddr), "Request from:", r.URL.Path)
	w.Write([]byte("Pong!"))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", getHello)
	mux.HandleFunc("/ping", getPong)

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx() // in case func return early

	server := &http.Server{
		Addr:    ":9090",
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server closed")
	} else if err != nil {
		fmt.Println("Error listening for server:", err)
	}
}
