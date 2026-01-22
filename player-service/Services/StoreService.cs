using player_service.Data;
using player_service.Models;
using Microsoft.EntityFrameworkCore;

public class StoreService(AppDbContext db) : IStoreService
{
    private readonly AppDbContext _db = db;

    public async Task<(bool Result, int NewBalance, string ItemName)> BuyItem(BuyRequest request)
    {
        var transaction = await _db.Database.BeginTransactionAsync();
        try
        {

            // Retrieve player and item
            var player = _db.Players.Find(request.PlayerId);
            var item = _db.Items.Find(request.ItemId);

            // Check if player and item exist
            if (player == null)
            {
                throw new PlayerNotFoundException("Player not found");
            }
            if (item == null)
            {
                throw new ItemNotFoundException("Item not found");
            }

            // Check if player has enough currency
            if (player.Currency < item.Price * request.Amount)
            {
                throw new NotEnoughCurrencyException("Not enough currency");
            }

            // Update player currency
            player.Currency -= item.Price * request.Amount;



            // Update inventory
            var inventory = await _db.Inventories
                .FirstOrDefaultAsync(i => i.PlayerId == player.Id && i.ItemId == item.Id);

            if (inventory != null)
            {
                inventory.Amount += request.Amount;
            }
            else
            {
                _db.Inventories.Add(new Inventory
                {
                    Player = player,
                    Item = item,
                    Amount = request.Amount,
                    AquiredAt = DateTime.Now
                });
            }

            // Save changes
            await _db.SaveChangesAsync();
            await transaction.CommitAsync();
            return (true, player.Currency, item.Name);

        }
        catch (Exception e)
        {
            await transaction.RollbackAsync();
            throw e;
        }

    }
}

