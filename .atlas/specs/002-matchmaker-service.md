# Matchmaker Service

## 1. Description

The high-performance matchmaking service built with Go. This service groups players into lobbies using goroutines and in-memory queues for maximum concurrency and low latency.

**Goal:** Efficiently match players into game lobbies when enough players are queued.

**Matching Rule:** When `len(queue) >= 4`, create a match and remove those players from the queue.

## 2. UI/UX

N/A - This is a backend API service with no user interface.

## 3. Backend API

### Endpoints

#### `POST /match/join`
Add a player to the matchmaking queue.

**Request:**
```json
{
  "playerId": "string",
  "playerName": "string"
}
```

**Response:**
```json
{
  "success": true,
  "queuePosition": 0,
  "ticketId": "string"
}
```

#### `GET /match/status/{ticketId}`
Long-polling endpoint to check if a match is ready.

**Response (Waiting):**
```json
{
  "status": "waiting",
  "queuePosition": 2
}
```

**Response (Matched):**
```json
{
  "status": "matched",
  "matchId": "string",
  "players": [
    {"id": "string", "name": "string"},
    {"id": "string", "name": "string"},
    {"id": "string", "name": "string"},
    {"id": "string", "name": "string"}
  ]
}
```

### Matchmaking Logic

1. Player joins queue via `/match/join`
2. Background goroutine continuously polls the queue
3. When queue size reaches 4:
   - Create new match with unique ID
   - Assign first 4 players to the match
   - Remove players from queue
   - Update match status for those players
4. Players polling `/match/status/{ticketId}` receive match info

## 4. File System

### Technology Stack
- **Language:** Go (Golang)
- **Concurrency:** Goroutines + Channels
- **Storage:** In-memory (maps with mutex protection)

### Directory Structure
```
/matchmaker-service
├── main.go
├── internal/
│   ├── queue/
│   │   ├── queue.go        # Queue data structure
│   │   └── matcher.go      # Matching logic goroutine
│   ├── handlers/
│   │   ├── join.go         # POST /match/join handler
│   │   └── status.go       # GET /match/status handler
│   └── models/
│       ├── player.go
│       └── match.go
└── go.mod
```
