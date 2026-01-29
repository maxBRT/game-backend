public interface IQueueManager
{
    Task AddPlayer(Player player);
    Task<Match?> GetPlayerStatus(string TicketID);
    Task<int> SurvivorCount();
    Task<int> KillerCount();
    Task<List<Player>> GetSurvivors();
    Task<Player?> GetKiller();
    Task CreateMatch(Match match);
}


