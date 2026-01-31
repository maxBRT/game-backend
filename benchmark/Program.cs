using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.DependencyInjection;
using System.CommandLine;

Option<int> survivorOption = new("-s")
{
    Description = "The number of survivor to queue"
};
Option<int> killerOption = new("-k")
{
    Description = "The number of killer to queue"
};
Option<int> purchaseOption = new("-b")
{
    Description = "The number of purchase request"
};



var rootCommand = new RootCommand("Game Backend Benchmark");
rootCommand.Options.Add(survivorOption);
rootCommand.Options.Add(killerOption);
rootCommand.Options.Add(purchaseOption);
ParseResult parsedResult = rootCommand.Parse(args);

var config = new BenchmarkConfig(
        parsedResult.GetValue(survivorOption) > 0 ? parsedResult.GetValue(survivorOption) : 800,
        parsedResult.GetValue(killerOption) > 0 ? parsedResult.GetValue(killerOption) : 200,
        parsedResult.GetValue(purchaseOption) > 0 ? parsedResult.GetValue(purchaseOption) : 500
        );

using IHost host = Host.CreateDefaultBuilder(args)
    .ConfigureServices((context, services) =>
    {
        services.AddSingleton<MetricsService>();
        services.AddHttpClient<IPlayerClient, PlayerClient>(client =>
        {
            client.BaseAddress = new Uri("http://localhost:5043");
            client.DefaultRequestHeaders.Add("Accept", "application/json");
        });
        services.AddHttpClient<IPurchaseClient, PurchaseClient>(client =>
        {
            client.BaseAddress = new Uri("http://localhost:5043");
            client.DefaultRequestHeaders.Add("Accept", "application/json");
        });
        services.AddHttpClient<IMatchClient, MatchClient>(client =>
        {
            client.BaseAddress = new Uri("http://localhost:8000");
            client.DefaultRequestHeaders.Add("Accept", "application/json");
        });
        services.AddTransient<App>();
    })
    .Build();

await host.Services.GetRequiredService<App>().Run(config);

