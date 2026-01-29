public class QueueManager(
    [FromKeyedServices("survivor")] IQueueService survivorQueueService,
    [FromKeyedServices("killer")] IQueueService killerQueueService,
    IMatchStore matchStore) : IQueueManager
{
    private readonly Lock _lock = new();

    public async Task AddPlayer(Player player)
    {
        if (player.IsKiller)
        {
            await killerQueueService.Enqueue(player);
        }
        else if (player.IsSurvivor)
        {
            await survivorQueueService.Enqueue(player);
        }
    }

    public async Task<Match?> GetPlayerStatus(string TicketID)
    {
        var match = await matchStore.GetMatch(TicketID);
        return match;
    }


    public List<Player>? GetSurvivors()
    {
        lock (_lock)
        {
            if (survivorQueueService.Count() < 4) return null;
            var survivors = new List<Player>();

            for (int i = 0; i < 4; i++)
            {
                var s = survivorQueueService.TryDequeue();
                if (s == null) return null;
                survivors.Add(s);
            }

            return survivors;
        }
    }

    public Player? GetKiller()
    {
        lock (_lock)
        {
            return killerQueueService.TryDequeue();
        }
    }

    public async Task CreateMatch(Match match)
    {
        await matchStore.AddMatch(match);
    }
}
