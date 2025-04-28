package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	mutex sync.Mutex
	state string
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	go func() {
		defer close(done)
		if err := runOperations(ctx); err != nil {
			log.Printf("Erro: %v", err)
		}
	}()

	select {
	case <-sigChan:
		fmt.Println("Interrupção recebida, limpando...")
		cancel()
		<-done
	case <-done:
		fmt.Println("Operações finalizadas")
	}
}

func runOperations(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	var (
		dbErr   error
		httpErr error
	)

	go func() {
		defer wg.Done()
		dbErr = dbConnection(ctx)
	}()

	go func() {
		defer wg.Done()
		httpErr = httpConnection(ctx)
	}()

	wg.Wait()

	if dbErr != nil {
		return fmt.Errorf("DB: %w", dbErr)
	}
	if httpErr != nil {
		return fmt.Errorf("HTTP: %w", httpErr)
	}

	return nil
}

func dbConnection(ctx context.Context) error {
	duration := time.Duration(rand.Intn(700)) * time.Millisecond

	select {
	case <-time.After(duration):
		mutex.Lock()
		state = "DB update"
		mutex.Unlock()

		fmt.Printf("Conexão Database finalizada em %v\n", duration)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("Conexão DB cancelada: %w", ctx.Err())
	}
}

func httpConnection(ctx context.Context) error {
	duration := time.Duration(rand.Intn(700)) * time.Millisecond

	select {
	case <-time.After(duration):
		mutex.Lock()
		state = "HTTP call done"
		mutex.Unlock()

		fmt.Printf("Conexão HTTP finalizada em %v\n", duration)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("Conexão HTTP cancelada: %w", ctx.Err())
	}
}
