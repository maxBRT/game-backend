public interface IMatchStore
{
    Task AddMatch(Match match);
    Task<Match> GetMatch(string TicketId);
    Task<Match?> RemoveMatch(string MatchId);
}
