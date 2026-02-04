namespace MatchMakingService.Tests;

public class InMemoryQueueServiceTests
{
    private readonly InMemoryQueueService _queueService;

    public InMemoryQueueServiceTests()
    {
        _queueService = new InMemoryQueueService();
    }

    [Fact]
    public async Task Enqueue_AddsPlayerToQueue()
    {
        // Arrange
        var player = new Player("1", "TestPlayer", "survivor", "ticket-1");

        // Act
        await _queueService.Enqueue(player);

        // Assert
        Assert.True(await _queueService.Contains("ticket-1"));
    }


    [Fact]
    public async Task Contains_WhenPlayerNotInQueue_ReturnsFalse()
    {
        Assert.False(await _queueService.Contains("wrong"));
    }

    [Fact]
    public async Task TryDequeue_WhenQueueHasPlayers_ReturnsPlayer()
    {
        // Arrange
        var player = new Player("1", "TestPlayer", "survivor", "ticket-1");
        await _queueService.Enqueue(player);

        // Act
        var p = _queueService.TryDequeue();

        // Assert
        Assert.NotNull(p);
        Assert.Equal(p.Id, player.Id);
    }

    [Fact]
    public void TryDequeue_WhenQueueEmpty_ReturnsNull()
    {
        // Act
        var x = _queueService.TryDequeue();

        // Assert
        Assert.Null(x);
    }

    [Fact]
    public async Task Count_ReturnsCorrectCount()
    {
        // Arrange
        var player1 = new Player("1", "TestPlayer", "survivor", "ticket-1");
        var player2 = new Player("2", "TestPlayer2", "survivor", "ticket-2");
        var player3 = new Player("3", "TestPlayer3", "survivor", "ticket-3");

        // Act
        await _queueService.Enqueue(player1);
        await _queueService.Enqueue(player2);
        await _queueService.Enqueue(player3);

        // Assert
        Assert.Equal(3, _queueService.Count());
    }
}

public class InMemoryMatchStoreTests
{
    private readonly InMemoryMatchStore _matchStore;

    public InMemoryMatchStoreTests()
    {
        _matchStore = new InMemoryMatchStore();
    }

    [Fact]
    public async Task AddMatch_StoresMatch()
    {
        // Arrange
        var killer = new Player("1", "Killer", "killer", "ticket-killer");
        var survivors = Enumerable.Range(1, 4)
            .Select(i => new Player($"s{i}", $"Survivor{i}", "survivor", $"ticket-s{i}"))
            .ToList();
        var match = new Match("match-1", survivors, killer);

        // Act
        await _matchStore.AddMatch(match);

        // Assert
        Assert.NotNull(_matchStore.GetMatch("ticket-killer"));

    }

    [Fact]
    public async Task GetMatch_WhenNoMatchExists_ReturnsNull()
    {
        // Act
        var m = await _matchStore.GetMatch("nothing");
        // Assert
        Assert.Null(m);
    }
}

public class QueueManagerTests
{
    private readonly InMemoryQueueService _survivorQueue;
    private readonly InMemoryQueueService _killerQueue;
    private readonly InMemoryMatchStore _matchStore;
    private readonly QueueManager _queueManager;

    public QueueManagerTests()
    {
        _survivorQueue = new InMemoryQueueService();
        _killerQueue = new InMemoryQueueService();
        _matchStore = new InMemoryMatchStore();
        _queueManager = new QueueManager(_survivorQueue, _killerQueue, _matchStore);
    }

    [Fact]
    public async Task AddPlayer_WhenSurvivor_AddsToSurvivorQueue()
    {
        // Arrange
        var survivor = new Player("1", "Survivor", "survivor", "ticket-1");

        // Act
        await _queueManager.AddPlayer(survivor);

        // Assert
        Assert.True(await _survivorQueue.Contains("ticket-1"));
    }

    [Fact]
    public async Task AddPlayer_WhenKiller_AddsToKillerQueue()
    {
        // Arrange
        var killer = new Player("1", "Killer", "killer", "ticket-1");

        // Act
        await _queueManager.AddPlayer(killer);

        // Assert
        Assert.True(await _killerQueue.Contains("ticket-1"));

    }

    [Fact]
    public async Task TryCreateMatch_WhenNotEnoughSurvivors_ReturnsNull()
    {
        // Arrange
        var killer = new Player("1", "Killer", "killer", "ticket-killer");
        await _queueManager.AddPlayer(killer);

        // Act
        var m = await _queueManager.TryCreateMatch();

        // Assert
        Assert.Null(m);
    }

    [Fact]
    public async Task TryCreateMatch_WhenEnoughPlayers_ReturnsMatch()
    {
        // Arrange
        var killer = new Player("1", "Killer", "killer", "ticket-killer");
        await _queueManager.AddPlayer(killer);

        for (int i = 0; i < 4; i++)
        {
            await _queueManager.AddPlayer(new Player($"s{i}", $"Survivor{i}", "survivor", $"ticket-s{i}"));
        }

        // Act
        var m = await _queueManager.TryCreateMatch();

        // Assert
        Assert.NotNull(m);
    }

    [Fact]
    public async Task TryCreateMatch_WhenMatchCreated_RemovesPlayersFromQueues()
    {
        // Arrange
        var killer = new Player("1", "Killer", "killer", "ticket-killer");
        await _queueManager.AddPlayer(killer);

        for (int i = 0; i < 4; i++)
        {
            await _queueManager.AddPlayer(new Player($"s{i}", $"Survivor{i}", "survivor", $"ticket-s{i}"));
        }

        // Act
        await _queueManager.TryCreateMatch();

        // Assert
        Assert.Equal(0, _survivorQueue.Count());
        Assert.Equal(0, _killerQueue.Count());
    }

    [Fact]
    public async Task GetPlayerStatus_WhenInMatch_ReturnsMatch()
    {
        // Arrange
        var killer = new Player("1", "Killer", "killer", "ticket-killer");
        await _queueManager.AddPlayer(killer);

        for (int i = 0; i < 4; i++)
        {
            await _queueManager.AddPlayer(new Player($"s{i}", $"Survivor{i}", "survivor", $"ticket-s{i}"));
        }
        await _queueManager.TryCreateMatch();

        // Act
        var m = await _queueManager.GetPlayerStatus("ticket-killer");

        // Assert
        Assert.NotNull(m);
    }

    [Fact]
    public async Task GetPlayerStatus_WhenNotInMatch_ReturnsNull()
    {
        // Arrange
        var player = new Player("1", "Survivor", "survivor", "ticket-1");
        await _queueManager.AddPlayer(player);

        // Act
        var m = await _queueManager.GetPlayerStatus("ticket-1");

        // Assert
        Assert.Null(m);
    }
}
