# Game Client CLI

## 1. Description

A developer tool built with Go that simulates the full user flow of a game client. This CLI connects to both the Player Service and Matchmaker Service to test the complete backend infrastructure.

**Goal:** Provide a command-line interface for testing and simulating player interactions with the backend services.

## 2. UI/UX

### Command Line Interface

#### `game-cli login --user "<username>"`
Simulates a player login by fetching their data from Service A (Player Service).

**Output:**
```
Logging in as "Nea"...
Player Data Retrieved:
  Level: 42
  XP: 15,230
  Currency: 1,500
Login successful!
```

#### `game-cli find-match`
Connects to Service B (Matchmaker) and waits for the matchmaking loop to complete.

**Output:**
```
Joining matchmaking queue...
Queue position: 3
Waiting for match...
Waiting for match...
Match found!
  Match ID: abc123
  Players: Nea, Dwight, Meg, Jake
Ready to play!
```

#### `game-cli buy --item "<item_id>"`
Purchases an item from the store via Service A.

**Output:**
```
Purchasing item "survivor_outfit_01"...
Purchase successful!
  New Balance: 1,200
  Item: Neon Survivor Outfit
```

## 3. Backend API

This component is a client, not a server. It consumes the following APIs:

### Player Service (Service A)
- `GET /api/player/{id}` - Retrieve player data
- `POST /api/store/buy` - Purchase items

### Matchmaker Service (Service B)
- `POST /match/join` - Join matchmaking queue
- `GET /match/status/{ticketId}` - Poll for match status

## 4. File System

### Technology Stack
- **Language:** Go (Golang)
- **CLI Framework:** cobra or standard flag package
- **HTTP Client:** net/http

### Directory Structure
```
/game-client
├── main.go
├── cmd/
│   ├── root.go         # Root command setup
│   ├── login.go        # login command
│   ├── findmatch.go    # find-match command
│   └── buy.go          # buy command
├── internal/
│   ├── client/
│   │   ├── player.go   # Player service client
│   │   └── match.go    # Matchmaker service client
│   └── config/
│       └── config.go   # Service URLs, timeouts
└── go.mod
```

### Configuration
The CLI should support configuration for service endpoints:
- `PLAYER_SERVICE_URL` (default: `http://localhost:5000`)
- `MATCHMAKER_SERVICE_URL` (default: `http://localhost:8080`)
