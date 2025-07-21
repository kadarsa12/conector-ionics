using IonicsConsole.Services;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Npgsql;

var builder = Host.CreateApplicationBuilder(args);

builder.Configuration
    .AddJsonFile("appsettings.json", optional: true, reloadOnChange: true)
    .AddEnvironmentVariables()
    .AddCommandLine(args);

builder.Services.AddHttpClient<ApiClient>();

builder.Services.AddSingleton<NpgsqlDataSource>(sp =>
{
    var connString = builder.Configuration.GetConnectionString("Default");
    return new NpgsqlDataSourceBuilder(connString).Build();
});

builder.Services.AddSingleton<JobService>();

using var host = builder.Build();

var job = host.Services.GetRequiredService<JobService>();
await job.RunAsync();
