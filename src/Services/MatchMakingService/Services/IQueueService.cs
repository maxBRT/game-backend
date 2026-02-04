public interface IQueueService
{
    Task Enqueue(Player player);
    Player? TryDequeue();
    Task<bool> Contains(string TicketID);
    int Count();
}
