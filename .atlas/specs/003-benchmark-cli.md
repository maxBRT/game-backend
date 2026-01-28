# Benchmark CLI

## 1. Description

A load-testing and demonstration tool built with Go that stress-tests both the matchmaking service and the player service store. It simulates hundreds or thousands of concurrent players joining queues and making store purchases. The CLI provides a real-time dashboard displaying performance metrics, making it ideal for showcasing the system's concurrency capabilities.

**Goal:** Provide a visually impressive demonstration of both services handling high-concurrency workloads—matchmaking (Go) and transactional store purchases (C#/SQL).

## 2. UI/UX

### Command Line Interface

#### `benchmark-cli run`
Runs a load test against both services with configurable parameters.

**Flags:**
- `--survivors <n>` - Number of simulated survivor players (default: 100)
- `--killers <n>` - Number of simulated killer players (default: 25)
- `--duration <duration>` - How long to run the benchmark (default: 30s)
- `--ramp-up <duration>` - Time to gradually spawn all players (default: 5s)
- `--purchases-per-sec <n>` - Target store purchase rate (default: 50)

**Example:**
```
benchmark-cli run --survivors 500 --killers 150 --purchases-per-sec 100 --duration 60s
```

**Real-time Dashboard Output:**
```
╔══════════════════════════════════════════════════════════════╗
║                    BACKEND BENCHMARK                         ║
╠══════════════════════════════════════════════════════════════╣
║  Duration: 00:34 / 01:00                                     ║
╠══════════════════════════════════════════════════════════════╣
║  MATCHMAKER                        │  STORE                  ║
║    Survivors waiting:  127         │    Purchases:    3,412  ║
║    Killers waiting:     31         │    Purchases/sec: 98.2  ║
║    Matches formed:     892         │    Avg latency:    24ms ║
║    Matches/sec:       26.4         │    Errors:            0 ║
║    Avg time-to-match: 142ms        │                         ║
╠══════════════════════════════════════════════════════════════╣
║  TOTALS                                                      ║
║    Active connections:   650                                 ║
║    Total requests:    21,833                                 ║
║    Error rate:          0.00%                                ║
╚══════════════════════════════════════════════════════════════╝
```

**Final Summary Output:**
```
Benchmark complete!

Matchmaker Results:
  Total matches formed:    1,847
  Peak matches/sec:         31.2
  Avg time-to-match:       138ms
  P95 time-to-match:       289ms
  P99 time-to-match:       412ms

Store Results:
  Total purchases:         5,892
  Peak purchases/sec:      102.4
  Avg latency:              22ms
  P95 latency:              48ms
  P99 latency:              71ms

Overall:
  Total requests:        48,042
  Error rate:              0.00%
```

## 3. Backend API

This component is a client, not a server. It consumes both service APIs:

### Player Service (Service A)
- `POST /api/store/purchase` - Execute a store purchase transaction

### Matchmaker Service (Service B)
- `POST /match/join` - Join matchmaking queue (with role: survivor/killer)
- `GET /match/status/{ticketId}` - Poll for match status

## 4. File System

### Technology Stack
- **Language:** Go (Golang)
- **CLI Framework:** cobra
- **HTTP Client:** net/http with connection pooling
- **TUI:** bubbletea or similar for real-time dashboard
- **Concurrency:** goroutines + channels for simulated players

### Directory Structure
```
/benchmark-cli
├── main.go
├── cmd/
│   ├── root.go           # Root command setup
│   └── run.go            # run command
├── internal/
│   ├── client/
│   │   ├── match.go      # Matchmaker service client
│   │   └── store.go      # Player service store client
│   ├── benchmark/
│   │   ├── runner.go     # Orchestrates simulated players
│   │   ├── player.go     # Simulated player goroutine (matchmaking)
│   │   ├── purchaser.go  # Simulated store purchaser goroutine
│   │   └── metrics.go    # Collects and aggregates stats
│   ├── dashboard/
│   │   └── dashboard.go  # Real-time TUI display
│   └── config/
│       └── config.go     # Service URLs, timeouts
└── go.mod
```

### Configuration
- `PLAYER_SERVICE_URL` (default: `http://localhost:5000`)
- `MATCHMAKER_SERVICE_URL` (default: `http://localhost:8000`)
