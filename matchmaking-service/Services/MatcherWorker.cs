public class MatcherWorker(IQueueManager queueManager) : BackgroundService
{
    private readonly IQueueManager _queueManager = queueManager;

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            await Task.Delay(2);

            var survivors = _queueManager.GetSurvivors();
            if (survivors == null) continue;

            var killer = _queueManager.GetKiller();
            if (killer == null) continue;

            var match = new Match(Guid.NewGuid().ToString(), survivors, killer);
            await _queueManager.CreateMatch(match);
        }
    }
}
