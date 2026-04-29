package main

import (
	"fmt"
	"gocasts/gameapp/config"
	"gocasts/gameapp/scheduler"
	"os"
	"os/signal"
	"time"
)

func main() {

	cfg := config.Load("config.yml")

	done := make(chan bool)

	go func() {
		sch := scheduler.New()
		sch.Start(done)
	}()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch
	fmt.Println("arrived interupt signal, shutting down gracefully")

	done <- true

	time.Sleep(cfg.Application.GracefulShutdownTimeout)
}
