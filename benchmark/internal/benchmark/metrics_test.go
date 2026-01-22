package benchmark

import (
	"sync"
	"testing"
	"time"
)

func TestConcurrentRecording(t *testing.T) {
	m := NewMetrics()
	var wg sync.WaitGroup

	numGoroutines := 100
	recordsPerGoroutine := 1000

	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			for range recordsPerGoroutine {
				m.RecordMatch(time.Duration(100) * time.Millisecond)
				m.RecordPurchase(time.Duration(50) * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	snapshot := m.Snapshot()
	expected := int64(recordsPerGoroutine * numGoroutines)
	if snapshot.MatchesStarted != expected {
		t.Errorf("Recorded %v matches, want %v", snapshot.MatchesStarted, expected)
	}
	if snapshot.Purchases != expected {
		t.Errorf("Recorded %v purchases, want %v", snapshot.Purchases, expected)
	}

}

func TestPercentile(t *testing.T) {
	m := NewMetrics()
	m.RecordMatch(10 * time.Millisecond)
	m.RecordMatch(20 * time.Millisecond)
	m.RecordMatch(30 * time.Millisecond)
	m.RecordMatch(40 * time.Millisecond)
	m.RecordMatch(50 * time.Millisecond)

	result := m.Percentile(m.matchLatancies, 0.10)
	expected := 10 * time.Millisecond
	if result != expected {
		t.Errorf("Percentile(0.25) on unsorted = %v, want %v", result, expected)
	}

	result = m.Percentile(m.matchLatancies, 0.25)
	expected = 20 * time.Millisecond
	if result != expected {
		t.Errorf("Percentile(0.25) on unsorted = %v, want %v", result, expected)
	}

	result = m.Percentile(m.matchLatancies, 0.5)
	expected = 30 * time.Millisecond
	if result != expected {
		t.Errorf("Percentile(0.5) on unsorted = %v, want %v", result, expected)
	}

	result = m.Percentile(m.matchLatancies, 0.75)
	expected = 40 * time.Millisecond
	if result != expected {
		t.Errorf("Percentile(0.9) on unsorted = %v, want %v", result, expected)
	}

	result = m.Percentile(m.matchLatancies, 0.9)
	expected = 50 * time.Millisecond
	if result != expected {
		t.Errorf("Percentile(0.9) on unsorted = %v, want %v", result, expected)
	}

}

func TestSnapshot(t *testing.T) {
	m := NewMetrics()
	m.RecordMatch(10 * time.Millisecond)
	m.RecordMatch(20 * time.Millisecond)
	m.RecordMatch(30 * time.Millisecond)
	m.RecordMatch(40 * time.Millisecond)
	m.RecordMatch(50 * time.Millisecond)
	m.RecordPurchase(10 * time.Millisecond)
	m.RecordPurchase(20 * time.Millisecond)
	m.RecordPurchase(30 * time.Millisecond)
	m.RecordPurchase(40 * time.Millisecond)
	m.RecordPurchase(50 * time.Millisecond)

	snapshot := m.Snapshot()
	expected := int64(5)
	if snapshot.MatchesStarted != expected {
		t.Errorf("Recorded %v matches, want %v", snapshot.MatchesStarted, expected)
	}
	if snapshot.Purchases != expected {
		t.Errorf("Recorded %v purchases, want %v", snapshot.Purchases, expected)
	}
	if snapshot.SurvivorsWaiting != 0 {
		t.Errorf("SurvivorsWaiting = %v, want 0", snapshot.SurvivorsWaiting)
	}
	if snapshot.KillersWaiting != 0 {
		t.Errorf("KillersWaiting = %v, want 0", snapshot.KillersWaiting)
	}
	if snapshot.TotalRequests != 10 {
		t.Errorf("TotalRequests = %v, want 5", snapshot.TotalRequests)
	}
	if snapshot.ActiveConnections != 0 {
		t.Errorf("ActiveConnections = %v, want 0", snapshot.ActiveConnections)
	}
}

func TestRateCalculation(t *testing.T) {
	m := NewMetrics()

	// Record some events
	m.RecordMatch(10 * time.Millisecond)
	m.RecordMatch(20 * time.Millisecond)
	m.RecordPurchase(10 * time.Millisecond)

	// Sleep for a second to ensure elapsed time for rate calculation
	time.Sleep(1 * time.Second)

	// Record more events
	m.RecordMatch(30 * time.Millisecond)
	m.RecordPurchase(20 * time.Millisecond)
	m.RecordPurchase(30 * time.Millisecond)

	snapshot := m.Snapshot()

	// Since we slept for 1 second and recorded 3 matches and 3 purchases
	// the rate should be approximately 3 matches/sec and 3 purchases/sec
	// We'll allow for some floating point inaccuracy.
	if snapshot.MatchesPerSecond < 2.5 || snapshot.MatchesPerSecond > 3.5 {
		t.Errorf("MatchesPerSecond = %v, want approx 3", snapshot.MatchesPerSecond)
	}
	if snapshot.PurchasesPerSecond < 2.5 || snapshot.PurchasesPerSecond > 3.5 {
		t.Errorf("PurchasesPerSecond = %v, want approx 3", snapshot.PurchasesPerSecond)
	}
}
