# Game Backend

A microservices backend for an asymmetrical multiplayer game (ex: Dead By Daylight). Built with .NET 10 and ASP.NET Core.

The goal of this project is to demonstrate my understanding of the different layers of a live service game backend and how they interact with each other.

## Services

| Service | Port | Description |
|---------|------|-------------|
| **player-service** | 5043 | Player profiles, inventory, and store |
| **matchmaking-service** | 8000 | Queues players and creates matches |
| **redis** | 6379 | Queues and stores match data |
| **benchmark** | - | Load testing CLI tool |

## Quick Start

### Using Docker Compose (Recommended)

```bash
docker-compose up -d
```

This starts both services with health checks and persistent storage.

## API Endpoints

Documentation for the API endpoints can be viewed at `/scalar/v1` for each service.

### Player Service (localhost:5043)

```
GET  /players/{id}     - Get player profile
POST /store/buy        - Purchase an item
```

### Matchmaking Service (localhost:8000)

```
POST /match/join              - Join matchmaking queue
GET  /match/status/{ticketId} - Check match status (long-polls up to 30s)
GET  /health                  - Health check
```

## Running Benchmarks

```bash
cd benchmark
dotnet run
```

**Options:**
```
-s <count>  Number of survivors (default: 100)
-k <count>  Number of killers (default: 20)
-b <count>  Number of purchases (default: 50)
```

**Example:**
```bash
dotnet run -s 10000 -k 2500 -b 100
```

## How Matchmaking Works

1. Players join queue via `/match/join` and receive a `ticketId`
2. Background worker matches 4 survivors + 1 killer
3. Players poll `/match/status/{ticketId}` until matched
4. Response includes match ID and all player info
