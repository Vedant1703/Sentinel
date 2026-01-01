# Sentinel — Distributed Rate Limiter in Go

Sentinel is a production-grade distributed rate limiter built in Go that enforces global request limits across multiple instances using Redis and atomic Lua scripts, while maintaining low latency through a local in-memory burst limiter.

## 🚀 Why This Project?

Rate limiting is a critical requirement in real-world backend systems to protect services from abuse, ensure fair usage, and maintain system stability under high traffic. However, building a correct rate limiter is non-trivial due to challenges such as:

- Multiple application instances running concurrently
- Race conditions while enforcing limits
- Different rate limits for different endpoints
- Handling traffic spikes without degrading core services
- Observing and verifying limiter behavior in production

Simple in-memory limiters fail in distributed environments, while naive shared counters introduce correctness and performance issues.

Sentinel addresses these challenges by combining a local in-memory burst limiter with a Redis-backed global rate limiter using atomic Lua scripts, enabling correct, scalable, and observable rate limiting suitable for production systems.

## 🧠 Core Ideas & Key Features

- **Distributed Rate Limiting**  
  Enforces request limits globally across multiple application instances using Redis as a single source of truth.

- **Local Burst Protection**  
  An in-memory burst limiter absorbs short traffic spikes, reducing load on Redis and improving overall latency.

- **Atomic Enforcement with Redis + Lua**  
  Rate limit counters are updated using Redis Lua scripts to guarantee atomicity and eliminate race conditions.

- **Per-Endpoint Configuration**  
  Supports different rate limits for different routes (e.g., strict limits for authentication endpoints and higher limits for read-heavy APIs).

- **Flexible Key Strategy**  
  Rate limiting can be applied per user or per IP address, enabling fair usage for both authenticated and anonymous traffic.

- **Failure-Aware Design**  
  Gracefully handles Redis outages with configurable fail-open or fail-closed behavior.

- **Observability Built-In**  
  Exposes metrics and logs for allowed requests, blocked requests, and Redis failures to enable monitoring and debugging.

- **High-Concurrency Tested**  
  Validated under concurrent load using stress testing tools to ensure correctness and performance.

## 🏗️ Architecture Overview

Sentinel follows a layered middleware-based architecture that separates concerns while ensuring high performance and correctness in a distributed setup.

Client
↓
HTTP Middleware
→ Local Burst Limiter (in-memory)
→ Global Rate Limiter (Redis + Lua)
→ Route Configuration Lookup
→ Metrics & Logging
↓
Application Handler

- **HTTP Middleware** serves as the entry point for all requests and orchestrates rate-limiting decisions.
- **Burst Limiter** provides fast, in-memory protection against short-lived traffic spikes and shields Redis from excessive load.
- **Redis Global Limiter** enforces strict rate limits across all instances using atomic Lua scripts.
- **Configuration Layer** determines which rate limit policy applies to each request based on the endpoint.
- **Metrics and Logs** capture decision outcomes and system behavior for observability.

## 🔄 Request Flow

Each incoming HTTP request passes through a well-defined sequence of checks before reaching the application handler:

1. The request enters the HTTP middleware.
2. A rate-limiting key is extracted based on request context (user ID if available, otherwise IP address).
3. The local in-memory burst limiter checks for short-term traffic spikes.
4. The request path is matched against configured route policies to determine applicable limits.
5. The Redis-backed global rate limiter executes an atomic Lua script to enforce the rate limit.
6. Based on the limiter result, the request is either:
   - Allowed and forwarded to the application handler, or
   - Rejected with HTTP `429 Too Many Requests`.
7. Metrics counters and logs are updated to record the decision and system state.

This layered flow ensures low latency for normal traffic while maintaining strict global correctness under high concurrency.

## 🧮 Algorithms Used

Sentinel combines two complementary rate-limiting techniques to achieve both performance and correctness in a distributed environment.

### 1️⃣ Local Burst Limiter (In-Memory Sliding Window)

- Maintains a sliding window of recent request timestamps per key.
- Allows short traffic bursts while enforcing a maximum request count within a small time window.
- Executes entirely in memory, resulting in extremely low latency.
- Acts as a protective layer that reduces load on the Redis global limiter.

This limiter is optimized for fast rejection of excessive bursts on a single instance.

