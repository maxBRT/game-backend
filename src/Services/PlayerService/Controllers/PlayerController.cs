using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.AspNetCore.Mvc;
using player_service.Exceptions;
using player_service.Models;
using player_service.Services;

[ApiController]
public class PlayerController(IPlayerService playerService)
{
    private readonly IPlayerService _playerService = playerService;

    [HttpGet("players/{id}")]
    public async Task<Results<Ok<Player>, NotFound, InternalServerError>> GetPlayer(int id)
    {
        try
        {
            var player = await _playerService.GetPlayer(id);
            return TypedResults.Ok(player);
        }
        catch (Exception e)
        {
            return e switch
            {
                PlayerNotFoundException => TypedResults.NotFound(),
                _ => TypedResults.InternalServerError(),
            };
        }
    }
}
