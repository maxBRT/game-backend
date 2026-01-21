public interface IStoreService
{
    Task<(bool Result, int NewBalance, string ItemName)> BuyItem(BuyRequest request);
}
