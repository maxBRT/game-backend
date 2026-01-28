using System.Collections.Concurrent;

public class InMemoryQueueService : IQueueService
{
    private readonly ConcurrentDictionary<string, Player> playerMap = new();
    private readonly ConcurrentQueue<Player> _queue = new();

    public async Task<bool> Contains(string TicketID)
    {
        return await Task.FromResult(playerMap.ContainsKey(TicketID));
    }
    public async Task<Player?> Dequeue()
    {
        _queue.TryDequeue(out var player);
        if (player == null) return null;
        playerMap.TryRemove(player.TicketID, out var _);
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
