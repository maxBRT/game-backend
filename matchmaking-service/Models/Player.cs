using System.Text.Json.Serialization;

public record Player(string Id, string Name, string Role, int Position, string TicketID)
{
    [JsonIgnore]
    public bool IsKiller => Role == "killer";
    [JsonIgnore]
    public bool IsSurvivor => Role == "survivor";

};

public record PlayerInfo(string Id, string Name, string Role);
