public class MatcherWorker(IQueueManager queueManager) : BackgroundService
{
    private readonly IQueueManager _queueManager = queueManager;

    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            await Task.Delay(2);
            if (await _queueManager.KillerCount() == 0) continue;
            if (await _queueManager.SurvivorCount() < 4) continue;

            var killer = await _queueManager.GetKiller();
            if (killer == null) continue;

            var survivors = await _queueManager.GetSurvivors();
            if (survivors.Count < 4) continue;

            var match = new Match(Guid.NewGuid().ToString(), survivors, killer);
            await _queueManager.CreateMatch(match);
        }

    }
}
