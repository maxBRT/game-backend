using System.Collections.Concurrent;
using System.Diagnostics;

public class MetricsService
{
    private long _systemErrorCount;
    private long _totalRequestCount;
    private long _clientErrorCount;
    private long _playerInQueueCount;
    private long _playerMatchedCount;

    private readonly ConcurrentDictionary<string, OperationMetrics> _operationMetrics = new();

    public void IncrementPlayerInQueueCount()
    {
        Interlocked.Increment(ref _playerInQueueCount);
    }

    public void IncrementPlayerMatchedCount()
    {
        Interlocked.Increment(ref _playerMatchedCount);
    }

    public long GetPlayerInQueueCount() => Interlocked.Read(ref _playerInQueueCount);
    public long GetPlayerMatchedCount() => Interlocked.Read(ref _playerMatchedCount);

    public void RecordRequest(string operation, double elapsedMs, bool isSystemError, bool isClientError = false)
    {
        Interlocked.Increment(ref _totalRequestCount);
        if (isSystemError) Interlocked.Increment(ref _systemErrorCount);
        if (isClientError) Interlocked.Increment(ref _clientErrorCount);

        var metrics = _operationMetrics.GetOrAdd(operation, _ => new OperationMetrics());
        metrics.RecordLatency(elapsedMs, isSystemError, isClientError);
    }

    public long GetTotalRequestCount() => Interlocked.Read(ref _totalRequestCount);
    public long GetSystemErrorCount() => Interlocked.Read(ref _systemErrorCount);
    public long GetClientErrorCount() => Interlocked.Read(ref _clientErrorCount);

    public OperationMetrics? GetOperationMetrics(string operation)
    {
        _operationMetrics.TryGetValue(operation, out var metrics);
        return metrics;
    }

    public IReadOnlyDictionary<string, OperationMetrics> GetAllOperationMetrics() => _operationMetrics;
}

public class OperationMetrics
{
    private readonly ConcurrentBag<double> _latencies = [];
    private long _count;
    private long _errorCount;
    private long _clientErrorCount;

    public void RecordLatency(double elapsedMs, bool isSystemError, bool isClientError)
    {
        _latencies.Add(elapsedMs);
        Interlocked.Increment(ref _count);
        if (isSystemError) Interlocked.Increment(ref _errorCount);
        if (isClientError) Interlocked.Increment(ref _clientErrorCount);
    }

    public long Count => Interlocked.Read(ref _count);
    public long ErrorCount => Interlocked.Read(ref _errorCount);
    public long ClientErrorCount => Interlocked.Read(ref _clientErrorCount);
    public long SuccessCount => Count - ErrorCount - ClientErrorCount;

    public LatencyStats GetLatencyStats()
    {
        var sorted = _latencies.OrderBy(x => x).ToList();
        if (sorted.Count == 0)
            return new LatencyStats(0, 0, 0);

        return new LatencyStats(
            Min: sorted[0],
            Max: sorted[^1],
            Avg: sorted.Average()
            );
    }


}

public record LatencyStats(double Min, double Max, double Avg);
