using Scalar.AspNetCore;
using StackExchange.Redis;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring OpenAPI at https://aka.ms/aspnet/openapi
builder.Services.AddOpenApi();


builder.Services.AddSingleton<IQueueManager, QueueManager>();
builder.Services.AddSingleton<IMatchStore, RedisMatchStore>();

var multiplexer = ConnectionMultiplexer.Connect("localhost:6379,allowAdmin=true");

// Hard reset the database
var endPoints = multiplexer.GetEndPoints();
var server = multiplexer.GetServer(endPoints[0]);
server.FlushAllDatabases();

builder.Services.AddSingleton<IConnectionMultiplexer>(multiplexer);
builder.Services.AddKeyedSingleton<IQueueService, RedisQueueService>("survivor");
builder.Services.AddKeyedSingleton<IQueueService, RedisQueueService>("killer");
builder.Services.AddHostedService<MatcherWorker>();
builder.Services.AddControllers();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
    app.MapScalarApiReference();
}


app.MapGet("/health", () => "OK");
app.UseHttpsRedirection();
app.MapControllers();


app.Run();


