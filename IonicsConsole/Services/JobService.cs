using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Configuration;
using Npgsql;

namespace IonicsConsole.Services;

public class JobService
{
    private readonly ApiClient _client;
    private readonly NpgsqlDataSource _dataSource;
    private readonly ILogger<JobService> _logger;
    private readonly int _batchSize;
    private readonly string _initialDate;

    public JobService(ApiClient client, NpgsqlDataSource dataSource, IConfiguration config, ILogger<JobService> logger)
    {
        _client = client;
        _dataSource = dataSource;
        _logger = logger;
        _batchSize = config.GetValue<int>("BatchSize", 30);
        _initialDate = config.GetValue("InitialDate", DateTime.UtcNow.ToString("yyyy-MM-dd"));
    }

    public async Task RunAsync(CancellationToken token = default)
    {
        _logger.LogInformation("Running job...");

        var login = await _client.LoginAsync();
        if (login?.Token == null)
        {
            _logger.LogWarning("Failed to authenticate");
            return;
        }

        await using var conn = await _dataSource.OpenConnectionAsync(token);
        int offset = 0;
        while (!token.IsCancellationRequested)
        {
            var records = await LoadRecordsAsync(conn, offset, token);
            if (records.Count == 0)
            {
                _logger.LogInformation("Process completed");
                break;
            }

            var sent = await _client.SendDataAsync(login.Token, records, offset, login.User.CustomerId, login.User.ClientId);
            if (!sent)
            {
                _logger.LogWarning("Failed to send records to API");
                break;
            }

            offset += _batchSize;
        }
    }

    private async Task<List<Dictionary<string, object>>> LoadRecordsAsync(NpgsqlConnection conn, int offset, CancellationToken token)
    {
        await using var cmd = conn.CreateCommand();
        cmd.CommandText = @"SELECT * FROM your_table LIMIT @batch OFFSET @offset"; // replace with real query
        cmd.Parameters.AddWithValue("batch", _batchSize);
        cmd.Parameters.AddWithValue("offset", offset);

        await using var reader = await cmd.ExecuteReaderAsync(token);
        var list = new List<Dictionary<string, object>>();
        while (await reader.ReadAsync(token))
        {
            var record = new Dictionary<string, object>();
            for (int i = 0; i < reader.FieldCount; i++)
            {
                record[reader.GetName(i)] = await reader.IsDBNullAsync(i, token) ? null : reader.GetValue(i);
            }
            list.Add(record);
        }
        return list;
    }
}
