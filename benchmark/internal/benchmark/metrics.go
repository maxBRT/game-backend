package benchmark

import (
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	survivorsWaiting      atomic.Int64
	killersWaiting        atomic.Int64
	matchesStarted        atomic.Int64
	purchases             atomic.Int64
	errorCount            atomic.Int64
	mu                    sync.Mutex
	matchLatancies        []time.Duration
	purchasesLatancies    []time.Duration
	startTime             time.Time
	lastMatchCount        int64
	lastPurchaseCount     int64
	lastRateUpdate        time.Time
	peakMatchPerSecond    float64
	peakPurchasePerSecond float64
}

func NewMetrics() *Metrics {
	return &Metrics{
		survivorsWaiting:      atomic.Int64{},
		killersWaiting:        atomic.Int64{},
		matchesStarted:        atomic.Int64{},
		purchases:             atomic.Int64{},
		errorCount:            atomic.Int64{},
		matchLatancies:        make([]time.Duration, 0, 1000),
		purchasesLatancies:    make([]time.Duration, 0, 1000),
		startTime:             time.Now(),
		lastMatchCount:        0,
		lastPurchaseCount:     0,
		lastRateUpdate:        time.Now(),
		peakMatchPerSecond:    0,
		peakPurchasePerSecond: 0,
	}
}

func (m *Metrics) IncrementSurvivorsWaiting() {
	m.survivorsWaiting.Add(1)
}

func (m *Metrics) DecrementSurvivorsWaiting() {
	m.survivorsWaiting.Add(-1)
}

func (m *Metrics) IncrementKillersWaiting() {
	m.killersWaiting.Add(1)
}

func (m *Metrics) DecrementKillersWaiting() {
	m.killersWaiting.Add(-1)
}

func (m *Metrics) RecordMatch(timeToMatch time.Duration) {
	m.matchesStarted.Add(1)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.matchLatancies = append(m.matchLatancies, timeToMatch)
}

func (m *Metrics) RecordPurchase(timeToPurchase time.Duration) {
	m.purchases.Add(1)

	m.mu.Lock()
	defer m.mu.Unlock()

	m.purchasesLatancies = append(m.purchasesLatancies, timeToPurchase)
}

func (m *Metrics) IncrementErrorCount() {
	m.errorCount.Add(1)
}

func (m *Metrics) Update() {
	now := time.Now()
	elapsed := now.Sub(m.startTime)
	if elapsed.Seconds() < 1 {
		return
	}

	matchesPerSecond := float64(m.matchesStarted.Load()) / elapsed.Seconds()
	purchasesPerSecond := float64(m.purchases.Load()) / elapsed.Seconds()

	if matchesPerSecond > m.peakMatchPerSecond {
		m.peakMatchPerSecond = matchesPerSecond
	}
	if purchasesPerSecond > m.peakPurchasePerSecond {
		m.peakPurchasePerSecond = purchasesPerSecond
	}

	m.lastMatchCount = m.matchesStarted.Load()
	m.lastPurchaseCount = m.purchases.Load()
	m.lastRateUpdate = now
}

func (m *Metrics) Percentile(samples []time.Duration, percentile float64) time.Duration {
	if len(samples) == 0 {
		return 0
	}
	samplesCopy := make([]time.Duration, len(samples))
	copy(samplesCopy, samples)
	slices.Sort(samplesCopy)
	index := int(float64(len(samplesCopy)) * percentile)
	return samplesCopy[index]
}

type MetricsSnapshot struct {
	SurvivorsWaiting      int64
	KillersWaiting        int64
	MatchesStarted        int64
	MatchesPerSecond      float64
	AvgTimeToMatch        time.Duration
	P95TimeToMatch        time.Duration
	P99TimeToMatch        time.Duration
	Purchases             int64
	PurchasesPerSecond    float64
	AvgPurchaseLatency    time.Duration
	P95PurchaseLatency    time.Duration
	P99PurchaseLatency    time.Duration
	ErrorCount            int64
	ErrorRate             float64
	TotalRequests         int64
	ActiveConnections     int64
	PeekMatchPerSecond    float64
	PeekPurchasePerSecond float64
	Elapsed               time.Duration
}

func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	elapsed := time.Since(m.startTime)
	return MetricsSnapshot{
		SurvivorsWaiting:      m.survivorsWaiting.Load(),
		KillersWaiting:        m.killersWaiting.Load(),
		MatchesStarted:        m.matchesStarted.Load(),
		MatchesPerSecond:      float64(m.matchesStarted.Load()) / elapsed.Seconds(),
		AvgTimeToMatch:        Average(m.matchLatancies),
		P95TimeToMatch:        m.Percentile(m.matchLatancies, 0.95),
		P99TimeToMatch:        m.Percentile(m.matchLatancies, 0.99),
		Purchases:             m.purchases.Load(),
		PurchasesPerSecond:    float64(m.purchases.Load()) / elapsed.Seconds(),
		AvgPurchaseLatency:    Average(m.purchasesLatancies),
		P95PurchaseLatency:    m.Percentile(m.purchasesLatancies, 0.95),
		P99PurchaseLatency:    m.Percentile(m.purchasesLatancies, 0.99),
		ErrorCount:            m.errorCount.Load(),
		ErrorRate:             float64(m.errorCount.Load()) / elapsed.Seconds(),
		TotalRequests:         m.matchesStarted.Load() + m.purchases.Load() + m.survivorsWaiting.Load() + m.killersWaiting.Load(),
		ActiveConnections:     m.survivorsWaiting.Load() + m.killersWaiting.Load(),
		PeekMatchPerSecond:    m.peakMatchPerSecond,
		PeekPurchasePerSecond: m.peakPurchasePerSecond,
		Elapsed:               elapsed,
	}
}

func Average(samples []time.Duration) time.Duration {
	if len(samples) == 0 {
		return 0
	}
	samplesCopy := make([]time.Duration, len(samples))
	copy(samplesCopy, samples)
	slices.Sort(samplesCopy)
	return samplesCopy[len(samplesCopy)/2]
}
