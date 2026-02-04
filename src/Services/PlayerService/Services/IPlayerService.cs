using player_service.Models;

namespace player_service.Services;

public interface IPlayerService
{
    Task<Player?> GetPlayer(int id);
}
