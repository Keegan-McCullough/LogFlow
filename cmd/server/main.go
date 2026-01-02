package main

import (
	"fmt"
	"sync"
	"time"
)

type LogEntry struct {
	Message string
	Level   string
	Source  string
}

type Metrics struct {
	mu          sync.Mutex
	ErrorCounts map[string]float64
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

	go producer(ch)

	wg.Add(2)
	go consumer(ch, metrics, &wg)
	go consumer(ch, metrics, &wg)

	wg.Wait()

	fmt.Println("-----------------")
	fmt.Println("Final stats:", metrics.ErrorCounts)

}
