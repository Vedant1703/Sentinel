
# Sentinel 🛡️  
### A Distributed Rate Limiting System built in Go

Sentinel is a **high-performance distributed rate limiter** designed to enforce **global request limits across multiple service instances**.  
It is built using **Golang**, focusing on **concurrency, scalability, and real-world backend system design**.

Sentinel mirrors the core ideas used in modern API gateways and large-scale systems to protect services from traffic spikes, abuse, and overload.

---

## ✨ Features

- Distributed rate limiting across multiple nodes
- Globally consistent limits using centralized coordination
- High concurrency support using Go primitives
- Multiple rate limiting algorithms
- Easy integration as middleware or standalone service
- Horizontally scalable design
- Production-inspired architecture

---

## 🧠 Supported Rate Limiting Algorithms

Sentinel currently supports:

- **Token Bucket**  
  Allows smooth traffic flow while supporting short bursts

- **Sliding Window**  
  Accurate request tracking over rolling time windows

- **Fixed Window**  
  Simple and fast rate limiting (used for comparison and benchmarking)

---

## 🏗️ System Architecture

```

Client
↓
Service / API Gateway
↓
Sentinel Middleware
↓
Distributed Store (Redis / etcd)
↓
Allow / Reject Decision

````

### Key Design Points
- Each Sentinel instance runs independently
- Shared state ensures **consistent rate limits across nodes**
- Atomic operations prevent race conditions
- Designed for horizontal scaling

---

## ⚙️ Tech Stack

- **Language:** Go (Golang)
- **Distributed Store:** Redis / etcd
- **API Interface:** REST / gRPC
- **Concurrency:** Goroutines, Mutexes, Atomic Operations
- **Deployment:** Docker (optional)

---

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or higher
- Redis or etcd running locally or remotely

---

### Clone the Repository

```bash
git clone https://github.com/your-username/sentinel.git
cd sentinel
````

---

### Run the Service

```bash
go run main.go
```

---

## 🔌 Usage Example

### Basic Rate Limiting

```go
limiter := sentinel.NewLimiter(sentinel.Config{
    Requests: 100,
    Window:   time.Minute,
    Strategy: sentinel.TokenBucket,
})

if !limiter.Allow(clientID) {
    return http.StatusTooManyRequests
}
```

---

## 🧩 Middleware Integration Example

```go
func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        clientID := r.RemoteAddr

        if !limiter.Allow(clientID) {
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

---

## 📊 Performance Considerations

* Uses atomic operations for safe concurrent access
* Minimizes network calls to the distributed store
* Designed to handle high QPS environments
* Benchmarked for latency and throughput (WIP)

---

## 🔐 Failure Handling

* Graceful degradation if the distributed store becomes unavailable
* Configurable fallback strategies
* Safe defaults to protect downstream services

---

## 📦 Project Structure

```
sentinel/
├── cmd/
│   └── server/
│       └── main.go
├── limiter/
│   ├── token_bucket.go
│   ├── sliding_window.go
│   └── fixed_window.go
├── store/
│   ├── redis.go
│   └── etcd.go
├── middleware/
│   └── http.go
├── config/
│   └── config.go
├── tests/
│   └── limiter_test.go
└── README.md
```

---

## 🧪 Testing

Run unit tests using:

```bash
go test ./...
```

---

## 🚧 Future Improvements

* Distributed leader election
* Rate limit dashboards and metrics
* Prometheus & Grafana integration
* Adaptive rate limiting
* gRPC interceptor support

---

## 👥 Team

This project was built collaboratively by a team of two, with responsibilities split across:

* Core rate limiting logic and algorithms
* Distributed coordination and API integration

---

## 📜 License

MIT License

---

## ⭐ Why Sentinel?

Sentinel is not a CRUD application.
It demonstrates **backend engineering fundamentals**, **distributed systems concepts**, and **production-grade design decisions**, making it an ideal project for learning and showcasing system-level expertise.

