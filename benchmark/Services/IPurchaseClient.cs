public interface IPurchaseClient
{
    Task<StoreResponse> PurchaseItem();
    StoreRequest GenerateStoreRequest();

}
