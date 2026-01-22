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
        _db.Items.RemoveRange(_db.Items);
        _db.Players.RemoveRange(_db.Players);
        _db.Inventories.RemoveRange(_db.Inventories);

        await _db.SaveChangesAsync();

        var PlayersPath = Path.Combine("Data", "Players.csv");
        var ItemsPath = Path.Combine("Data", "Items.csv");
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
        using var itemsReader = new StreamReader(ItemsPath);
        using var itemsCsv = new CsvReader(itemsReader, CultureInfo.InvariantCulture);
        {
            var item = itemsCsv.GetRecords<Item>().ToList();
            await _db.Items.AddRangeAsync(item);
            await _db.SaveChangesAsync();
        }

    }
}
