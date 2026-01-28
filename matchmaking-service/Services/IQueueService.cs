public interface IQueueService
{
    Task Enqueue(Player player);
    Task<Player?> Dequeue();
    Task<bool> Contains(string TicketID);
    int Count();
}
