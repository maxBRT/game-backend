using System.Diagnostics;
using System.Net;
using System.Net.Http.Json;

public class MatchClient(HttpClient httpClient, MetricsService metrics) : IMatchClient
{
    private readonly HttpClient _httpClient = httpClient;
    private readonly MetricsService _metrics = metrics;

    public async Task<StatusResponse> GetStatus(string tikeckId)
    {
        var sw = Stopwatch.StartNew();
        var response = await _httpClient.GetAsync($"/match/status/{tikeckId}");
        sw.Stop();

        bool isSystemError = false;

        StatusResponse? statusResponse = response.StatusCode switch
        {
            HttpStatusCode.OK => await response.Content.ReadFromJsonAsync<StatusResponse>(),
            _ => HandleSystemError(),
        };

        _metrics.RecordRequest("MatchStatus", sw.Elapsed.TotalMilliseconds, isSystemError);

        StatusResponse HandleSystemError()
        {
            isSystemError = true;
            return new StatusResponse("Something went wrong", null, null, null, null);
        }

        if (statusResponse == null) return new StatusResponse("Something went wrong", null, null, null, null);
        if (statusResponse.IsMatched) _metrics.IncrementPlayerMatchedCount();

        return statusResponse;
    }

    public async Task<JoinResponse> JoinQueue(JoinRequest request)
    {
        var sw = Stopwatch.StartNew();
        var response = await _httpClient.PostAsJsonAsync("/match/join", request);
        sw.Stop();

        bool isSystemError = false;
        bool isClientError = false;

        JoinResponse? joinResponse = response.StatusCode switch
        {
            HttpStatusCode.OK => await response.Content.ReadFromJsonAsync<JoinResponse>(),
            HttpStatusCode.BadRequest => HandleClientError(),
            _ => HandleSystemError(),
        };

        _metrics.RecordRequest("MatchJoin", sw.Elapsed.TotalMilliseconds, isSystemError, isClientError);

        JoinResponse HandleClientError()
        {
            isClientError = true;
            return new JoinResponse(false, 0, "", "");
        }

        JoinResponse HandleSystemError()
        {
            isSystemError = true;
            return new JoinResponse(false, 0, "", "");
        }

        if (!isSystemError && !isClientError) _metrics.IncrementPlayerInQueueCount();
        return joinResponse ?? new JoinResponse(false, 0, "", "");
    }
}
