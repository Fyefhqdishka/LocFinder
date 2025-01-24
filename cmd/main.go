package main

import (
	"fmt"
	"github.com/Fyefhqdishka/LocFinder/internal/app"
	"github.com/Fyefhqdishka/LocFinder/internal/config"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("can't load config, err: %v", err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("can't load server, err: %v", err)
	}

	go func() {
		if err = app.Run(); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()
	log.Println("shutting down...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	if err := app.Stop(); err != nil {
		fmt.Errorf("error during shutdown: %v", err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("can't load env file, err=%v", err)
	}
}

func JoinChannels(s ...<-chan int) <-chan int {
	wg := &sync.WaitGroup{}
	result := make(chan int)

	wg.Add(len(s))
	for _, ch := range s {
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				result <- v
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}
