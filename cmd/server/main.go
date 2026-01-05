package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type LogEntry struct {
	Message string `json:"message"`
	Level   string `json:"level"`
	Source  string `json:"source"`
}

type Metrics struct {
	mu          sync.Mutex
	ErrorCounts map[string]float64
}

func handleLog(ch chan<- LogEntry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Only allow POST requests
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		// 2. Decode JSON payload
		var entry LogEntry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			http.Error(w, "Bad JSON", http.StatusBadRequest)
			return
		}

		// 3. Push to Channel (Non-blocking if buffer has space)
		select {
		case ch <- entry:
			w.WriteHeader(http.StatusAccepted) // 202 Accepted
			w.Write([]byte("Log ingested"))
		default:
			// If buffer is full, drop the log so we don't crash the server
			http.Error(w, "Queue full", http.StatusServiceUnavailable)
		}
	}
}

func (m *Metrics) IncCounter(source string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorCounts[source]++
}

func producer(ch chan<- LogEntry) {
	sources := []string{"Auth", "Payment", "Auth", "Database", "Payment"}
	levels := []string{"INFO", "ERROR", "INFO", "ERROR", "ERROR"}

	for i := 0; i < 5; i++ {
		entry := LogEntry{
			Level:   levels[i],
			Source:  sources[i],
			Message: fmt.Sprint(" Something happened %d", i),
		}
		fmt.Printf("->Producing: [%s] %s %s\n", entry.Level, entry.Source, time.StampMilli)
		ch <- entry
		time.Sleep(200 * time.Millisecond)
	}
	close(ch) // signal no more values
}

func consumer(ch <-chan LogEntry, metrics *Metrics, wg *sync.WaitGroup) {
	defer wg.Done()

	for entry := range ch {
		time.Sleep(300 * time.Millisecond)
		if entry.Level == "ERROR" {
			fmt.Printf(" <- Worker processing ERROR from %s %s\n", entry.Source, time.StampMilli)
			metrics.IncCounter(entry.Source)
		}
	}
}

func main() {
	ch := make(chan LogEntry, 100)
	metrics := &Metrics{ErrorCounts: make(map[string]float64)}

	var wg sync.WaitGroup

	// 1. Start Workers (Consumers)
	// These sit in the background waiting for work
	wg.Add(2)
	go consumer(ch, metrics, &wg)
	go consumer(ch, metrics, &wg)

	// 2. Start HTTP Server
	// This replaces the old 'producer()' function
	http.HandleFunc("/log", handleLog(ch))

	fmt.Println("Server started on port 8080...")

	// ListenAndServe blocks forever, so we don't need wg.Wait() right now
	// (In Phase 4 we will handle graceful shutdown)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server failed:", err)
	}
}
