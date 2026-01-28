using System.Diagnostics;
using Spectre.Console;

public class App(IPurchaseClient purchaseClient, MetricsService metricsService, IPlayerClient playerClient, IMatchClient matchClient)
{
    private readonly IPurchaseClient purchaseClient = purchaseClient;
    private readonly MetricsService metricsService = metricsService;
    private readonly IPlayerClient playerClient = playerClient;
    private readonly IMatchClient matchClient = matchClient;

    public async Task Run(BenchmarkConfig config)
    {
        var stopwatch = Stopwatch.StartNew();

        AnsiConsole.MarkupLine("[bold blue]Starting Benchmark[/]");
        AnsiConsole.WriteLine();

        await AnsiConsole.Progress()
            .AutoRefresh(true)
            .AutoClear(false)
            .HideCompleted(false)
            .Columns(
                new TaskDescriptionColumn(),
                new ProgressBarColumn(),
                new PercentageColumn(),
                new SpinnerColumn())
            .StartAsync(async ctx =>
            {
                var survivorTask = ctx.AddTask("[green]Queueing Survivors[/]", maxValue: config.SurvivorCount);
                var killerTask = ctx.AddTask("[red]Queueing Killers[/]", maxValue: config.KillerCount);
                var purchaseTask = ctx.AddTask("[yellow]Processing Purchases[/]", maxValue: config.PurchaseCount);

                await QueueSurvivors(config.SurvivorCount, survivorTask);
                await QueueKillers(config.KillerCount, killerTask);
                await PurchaseItems(config.PurchaseCount, purchaseTask);
            });

        stopwatch.Stop();

        PrintSummary(stopwatch.Elapsed, config);
    }

    private async Task QueueSurvivors(int n, ProgressTask task)
    {
        for (int i = 0; i < n; i++)
        {
            var p = await playerClient.GetPlayer();
            await matchClient.JoinQueue(new JoinRequest(p.ToPlayerInfoSurvivor()));
            task.Increment(1);
        }
    }

    private async Task QueueKillers(int n, ProgressTask task)
    {
        for (int i = 0; i < n; i++)
        {
            var p = await playerClient.GetPlayer();
            await matchClient.JoinQueue(new JoinRequest(p.ToPlayerInfoKiller()));
            task.Increment(1);
        }
    }

    private async Task PurchaseItems(int n, ProgressTask task)
    {
        for (int i = 0; i < n; i++)
        {
            await purchaseClient.PurchaseItem();
            task.Increment(1);
        }
    }

    private void PrintSummary(TimeSpan elapsed, BenchmarkConfig config)
    {
        AnsiConsole.WriteLine();
        AnsiConsole.Write(new Rule("[bold blue]Benchmark Summary[/]").RuleStyle("blue"));
        AnsiConsole.WriteLine();

        var totalRequests = metricsService.GetTotalRequestCount();
        var throughput = totalRequests / elapsed.TotalSeconds;

        var overviewTable = new Table()
            .Border(TableBorder.Rounded)
            .AddColumn("Metric")
            .AddColumn("Value");

        overviewTable.AddRow("Duration", $"{elapsed.TotalSeconds:F2}s");
        overviewTable.AddRow("Total Requests", $"{totalRequests}");
        overviewTable.AddRow("Throughput", $"{throughput:F2} req/s");
        overviewTable.AddRow("[green]Success[/]", $"{totalRequests - metricsService.GetSystemErrorCount() - metricsService.GetClientErrorCount()}");
        overviewTable.AddRow("[yellow]Client Errors[/]", $"{metricsService.GetClientErrorCount()}");
        overviewTable.AddRow("[red]System Errors[/]", $"{metricsService.GetSystemErrorCount()}");

        AnsiConsole.Write(overviewTable);
        AnsiConsole.WriteLine();

        AnsiConsole.Write(new Rule("[bold]Latency by Operation (ms)[/]").RuleStyle("grey"));
        AnsiConsole.WriteLine();

        var latencyTable = new Table()
            .Border(TableBorder.Rounded)
            .AddColumn("Operation")
            .AddColumn("Count")
            .AddColumn("Min")
            .AddColumn("Avg")
            .AddColumn("Max");

        foreach (var (operation, metrics) in metricsService.GetAllOperationMetrics().OrderBy(x => x.Key))
        {
            var stats = metrics.GetLatencyStats();
            latencyTable.AddRow(
                operation,
                metrics.Count.ToString(),
                $"{stats.Min:F2}",
                $"{stats.Avg:F2}",
                $"{stats.Max:F2}"
            );
        }

        AnsiConsole.Write(latencyTable);
        AnsiConsole.WriteLine();
    }
}
