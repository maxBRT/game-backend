using player_service.Models;

public interface IPlayerService
{
    Task<Player?> GetPlayer(int id);
}
