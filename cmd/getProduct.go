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

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	<-quit

	log.Println("Shutdown GetProduct...")

}
