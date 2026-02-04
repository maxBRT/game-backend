# Game Backend

A microservices backend for an asymmetrical multiplayer game (ex: Dead By Daylight). Built with .NET 10 and ASP.NET Core.

The goal of this project is to demonstrate my understanding of the different layers of a live service game backend and how they interact with each other.

While this project is simple in terms of features, I put a lot of focus on making the architecture as modular as possible to showcase my ability to design and architect a service-oriented architecture.

- The code has an automated test suite.
- The queue is built in a way that redis could be swap for anything with minimal overhead
- The application is fully containerized and all the services can be launched with a single command.

## Services

| Service | Port | Description |
|---------|------|-------------|
| **player-service** | 5043 | Player profiles |
| **player-database** | 5432 | Postgres database for player profiles |
| **matchmaking-service** | 8000 | Queues players and creates matches |
| **redis** | 6379 | Queues and stores match data |
| **benchmark** | - | Load testing CLI tool |

## Quick Start

### Using Docker Compose

```bash
docker-compose up -d --build
```

This starts both services with health checks and persistent storage.

## API Endpoints

Documentation for the API endpoints can be viewed at `/scalar/v1` for each service.

### Player Service (localhost:5043)

```
GET  /players/{id}     - Get player profile
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
-s <count>  Number of survivors (default: 800)
-k <count>  Number of killers (default: 200)
```

**Example:**
```bash
dotnet run -s 10000 -k 2500 
```

## How Matchmaking Works

1. Players join queue via `/match/join` and receive a `ticketId`
2. Background worker matches 4 survivors + 1 killer
3. Players poll `/match/status/{ticketId}` until matched
4. Response includes match ID and all player info
