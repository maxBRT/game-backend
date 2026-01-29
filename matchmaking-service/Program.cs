using Scalar.AspNetCore;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring OpenAPI at https://aka.ms/aspnet/openapi
builder.Services.AddOpenApi();

builder.Services.AddKeyedSingleton<IQueueService, InMemoryQueueService>("survivor");
builder.Services.AddKeyedSingleton<IQueueService, InMemoryQueueService>("killer");
builder.Services.AddSingleton<IQueueManager, QueueManager>();
builder.Services.AddSingleton<IMatchStore, InMemoryMatchStore>();
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


