package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type LogEntry struct {
	Clock   time.Time
	Message string
}

type LogReader struct {
	Logs []LogEntry
	mu   sync.Mutex
}

func (lr *LogReader) ReadLogs() <-chan LogEntry {
	logCh := make(chan LogEntry)
	go func() {
		defer close(logCh)
		lr.mu.Lock()
		defer lr.mu.Unlock()
		for _, entry := range lr.Logs {
			logCh <- entry
		}
	}()
	return logCh
}

type LogWriter struct{}

func (lw *LogWriter) WriteLog(entry LogEntry) error {
	fmt.Printf("Writing log: %v\n", entry)
	return nil
}

func handleLogEntry(entry LogEntry, writer *LogWriter, ctx context.Context, workerID int) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		if err := writer.WriteLog(entry); err != nil {
			return fmt.Errorf("worker %d: %w", workerID, err)
		}
		return nil
	}
}

func startWorkers(logCh <-chan LogEntry, writer *LogWriter, ctx context.Context, numWorkers int) *errgroup.Group {
	var eg errgroup.Group
	for i := 0; i < numWorkers; i++ {
		workerID := i
		eg.Go(func() error {
			for {
				select {
				case entry, ok := <-logCh:
					if !ok {
						return nil
					}
					if err := handleLogEntry(entry, writer, ctx, workerID); err != nil {
						return err
					}
				case <-ctx.Done():
					return nil
				}
			}
		})
	}
	return &eg
}

func startLogProcessing() error {
	reader := &LogReader{
		Logs: []LogEntry{
			{Clock: time.Now(), Message: "Log 1"},
			{Clock: time.Now(), Message: "Log 2"},
		},
	}

	writer := &LogWriter{}

	logCh := reader.ReadLogs()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numWorkers := 4

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	eg := startWorkers(logCh, writer, ctx, numWorkers)

	go func() {
		select {
		case <-ctx.Done():
		case <-sigCh:
			cancel()
		}
	}()

	return eg.Wait()
}

func main() {
	if err := startLogProcessing(); err != nil {
		log.Printf("Error: %v\n", err)
	}
}
