
namespace player_service.Models;

public class Player
{
    public int Id { get; set; }
    public required string Username { get; set; }
    public int Level { get; set; }
    public int Experience { get; set; }
    public int Currency { get; set; }

    public ICollection<Inventory>? Inventories { get; set; }
}
