using player_service.Data;
using CsvHelper;
using System.Globalization;
using player_service.Models;
using CsvHelper.Configuration;

namespace player_service.Services;

public class SeedingService(AppDbContext db)
{
    private readonly AppDbContext _db = db;

    public async Task Seed()
    {
        _db.Players.RemoveRange(_db.Players);

        await _db.SaveChangesAsync();

        var PlayersPath = Path.Combine("Data", "Players.csv");
        var config = new CsvConfiguration(CultureInfo.InvariantCulture)
        {
            PrepareHeaderForMatch = args => args.Header.ToLower()
        };

        using var playerReader = new StreamReader(PlayersPath);
        using var playerCsv = new CsvReader(playerReader, CultureInfo.InvariantCulture);
        {
            var player = playerCsv.GetRecords<Player>().ToList();
            await _db.Players.AddRangeAsync(player);
            await _db.SaveChangesAsync();
        }
    }
}
