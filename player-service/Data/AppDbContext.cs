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
        var connectionString = Environment.GetEnvironmentVariable("CONNECTION_STRING");
        optionsBuilder.UseNpgsql(connectionString);
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

    }
}


