using System.Diagnostics;
using System.Net;
using System.Net.Http.Json;

public class PlayerClient(HttpClient httpClient, MetricsService metrics) : IPlayerClient
{
    private readonly HttpClient _httpClient = httpClient;
    private readonly MetricsService _metrics = metrics;

    public async Task<GetPlayerResponse> GetPlayer()
    {
        var sw = Stopwatch.StartNew();
        var response = await _httpClient.GetAsync($"/players/{Random.Shared.Next(1, 100)}");
        sw.Stop();

        bool isSystemError = false;
        bool isClientError = false;

        GetPlayerResponse? playerResponse = response.StatusCode switch
        {
            HttpStatusCode.OK => await response.Content.ReadFromJsonAsync<GetPlayerResponse>(),
            HttpStatusCode.NotFound => HandleClientError(),
            _ => HandleSystemError(),
        };

        _metrics.RecordRequest("GetPlayer", sw.Elapsed.TotalMilliseconds, isSystemError, isClientError);

        GetPlayerResponse HandleClientError()
        {
            isClientError = true;
            return new GetPlayerResponse(0, "Not Found", 0, 0, 0);
        }

        GetPlayerResponse HandleSystemError()
        {
            isSystemError = true;
            return new GetPlayerResponse(0, "Something went wrong", 0, 0, 0);
        }

        return playerResponse ?? new GetPlayerResponse(0, "Something went wrong", 0, 0, 0);
    }
}
