using System.Collections.Concurrent;

public class InMemoryQueueService : IQueueService
{
    private readonly ConcurrentDictionary<string, Player> playerMap = new();
    private readonly ConcurrentQueue<Player> _queue = new();

    public Task<bool> Contains(string TicketID)
    {
        return Task.FromResult(playerMap.ContainsKey(TicketID));
    }
    public Player? TryDequeue()
    {
        _queue.TryDequeue(out var player);
        return player;
    }

    public async Task Enqueue(Player player)
    {
        playerMap.TryAdd(player.TicketID, player);
        _queue.Enqueue(player);
    }

    public int Count()
    {
        return _queue.Count;
    }
}
