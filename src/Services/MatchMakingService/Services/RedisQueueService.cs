using System.Text.Json;
using StackExchange.Redis;

public class RedisQueueService(IConnectionMultiplexer multiplexer, [ServiceKey] string queueKey) : IQueueService
{

    private readonly IDatabase _db = multiplexer.GetDatabase();
    private readonly string _queueKey = queueKey;

    public async Task<bool> Contains(string TicketID)
    {
        return await _db.KeyExistsAsync(TicketID);
    }

    public int Count()
    {
        return Convert.ToInt32(_db.ListLength(_queueKey));
    }

    public async Task Enqueue(Player player)
    {
        var key = player.TicketID;
        var value = JsonSerializer.Serialize(player);

        await _db.SetAddAsync(key, value);
        await _db.ListRightPushAsync(_queueKey, value);
    }

    public Player? TryDequeue()
    {
        RedisValue value = _db.ListLeftPop(_queueKey);
        if (value.IsNull)
        {
            return null;
        }
        var player = JsonSerializer.Deserialize<Player>((byte[])value);

        if (player == null) return null;

        _db.KeyDelete(player.TicketID);

        return player;
    }
}
