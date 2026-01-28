# Behaviour Interactive Backend Mock

A simplified backend implementation inspired by **Dead by Daylight**, the asymmetric multiplayer horror game by Behaviour Interactive. This project demonstrates core backend patterns for live-service games with role-based matchmaking (4 survivors vs 1 killer).

## Problem Statement

Asymmetric multiplayer games like Dead by Daylight require robust backend infrastructure that handles two fundamentally different concerns: persistent player data (inventory, economy, progression) and real-time game logic (role-based matchmaking, lobbies). These concerns have different scalability and consistency requirements, making a monolithic approach suboptimal.

## Proposed Solution

A microservices architecture that separates persistent data management from real-time matchmaking logic, using purpose-built technologies for each domain:

- **C# / ASP.NET Core** for ACID-compliant player data and transactions
- **Go (Golang)** for high-concurrency matchmaking with in-memory queues
- **Go CLI** for benchmarking and demonstrating system performance

## Core Features

### Player Data Management
- Player profile retrieval (level, XP, currency)
- Transactional store purchases with balance validation
- Inventory management

### Matchmaking System
- Role-based queue management (survivor queue + killer queue)
- Lobby creation when thresholds are met (4 survivors + 1 killer)
- Long-polling status checks for match readiness

### Benchmark CLI
- Stress-tests both matchmaking and store purchase endpoints
- Simulates hundreds/thousands of concurrent players and transactions
- Real-time TUI dashboard showing queue depths, matches/sec, purchases/sec, and latencies
- Final summary with P95/P99 latencies and throughput metrics for both services
- Designed for impressive live demonstrations showcasing Go and C#/SQL under load

## Technical Architecture

| Component | Role | Tech Stack | Description |
| :--- | :--- | :--- | :--- |
| **Service A** | **Player Data** | **C# / ASP.NET Core** | Handles persistent data, ACID transactions, and inventory management (SQL). |
| **Service B** | **Matchmaker** | **Go (Golang)** | Handles high-concurrency connections, in-memory queues, and matchmaking logic. |
| **Benchmark CLI** | **Load Testing** | **Go CLI** | Stress-tests both services (matchmaking + store). Displays real-time performance dashboard for demonstrations. |

## Data Model

### Player Service (SQL)
- **Player**: ID, Username, Level, XP, Currency
- **Inventory**: PlayerID, ItemID, Quantity, AcquiredAt
- **Item**: ID, Name, Price, Description

### Matchmaker Service (In-Memory)
- **SurvivorQueue**: List of player IDs queuing as survivors
- **KillerQueue**: List of player IDs queuing as killers
- **Match**: MatchID, SurvivorIDs[4], KillerID, Status, CreatedAt

## Project Structure

```text
/behaviour-backend-mock
├── /player-service      # C# ASP.NET Core API
│   ├── Controllers/
│   ├── Models/
│   └── Data/
├── /matchmaker-service  # Go Service
│   ├── main.go
│   └── internal/
│       └── queue/       # Matchmaking logic
├── /benchmark-cli       # Go Benchmark Tool
│   ├── main.go
│   └── internal/
│       ├── benchmark/   # Load test orchestration
│       └── dashboard/   # Real-time TUI
├── docker-compose.yml   # Orchestration
└── README.md
```

## Demonstration

The benchmark CLI showcases both services' performance with a real-time TUI:

```
benchmark-cli run --survivors 500 --killers 150 --purchases-per-sec 100 --duration 60s
```

**Real-time dashboard displays:**
- Matchmaker: queue depths, matches formed, matches/sec, time-to-match latency
- Store: total purchases, purchases/sec, transaction latency
- Overall: active connections, total requests, error rate

**Final summary includes:**
- Peak throughput for both services
- P95/P99 latency percentiles for matchmaking and store transactions
- Total requests and error rate

This provides a visually impressive demonstration of Go's concurrency (matchmaking) and C#/SQL's transactional reliability (store) working together under heavy load.
