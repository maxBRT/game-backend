using System.Diagnostics;
using System.Net;
using System.Net.Http.Json;

public class PurchaseClient(HttpClient httpClient, MetricsService metrics) : IPurchaseClient
{
    private readonly HttpClient _httpClient = httpClient;
    private readonly MetricsService _metrics = metrics;

    public async Task<StoreResponse> PurchaseItem()
    {
        var sw = Stopwatch.StartNew();
        var request = GenerateStoreRequest();
        var response = await _httpClient.PostAsJsonAsync("/store/buy", request);
        sw.Stop();

        bool isSystemError = false;
        bool isClientError = false;

        StoreResponse? storeResponse = response.StatusCode switch
        {
            HttpStatusCode.OK => await response.Content.ReadFromJsonAsync<StoreResponse>(),
            HttpStatusCode.NotFound => HandleClientError("Not Found"),
            HttpStatusCode.BadRequest => HandleClientError("Bad Request"),
            HttpStatusCode.UnprocessableEntity => HandleClientError("Unprocessable Entity"),
            _ => HandleSystemError(),
        };

        _metrics.RecordRequest("Purchase", sw.Elapsed.TotalMilliseconds, isSystemError, isClientError);

        StoreResponse HandleClientError(string message)
        {
            isClientError = true;
            return new StoreResponse(false, 0, message);
        }

        StoreResponse HandleSystemError()
        {
            isSystemError = true;
            return new StoreResponse(false, 0, "Something went wrong");
        }

        return storeResponse ?? new StoreResponse(false, 0, "Something went wrong");
    }

    public StoreRequest GenerateStoreRequest()
    {
        return new StoreRequest(
                Random.Shared.Next(1, 100),
                Random.Shared.Next(1, 100),
                Random.Shared.Next(1, 3)
                );
    }
}
