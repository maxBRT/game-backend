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

    public async Task<Match?> TryCreateMatch()
    {
        lock (_lock)
        {
            if (survivorQueueService.Count() < 4) return null;
            if (killerQueueService.Count() < 1) return null;

            var killer = killerQueueService.TryDequeue();
            if (killer == null) return null;

            var survivors = new List<Player>();
            for (int i = 0; i < 4; i++)
            {
                var s = survivorQueueService.TryDequeue();
                if (s == null) return null;
                survivors.Add(s);
            }

            var match = new Match(Guid.NewGuid().ToString(), survivors, killer);
            matchStore.AddMatch(match).GetAwaiter().GetResult();
            return match;
        }
    }
}
