public class MatcherWorker(IQueueManager queueManager) : BackgroundService
{
    private readonly IQueueManager _queueManager = queueManager;

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            await Task.Delay(2);
            await _queueManager.TryCreateMatch();
        }
    }
}
