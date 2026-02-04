using Microsoft.EntityFrameworkCore;
using player_service.Data;
using player_service.Exceptions;
using player_service.Models;

namespace PlayerService.Tests;

public class PlayerServiceTests : IDisposable
{
    private readonly AppDbContext _context;
    private readonly player_service.Services.PlayerService _service;

    public PlayerServiceTests()
    {
        var options = new DbContextOptionsBuilder<AppDbContext>()
            .UseInMemoryDatabase(databaseName: Guid.NewGuid().ToString())
            .Options;

        _context = new AppDbContext(options);
        _service = new player_service.Services.PlayerService(_context);
    }

    public void Dispose()
    {
        _context.Dispose();
    }

    [Fact]
    public async Task GetPlayer_WhenPlayerExists_ReturnsPlayer()
    {
        // Arrange
        var player = new Player { Id = 1, Username = "Test Player", Level = 1, Experience = 0, Currency = 0 };
        _context.Players.Add(player);
        await _context.SaveChangesAsync();

        // Act
        var result = await _service.GetPlayer(1);

        // Assert
        Assert.NotNull(result);
        Assert.Equal(1, result.Id);
        Assert.Equal("Test Player", result.Username);
    }

    [Fact]
    public async Task GetPlayer_WhenPlayerDoesNotExist_ThrowsPlayerNotFoundException()
    {
        // Act & Assert
        await Assert.ThrowsAsync<PlayerNotFoundException>(() => _service.GetPlayer(999));
    }
}
