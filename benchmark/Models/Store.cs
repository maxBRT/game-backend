
public record StoreResponse(bool Success, int NewBalance, string ItemName);

public record StoreRequest(int PlayerId, int ItemId, int Amount);

