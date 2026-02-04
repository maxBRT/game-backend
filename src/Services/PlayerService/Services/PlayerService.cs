using player_service.Data;
using player_service.Exceptions;
using player_service.Models;

namespace player_service.Services;

public class PlayerService(AppDbContext db) : IPlayerService
{
    private readonly AppDbContext _db = db;

    public async Task<Player?> GetPlayer(int id)
    {
        var player = await _db.Players.FindAsync(id);
        return player ?? throw new PlayerNotFoundException("Player not found");
    }
}
