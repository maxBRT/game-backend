using Microsoft.EntityFrameworkCore;
using player_service.Models;

namespace player_service.Data;

public class AppDbContext : DbContext
{

    public DbSet<Player> Players { get; set; }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        var connectionString = Environment.GetEnvironmentVariable("CONNECTION_STRING");
        optionsBuilder.UseNpgsql(connectionString);
    }


}