### 2️⃣ Global Rate Limiter (Redis Fixed Window with Lua)

- Uses Redis as a centralized store for request counters.
- Implements fixed-window rate limiting.
- Enforces limits using a Lua script that performs `GET`, `INCR`, and `EXPIRE` operations atomically.
- Guarantees correctness and eliminates race conditions across multiple application instances.

The combination of local burst limiting and atomic global enforcement provides a balanced trade-off between performance and strong consistency.

## ⚙️ Configuration

Sentinel supports configurable rate-limit policies to accommodate different traffic patterns across endpoints without requiring code changes.

Rate limits are defined using a simple rule-based configuration:

```go
cfg := config.Config{
    Routes: map[string]config.Rule{
        "/login":  {Limit: 5, Window: time.Minute},
        "/search": {Limit: 50, Window: time.Minute},
    },
    Default: config.Rule{
        Limit: 100,
        Window: time.Minute,
    },
}
```
Configuration Capabilities
Per-Endpoint Limits
Assign strict limits to sensitive routes (e.g., authentication) and more lenient limits to read-heavy APIs.

Default Fallback Rule
Ensures all routes are protected even if no explicit rule is defined.

Prefix-Based Matching
Supports scalable policies such as /api/*, with more specific routes taking precedence.

Extensible Design
The configuration structure can be easily extended to support role-based or API-key-based limits.

This approach allows Sentinel to adapt to real-world usage patterns without redeploying the application.

## 🔑 Key Strategy

Sentinel applies rate limits based on a dynamically derived request identity to ensure fair usage.

- **User-Based Limiting**  
  If a request includes an `X-User-ID` header, rate limits are applied per user:
user:<id>

- **IP-Based Limiting**  
For anonymous requests, limits are enforced per client IP:
ip:<address>

- **Namespaced Keys**  
Keys are prefixed (e.g., `user:` or `ip:`) to avoid collisions and allow future extensions such as API keys or service accounts.

This strategy supports both authenticated and unauthenticated traffic while remaining extensible and production-safe.

## 📊 Observability

Sentinel exposes internal state to make rate-limiting behavior visible and debuggable in production.

- **Metrics**  
  Atomic counters track allowed requests, blocked requests, and Redis errors.  
  Metrics are exposed via a lightweight HTTP endpoint:

GET /metrics

- **Logging**  
The system logs key events including:
- Requests blocked by the burst limiter or Redis limiter
- Request key and endpoint
- Redis failures

These signals allow operators to validate limiter behavior, detect anomalies, and understand traffic patterns under load.

## ⚠️ Failure Handling

Sentinel is designed to handle Redis failures gracefully.

- **Fail-Open (Default)**  
  Requests are allowed if Redis is unavailable to preserve availability.

- **Fail-Closed**  
  Requests are blocked during Redis outages for stricter protection.

The strategy is configurable, allowing trade-offs between availability and safety depending on system requirements.

## 🚀 Performance & Load Testing

The system was validated under concurrent load to ensure correctness and low latency.

Example test:
```bash
ab -n 1000 -c 50 http://localhost:8080/login
Observed behavior:

Strict enforcement of configured limits

Consistent HTTP 429 responses when limits are exceeded

Stable latency under high concurrency

```
## 📈 Verified Results

After load testing the `/login` endpoint:

allowed_requests 5
blocked_requests 995
redis_errors 0

These metrics confirm correct global enforcement, accurate blocking behavior, and error-free Redis interaction under concurrent load.

## ▶️ How to Run

Start Redis using Docker Compose:

```bash
docker-compose up
Run the application:

bash
Copy code
go run ./cmd/server
The server listens on:

arduino
Copy code
http://localhost:8080
```
## 📁 Project Structure

cmd/server → application entry point
config/ → rate-limit configuration
limiter/
├── burst/ → in-memory burst limiter
└── redis/ → Redis Lua global limiter
middleware/ → HTTP middleware
metrics/ → observability counters
redis/ → Lua scripts

## 🎯 Key Learnings

- Designing distributed systems with strong correctness guarantees
- Preventing race conditions using Redis and atomic Lua scripts
- Building low-latency middleware in Go
- Handling failures and trade-offs between availability and safety
- Adding observability to backend systems through metrics and logs

## 🧾 License

MIT
