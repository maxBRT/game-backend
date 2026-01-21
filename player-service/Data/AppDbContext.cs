using Microsoft.EntityFrameworkCore;

namespace player_service.Data;

public class AppDbContext : DbContext
{

    public DbSet<Player> Players { get; set; }
    public DbSet<Item> Items { get; set; }
    public DbSet<Inventory> Inventories { get; set; }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        optionsBuilder.UseSqlite("Data Source=player-service.db");
    }
}


