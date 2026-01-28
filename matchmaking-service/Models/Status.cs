using System.Text.Json.Serialization;

public record StatusResponse(string Status, string? MatchID, List<PlayerInfo>? Survivors, PlayerInfo? Killer)
{
    [JsonIgnore]
    public bool IsMatched => MatchID != null;
};
