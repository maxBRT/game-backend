public interface IQueueManager
{
    Task AddPlayer(Player player);
    Task<Match?> GetPlayerStatus(string TicketID);
    List<Player>? GetSurvivors();
    Player? GetKiller();
    Task CreateMatch(Match match);
}


