using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

[ApiController]
public class StoreController(IStoreService storeService)
{
    private readonly IStoreService _storeService = storeService;

    [HttpPost("store/buy")]
    public async
        Task<Results<Ok<BuyResponse>, BadRequest, NotFound, UnprocessableEntity, InternalServerError>>
        BuyItem(BuyRequest request)
    {
        try
        {
            var (result, currency, itemName) = await _storeService.BuyItem(request);
            return TypedResults.Ok(new BuyResponse(result, currency, itemName));
        }
        catch (Exception e)
        {
            return e switch
            {
                PlayerNotFoundException => TypedResults.NotFound(),
                ItemNotFoundException => TypedResults.BadRequest(),
                NotEnoughCurrencyException => TypedResults.UnprocessableEntity(),
                _ => TypedResults.InternalServerError(),
            };
        }

    }
}
