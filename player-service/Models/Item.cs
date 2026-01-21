namespace player_service.Models;

public class Item
{
    public int Id { get; set; }
    public required string Name { get; set; }
    public int Price { get; set; }


    public ICollection<Inventory>? Inventories { get; set; }
}
