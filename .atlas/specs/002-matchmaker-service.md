# Matchmaker Service

## 1. Description

The high-performance matchmaking service built with Go. This service groups players into lobbies using goroutines and in-memory queues for maximum concurrency and low latency.

**Goal:** Efficiently match players into game lobbies when enough players are queued in both role-specific queues.

**Queues:**
- **Survivor Queue:** Players queuing as survivors
- **Killer Queue:** Players queuing as killers

**Matching Rule:** When `len(survivorQueue) >= 4` AND `len(killerQueue) >= 1`, create a match with 4 survivors and 1 killer, removing those players from their respective queues.

## 2. UI/UX

N/A - This is a backend API service with no user interface.

## 3. Backend API

### Endpoints

#### `POST /match/join`
Add a player to the matchmaking queue for their chosen role.

**Request:**
```json
{
  "playerId": "string",
  "playerName": "string",
  "role": "survivor" | "killer"
}
```

**Response:**
```json
{
  "success": true,
  "queuePosition": 0,
  "ticketId": "string",
  "role": "survivor" | "killer"
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
  "survivors": [
    {"id": "string", "name": "string"},
    {"id": "string", "name": "string"},
    {"id": "string", "name": "string"},
    {"id": "string", "name": "string"}
  ],
  "killer": {"id": "string", "name": "string"}
}
```

### Matchmaking Logic

1. Player joins appropriate queue via `/match/join` based on chosen role
2. Background goroutine continuously monitors both queues
3. When survivor queue has 4+ players AND killer queue has 1+ player:
   - Create new match with unique ID
   - Assign first 4 survivors from survivor queue
   - Assign first 1 killer from killer queue
   - Remove matched players from their respective queues
   - Update match status for all 5 players
4. Players polling `/match/status/{ticketId}` receive match info with role assignments

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
│   │   ├── queue.go        # Queue data structure (survivor + killer queues)
│   │   └── matcher.go      # Matching logic goroutine (monitors both queues)
│   ├── handlers/
│   │   ├── join.go         # POST /match/join handler
│   │   └── status.go       # GET /match/status handler
│   └── models/
│       ├── player.go       # Player model with role field
│       └── match.go        # Match model with survivors/killer fields
└── go.mod
```
