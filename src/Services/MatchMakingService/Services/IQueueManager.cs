public interface IQueueManager
{
    Task AddPlayer(Player player);
    Task<Match?> GetPlayerStatus(string TicketID);
    Task<Match?> TryCreateMatch();
}


