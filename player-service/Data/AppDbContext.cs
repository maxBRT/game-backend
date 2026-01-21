using Microsoft.EntityFrameworkCore;
using player_service.Models;

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

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);

        // Composite key
        modelBuilder.Entity<Inventory>()
            .HasKey(i => new { i.PlayerId, i.ItemId });

        // Navigation properties
        modelBuilder.Entity<Inventory>()
            .HasOne(i => i.Player)
            .WithMany(p => p.Inventories)
            .HasForeignKey(i => i.PlayerId);

        modelBuilder.Entity<Inventory>()
            .HasOne(i => i.Item)
            .WithMany(i => i.Inventories)
            .HasForeignKey(i => i.ItemId);

        // Seed data
        modelBuilder.Entity<Player>().HasData(
            new Player { Id = 1, Username = "Max", Level = 1, Experience = 0, Currency = 100 },
            new Player { Id = 2, Username = "John", Level = 1, Experience = 0, Currency = 100 },
            new Player { Id = 3, Username = "Jane", Level = 1, Experience = 0, Currency = 100 },
            new Player { Id = 4, Username = "Bob", Level = 1, Experience = 0, Currency = 100 },
            new Player { Id = 5, Username = "Alice", Level = 1, Experience = 0, Currency = 100 }
        );

        modelBuilder.Entity<Item>().HasData(
            new Item { Id = 1, Name = "Sword", Price = 20 },
            new Item { Id = 2, Name = "Boots", Price = 5 },
            new Item { Id = 3, Name = "Shield", Price = 10 }
        );
    }
}


