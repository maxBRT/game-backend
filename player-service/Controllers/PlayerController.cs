using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;
using player_service.Models;

[ApiController]
public class PlayerController(IPlayerService playerService)
{
    private readonly IPlayerService _playerService = playerService;

    [HttpGet("players/{id}")]
    public async Task<Results<Ok<Player>, NotFound>> GetPlayer(int id)
    {
        var player = await _playerService.GetPlayer(id);

        if (player == null)
        {
            return TypedResults.NotFound();
        }

        return TypedResults.Ok(player);
    }
}
