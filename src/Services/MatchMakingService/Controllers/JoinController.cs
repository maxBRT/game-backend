using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;

[ApiController]
public class JoinController
{
    [HttpPost("match/join")]
    public async Task<Results<Ok<JoinResponse>, InternalServerError>> Join(JoinRequest request, IQueueManager _queueManager)
    {
        try
        {
            var player = request.Player.ToPlayer();
            await _queueManager.AddPlayer(player);
            return TypedResults.Ok(new JoinResponse(true, player.TicketID, player.Role));
        }
        catch (Exception e)
        {
            Console.Error.WriteLine(e.Message);
            return TypedResults.InternalServerError();
        }
    }
}
