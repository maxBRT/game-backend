# Player Data Service

## 1. Description

The "Persistence" Service built with ASP.NET Core. This service manages the player's profile and economy, handling all persistent data operations with ACID compliance.

**Goal:** Provide a reliable API for player data retrieval and transactional purchases.

## 2. UI/UX

N/A - This is a backend API service with no user interface.

## 3. Backend API

### Endpoints

#### `GET /api/player/{id}`
Retrieve player profile data.

**Response:**
```json
{
  "id": "string",
  "username": "string",
  "level": 0,
  "xp": 0,
  "currency": 0
}
```

#### `POST /api/store/buy`
Transactional purchase of items.

**Request:**
```json
{
  "playerId": "string",
  "itemId": "string",
  "quantity": 1
}
```

**Response:**
```json
{
  "success": true,
  "newBalance": 0,
  "item": {
    "id": "string",
    "name": "string"
  }
}
```

**Transaction Flow:**
1. Check player balance
2. Verify item exists and is purchasable
3. Deduct currency from player
4. Add item to player inventory
5. Commit transaction (rollback on any failure)

## 4. File System

### Technology Stack
- **Framework:** ASP.NET Core
- **ORM:** Entity Framework Core
- **Database:** SQLite (development) / PostgreSQL (production)

### Directory Structure
```
/player-service
├── Controllers/
│   ├── PlayerController.cs
│   └── StoreController.cs
├── Models/
│   ├── Player.cs
│   ├── Item.cs
│   └── Inventory.cs
├── Data/
│   └── AppDbContext.cs
├── Services/
│   ├── PlayerService.cs
│   └── StoreService.cs
└── Program.cs
```
