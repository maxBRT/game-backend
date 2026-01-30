# Code Review - Enhancement Observations

This document outlines improvements to make the codebase interview-ready.

---

## Code Quality Issues

### 8. Redundant field assignments with primary constructors
**Files:** Multiple controllers and services
```csharp
public class StoreController(IStoreService storeService)
{
    private readonly IStoreService _storeService = storeService;  // ❌ Redundant
```
With primary constructors, `storeService` is already a field. Remove the explicit field.

### 10. Role as string instead of enum
**File:** `matchmaking-service/Models/Player.cs`
```csharp
public record Player(string Id, string Name, string Role, string TicketID)
```
Using strings allows invalid values. A `Role` enum already exists but isn't used.

---

## Naming & Consistency

### 11. Inconsistent casing
- `TicketID` vs `TicketId` - Pick one convention (C# standard: `TicketId`)
- `MatchId` vs `matchId` - Parameters should use camelCase

### 12. Project naming uses kebab-case
- `player-service` should be `PlayerService` for C# conventions
- Namespace `player_service` uses snake_case (inconsistent)

### 13. Magic numbers
**File:** `matchmaking-service/Services/MatcherWorker.cs`
```csharp
await Task.Delay(2);           // What does 2 mean?
if (await _queueManager.SurvivorCount() < 4)  // Why 4?
```
Extract to named constants:
```csharp
private const int MatchPollingIntervalMs = 2;
private const int SurvivorsPerMatch = 4;
```

---

## Missing Patterns

### 14. No input validation
**Files:** All controllers

Requests are accepted without validation. Add:
```csharp
public record BuyRequest(
    [Required] int PlayerId,
    [Required] int ItemId,
    [Range(1, 100)] int Amount
);
```

### 15. No logging infrastructure
Using `Console.Error.WriteLine` instead of `ILogger<T>`. Add structured logging.

### 16. Missing health endpoint in player-service
**File:** `player-service/Program.cs`

Matchmaking has `/health`, player-service doesn't. Add for consistency.

### 17. No cancellation token propagation
**File:** `matchmaking-service/Services/MatcherWorker.cs`

The `stoppingToken` isn't passed to async operations inside the loop.

---

## Minor Cleanup

### 18. Extra blank lines
**File:** `player-service/Services/StoreService.cs:14, 37-38`

Remove extra blank lines for consistency.

### 19. Transaction could use using declaration
**File:** `player-service/Services/StoreService.cs:11`
```csharp
var transaction = await _db.Database.BeginTransactionAsync();
```
Could be:
```csharp
await using var transaction = await _db.Database.BeginTransactionAsync();
```
This ensures disposal even if early return is added later.

### 20. Sync Find() in async method
**File:** `player-service/Services/StoreService.cs:16-17`
```csharp
var player = _db.Players.Find(request.PlayerId);      // Sync
var item = _db.Items.Find(request.ItemId);            // Sync
```
Use `FindAsync()` for consistency in an async method.

---

## Suggested Priority

1. **High:** Items 1-5 (actual bugs)
2. **Medium:** Items 6-10 (code quality that interviewers will notice)
3. **Low:** Items 11-20 (polish)

---

## Interview Talking Points

After fixing these, you can discuss:
- **Why** you chose primary constructors and when the redundant fields were an oversight
- The **trade-offs** of in-memory storage vs persistent backends
- How you'd handle the **race condition** in a production system
- Your approach to **structured logging** and **observability**
