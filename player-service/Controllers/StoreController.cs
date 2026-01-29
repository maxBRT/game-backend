using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

[ApiController]
public class StoreController(IStoreService storeService)
{
    [HttpPost("store/buy")]
    public async
        Task<
        Results<Ok<BuyResponse>,
        BadRequest<string>,
        NotFound<string>,
        UnprocessableEntity<string>,
        InternalServerError<string>>>
        BuyItem(BuyRequest request)
    {
        try
        {
            var (result, currency, itemName) = await storeService.BuyItem(request);
            return TypedResults.Ok(new BuyResponse(result, currency, itemName));
        }
        catch (Exception e)
        {
            return e switch
            {
                PlayerNotFoundException => TypedResults.NotFound(e.Message),
                ItemNotFoundException => TypedResults.BadRequest(e.Message),
                NotEnoughCurrencyException => TypedResults.UnprocessableEntity(e.Message),
                _ => TypedResults.InternalServerError(e.Message),
            };
        }

    }
}
