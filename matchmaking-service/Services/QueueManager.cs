public class QueueManager(
    [FromKeyedServices("survivor")] IQueueService survivorQueueService,
    [FromKeyedServices("killer")] IQueueService killerQueueService,
    IMatchStore matchStore) : IQueueManager
{
    private readonly IQueueService _survivorQueueService = survivorQueueService;
    private readonly IQueueService _killerQueueService = killerQueueService;
    private readonly IMatchStore _matchStore = matchStore;

    public async Task AddPlayer(Player player)
    {
        if (player.IsKiller)
        {
            await _killerQueueService.Enqueue(player);
        }
        if (player.IsSurvivor)
        {
            await _survivorQueueService.Enqueue(player);
        }
    }

    public async Task<Match?> GetPlayerStatus(string TicketID)
    {
        var match = await _matchStore.GetMatch(TicketID);
        return match ?? null;
    }

    public async Task<int> SurvivorCount()
    {
        return _survivorQueueService.Count();
    }

    public async Task<int> KillerCount()
    {
        return _killerQueueService.Count();
    }

    public async Task<List<Player>> GetSurvivors()
    {
        var survivors = new List<Player>();
        for (int i = 0; i < 4; i++)
        {
            var s = await _survivorQueueService.Dequeue();
            if (s == null) return [];
            survivors.Add(s);
        }
        return survivors;
    }

    public Task<Player?> GetKiller()
    {
        return _killerQueueService.Dequeue();
    }

    public async Task CreateMatch(Match match)
    {
        await _matchStore.AddMatch(match);
    }
}
