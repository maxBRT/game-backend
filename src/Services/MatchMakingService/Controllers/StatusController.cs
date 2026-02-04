using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

[ApiController]
public class StatusController
{
    [HttpGet("match/status/{TicketID}")]
    public async Task<Results<Ok<StatusResponse>, InternalServerError>> Status(string TicketID, CancellationToken cancellationToken, IQueueManager _queueManager)
    {
        var cts = CancellationTokenSource.CreateLinkedTokenSource(cancellationToken);
        cts.CancelAfter(TimeSpan.FromSeconds(30));
        try
        {
            while (!cts.IsCancellationRequested)
            {
                await Task.Delay(10);
                var match = await _queueManager.GetPlayerStatus(TicketID);

                if (match == null) continue;

                var survivors = match.Survivors.Select(p => p.ToPlayerInfo()).ToList();

                return TypedResults.Ok(new StatusResponse("matched", match.Id, survivors, match.Killer.ToPlayerInfo()));
            }

            return TypedResults.Ok(new StatusResponse("waiting", null, null, null));
        }
        catch (TaskCanceledException)
        {
            return TypedResults.Ok(new StatusResponse("waiting", null, null, null));
        }
        catch (Exception e)
        {
            Console.Error.WriteLine(e.Message, "Error getting status");
            return TypedResults.InternalServerError();
        }
    }
}
