using Microsoft.EntityFrameworkCore;
using player_service.Data;
using player_service.Services;
using Scalar.AspNetCore;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring OpenAPI at https://aka.ms/aspnet/openapi
builder.Services.AddOpenApi();
builder.Services.AddDbContext<AppDbContext>();
builder.Services.AddScoped<IPlayerService, PlayerService>();
builder.Services.AddScoped<SeedingService>();
builder.Services.AddControllers();

var app = builder.Build();

using var scope = app.Services.CreateScope();
{
    var db = scope.ServiceProvider.GetRequiredService<AppDbContext>();
    await db.Database.MigrateAsync();
    SeedingService seedingService = scope.ServiceProvider.GetRequiredService<SeedingService>();
    await seedingService.Seed();
}

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
    app.MapScalarApiReference();
}

app.UseHttpsRedirection();
app.MapControllers();

app.Run();


