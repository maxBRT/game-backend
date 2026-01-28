using System.Collections.Concurrent;

public class InMemoryMatchStore : IMatchStore
{
    private readonly ConcurrentDictionary<string, Match> matchMap = new();
    private readonly ConcurrentDictionary<string, string> playerTicketToMatchIdMap = new();

    public async Task AddMatch(Match match)
    {
        matchMap.TryAdd(match.Id, match);

        foreach (var player in match.Survivors)
        {
            playerTicketToMatchIdMap.TryAdd(player.TicketID, match.Id);
        }

        playerTicketToMatchIdMap.TryAdd(match.Killer.TicketID, match.Id);

    }

    public async Task<Match> GetMatch(string TicketId)
    {
        return await Task.FromResult(matchMap[playerTicketToMatchIdMap[TicketId]]);
    }

    public async Task<Match?> RemoveMatch(string MatchId)
    {
        matchMap.TryRemove(MatchId, out var match);
        if (match == null) return null;
        foreach (var player in match.Survivors)
        {
            playerTicketToMatchIdMap.TryRemove(player.TicketID, out var _);
        }
        playerTicketToMatchIdMap.TryRemove(match.Killer.TicketID, out var _);
        return match;
    }
}
