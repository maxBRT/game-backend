public interface IQueueService
{
    Task Enqueue(Player player);
    Player? TryDequeue();
    bool Contains(string TicketID);
    int Count();
}
