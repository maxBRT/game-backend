using System.Text.Json.Serialization;

public record JoinRequest(PlayerInfo Player);

public record JoinResponse(bool Success, int QueuePosition, string TicketID, string Role);

public record StatusResponse(
        string Status,
        // If Waiting
        int? QueuePosition,
        // If Matched
        string? MatchId,
        List<PlayerInfo>? Survivors,
        PlayerInfo? Killer
    )
{
    [JsonIgnore]
    public bool IsWaiting => Status == "waiting";
    [JsonIgnore]
    public bool IsMatched => Status == "matched";
}
