package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MrDavudov/AirtableToDraw/pkg/handler"
)


func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	
	srv := new(handler.Server)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("errors occured while running http server: %s", err)
	}
}