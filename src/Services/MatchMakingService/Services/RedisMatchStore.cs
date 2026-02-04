using System.Text.Json;
using StackExchange.Redis;

public class RedisMatchStore(IConnectionMultiplexer connectionMultiplexer) : IMatchStore
{
    private readonly IDatabase _db = connectionMultiplexer.GetDatabase();
    private readonly string matchMap = "matchMap";
    private readonly string playerTicketToMatchIdMap = "playerTicketToMatchIdMap";

    public async Task AddMatch(Match match)
    {
        await _db.HashSetAsync(matchMap, match.Id, JsonSerializer.Serialize(match));

        foreach (var player in match.Survivors)
        {
            await _db.HashSetAsync(playerTicketToMatchIdMap, player.TicketID, match.Id);
        }

        await _db.HashSetAsync(playerTicketToMatchIdMap, match.Killer.TicketID, match.Id);
    }

    public async Task<Match?> GetMatch(string TicketId)
    {
        var matchId = await _db.HashGetAsync(playerTicketToMatchIdMap, TicketId);
        if (matchId.IsNull) return null;

        RedisValue value = await _db.HashGetAsync(matchMap, matchId);
        if (value.IsNull) return null;

        return JsonSerializer.Deserialize<Match>((byte[])value);
    }

}
