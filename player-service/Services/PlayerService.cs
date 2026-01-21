using player_service.Data;
using player_service.Models;

public class PlayerService(AppDbContext db) : IPlayerService
{
    private readonly AppDbContext _db = db;

    public async Task<Player?> GetPlayer(int id)
    {
        return await _db.Players.FindAsync(id);
    }
}
