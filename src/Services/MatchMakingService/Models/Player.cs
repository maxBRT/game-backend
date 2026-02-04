using System.Text.Json.Serialization;

public record Player(string Id, string Name, string Role, string TicketID)
{
    [JsonIgnore]
    public bool IsKiller => Role == "killer";
    [JsonIgnore]
    public bool IsSurvivor => Role == "survivor";
    public PlayerInfo ToPlayerInfo() => new(Id, Name, Role);
};

public record PlayerInfo(string Id, string Name, string Role)
{
    public Player ToPlayer() => new(Id, Name, Role, Guid.NewGuid().ToString());
};
