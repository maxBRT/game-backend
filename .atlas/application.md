# Behaviour Interactive Backend Mock

## Problem Statement

Live-service multiplayer games require robust backend infrastructure that handles two fundamentally different concerns: persistent player data (inventory, economy, progression) and real-time game logic (matchmaking, lobbies). These concerns have different scalability and consistency requirements, making a monolithic approach suboptimal.

## Proposed Solution

A microservices architecture that separates persistent data management from real-time matchmaking logic, using purpose-built technologies for each domain:

- **C# / ASP.NET Core** for ACID-compliant player data and transactions
- **Go (Golang)** for high-concurrency matchmaking with in-memory queues
- **Go CLI** for simulating game client behavior during development

## Core Features

### Player Data Management
- Player profile retrieval (level, XP, currency)
- Transactional store purchases with balance validation
- Inventory management

### Matchmaking System
- Player queue management using goroutines
- Lobby creation when queue threshold is met (4 players)
- Long-polling status checks for match readiness

### Developer Tooling
- CLI tool to simulate player login flow
- CLI tool to test matchmaking queue behavior

## Technical Architecture

| Component | Role | Tech Stack | Description |
| :--- | :--- | :--- | :--- |
| **Service A** | **Player Data** | **C# / ASP.NET Core** | Handles persistent data, ACID transactions, and inventory management (SQL). |
| **Service B** | **Matchmaker** | **Go (Golang)** | Handles high-concurrency connections, in-memory queues, and matchmaking logic. |
| **Client** | **Game Client** | **Go CLI** | A terminal tool that simulates a game client connecting to both services. |

## Data Model

### Player Service (SQL)
- **Player**: ID, Username, Level, XP, Currency
- **Inventory**: PlayerID, ItemID, Quantity, AcquiredAt
- **Item**: ID, Name, Price, Description

### Matchmaker Service (In-Memory)
- **LobbyQueue**: List of player IDs waiting for match
- **Match**: MatchID, PlayerIDs[], Status, CreatedAt

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
├── /game-client         # Go CLI Tool
│   └── main.go
├── docker-compose.yml   # Orchestration 
└── README.md
```
