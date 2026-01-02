# LogFlow 🌊

**A high-performance, concurrent log aggregation service built in Go.**

LogFlow is a backend system designed to ingest, process, and analyze log data in real-time. It leverages Go's native concurrency primitives (Goroutines and Channels) to handle high-throughput log streams without blocking, making it suitable for simulating production-level metrics collection.

---

## 🚀 Key Features

* **Concurrent Processing:** Uses a worker pool pattern to process incoming logs in parallel.
* **Buffered Pipeline:** Implements buffered channels to decouple log ingestion (Producer) from processing (Consumer), handling backpressure gracefully.
* **Thread-Safe State:** Uses `sync.Mutex` to aggregate metrics safely across multiple workers preventing race conditions.
* **Scalable Architecture:** Designed to easily scale the number of consumers based on system load.

---

## 🛠️ Architecture

The system follows a classic **Producer-Consumer** model:

1.  **Ingestion (Producer):** Simulates a high-speed log stream.
2.  **The Pipe (Channel):** A buffered Go Channel (`chan LogEntry`) acts as a thread-safe queue.
3.  **Processing (Consumers):** Multiple worker Goroutines pull logs from the channel concurrently.
4.  **Aggregation (State):** Workers update a shared `Metrics` struct, protected by a Mutex.

```mermaid
graph LR
    A[Log Producer] -->|Push Log| B(Buffered Channel)
    B -->|Pull Log| C[Worker 1]
    B -->|Pull Log| D[Worker 2]
    B -->|Pull Log| E[Worker 3]
    C --> F{Shared Metrics}
    D --> F
    E --> F
