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

    public async Task<Match?> GetMatch(string TicketId)
    {
        playerTicketToMatchIdMap.TryGetValue(TicketId, out var matchId);
        if (matchId == null) return null;
        return matchMap.TryGetValue(matchId, out var match) ? match : null;
    }

}
