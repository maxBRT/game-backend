namespace player_service.Models;

public class Inventory
{
    public int PlayerId { get; set; }
    public int ItemId { get; set; }
    public int Amount { get; set; }
    public DateTime AquiredAt { get; set; }

    public Player Player { get; set; }
    public Item Item { get; set; }
}

