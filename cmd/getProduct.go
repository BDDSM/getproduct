package main

import (
	"github.com/korableg/getproduct/internal/api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	log.Println("Starting GetProduct...")

	api.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Println("GetProduct has started")

	<-quit

	log.Println("Shutdown GetProduct...")

}
