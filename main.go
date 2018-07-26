package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WuHan0608/webhook-listener/webhook"
)

var (
	rootMux  = http.NewServeMux()
	httpSrv  = &http.Server{}
	stopChan = make(chan struct{}, 1)
)

func startListener() {
	rootMux.Handle("/webhook/dockerhub", http.StripPrefix("/webhook", webhook.DockerHubHandler()))
	httpSrv.Addr = os.Getenv("LISTEN_ADDRESS")
	httpSrv.Handler = rootMux

	go func() {
		<-stopChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatalf("http server shutdown error: %v", err)
		}
		log.Println("shutdown http server gracefully")
	}()

	log.Printf("Starting http server on %s\n", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		switch <-sigChan {
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
			signal.Stop(sigChan)
			close(stopChan)
		}
	}()

	startListener()
}
