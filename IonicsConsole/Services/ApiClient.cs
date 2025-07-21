using System.Net.Http.Json;
using Microsoft.Extensions.Configuration;

namespace IonicsConsole.Services;

public class ApiClient
{
    private readonly HttpClient _httpClient;
    private readonly IConfiguration _config;

    public ApiClient(HttpClient httpClient, IConfiguration config)
    {
        _httpClient = httpClient;
        _config = config;
    }

    public async Task<AuthResponse?> LoginAsync()
    {
        var payload = new
        {
            clientId = _config["ClientId"],
            clientSecret = _config["ClientSecret"]
        };

        var url = $"{_config["BaseUrl"]}{_config["AuthEndpoint"]}";
        var resp = await _httpClient.PostAsJsonAsync(url, payload);
        if (!resp.IsSuccessStatusCode)
        {
            return null;
        }

        return await resp.Content.ReadFromJsonAsync<AuthResponse>();
    }

    public async Task<bool> SendDataAsync(string token, IEnumerable<Dictionary<string, object>> records, int offset, int customerId, string clientId)
    {
        var url = $"{_config["BaseUrl"]}{_config["WriterEndpoint"]}";
        var req = new HttpRequestMessage(HttpMethod.Post, url)
        {
            Content = JsonContent.Create(new { records })
        };
        req.Headers.Add("Authorization", $"Bearer {token}");
        req.Headers.Add("CustomerID", customerId.ToString());
        req.Headers.Add("ClientID", clientId);
        req.Headers.Add("DeleteData", offset == 0 ? "1" : "0");

        var resp = await _httpClient.SendAsync(req);
        if (!resp.IsSuccessStatusCode) return false;
        var body = await resp.Content.ReadFromJsonAsync<SendDataResponse>();
        return body?.Message == "success";
    }
}

public record User(int Id, string Name, string ClientId, int CustomerId, string Role);
public record AuthResponse(string Message, User User, string Token);
public record SendDataResponse(string Message);
