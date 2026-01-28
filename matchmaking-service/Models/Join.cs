public record JoinRequest(PlayerInfo Player);
public record JoinResponse(bool Success, string TicketID, string Role);
