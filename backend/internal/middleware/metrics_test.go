package middleware

import (
	"fmt"
	"testing"
	"time"
)

func TestMetrics_RecordRequest(t *testing.T) {
	// Create a fresh metrics instance for testing
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Test recording a successful request
	metrics.RecordRequest("GET", "/api/restaurants", 200, 10*time.Millisecond)

	if metrics.TotalRequests != 1 {
		t.Errorf("Expected TotalRequests to be 1, got %d", metrics.TotalRequests)
	}

	if metrics.TotalErrors != 0 {
		t.Errorf("Expected TotalErrors to be 0, got %d", metrics.TotalErrors)
	}

	// Test recording an error request
	metrics.RecordRequest("POST", "/api/restaurants", 500, 50*time.Millisecond)

	if metrics.TotalRequests != 2 {
		t.Errorf("Expected TotalRequests to be 2, got %d", metrics.TotalRequests)
	}

	if metrics.TotalErrors != 1 {
		t.Errorf("Expected TotalErrors to be 1, got %d", metrics.TotalErrors)
	}
}

func TestMetrics_RecordRequest_ByMethod(t *testing.T) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Record multiple requests with different methods
	metrics.RecordRequest("GET", "/api/restaurants", 200, 10*time.Millisecond)
	metrics.RecordRequest("GET", "/api/categories", 200, 5*time.Millisecond)
	metrics.RecordRequest("POST", "/api/restaurants", 201, 20*time.Millisecond)
	metrics.RecordRequest("PUT", "/api/restaurants/1", 200, 15*time.Millisecond)

	// Verify counts by method
	if metrics.RequestsByMethod["GET"] == nil {
		t.Error("Expected GET requests to be recorded")
	} else if *metrics.RequestsByMethod["GET"] != 2 {
		t.Errorf("Expected 2 GET requests, got %d", *metrics.RequestsByMethod["GET"])
	}

	if metrics.RequestsByMethod["POST"] == nil {
		t.Error("Expected POST requests to be recorded")
	} else if *metrics.RequestsByMethod["POST"] != 1 {
		t.Errorf("Expected 1 POST request, got %d", *metrics.RequestsByMethod["POST"])
	}

	if metrics.RequestsByMethod["PUT"] == nil {
		t.Error("Expected PUT requests to be recorded")
	} else if *metrics.RequestsByMethod["PUT"] != 1 {
		t.Errorf("Expected 1 PUT request, got %d", *metrics.RequestsByMethod["PUT"])
	}
}

func TestMetrics_RecordRequest_ByStatus(t *testing.T) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Record requests with different status codes
	metrics.RecordRequest("GET", "/api/restaurants", 200, 10*time.Millisecond)
	metrics.RecordRequest("GET", "/api/restaurants/999", 404, 5*time.Millisecond)
	metrics.RecordRequest("POST", "/api/restaurants", 500, 20*time.Millisecond)

	// Verify counts by status
	if metrics.RequestsByStatus[200] == nil {
		t.Error("Expected 200 status to be recorded")
	} else if *metrics.RequestsByStatus[200] != 1 {
		t.Errorf("Expected 1 request with status 200, got %d", *metrics.RequestsByStatus[200])
	}

	if metrics.RequestsByStatus[404] == nil {
		t.Error("Expected 404 status to be recorded")
	} else if *metrics.RequestsByStatus[404] != 1 {
		t.Errorf("Expected 1 request with status 404, got %d", *metrics.RequestsByStatus[404])
	}

	if metrics.RequestsByStatus[500] == nil {
		t.Error("Expected 500 status to be recorded")
	} else if *metrics.RequestsByStatus[500] != 1 {
		t.Errorf("Expected 1 request with status 500, got %d", *metrics.RequestsByStatus[500])
	}

	// Verify error count
	if metrics.TotalErrors != 1 {
		t.Errorf("Expected 1 error (status >= 500), got %d", metrics.TotalErrors)
	}
}

func TestMetrics_CalculatePercentiles(t *testing.T) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Add response times in a known pattern
	// 100 requests from 1ms to 100ms
	for i := 1; i <= 100; i++ {
		metrics.RecordRequest("GET", "/api/test", 200, time.Duration(i)*time.Millisecond)
	}

	p50, p95, p99 := metrics.calculatePercentiles()

	// p50 (median) should be around 50ms
	if p50 < 40*time.Millisecond || p50 > 60*time.Millisecond {
		t.Errorf("Expected p50 to be around 50ms, got %v", p50)
	}

	// p95 should be around 95ms
	if p95 < 85*time.Millisecond || p95 > 100*time.Millisecond {
		t.Errorf("Expected p95 to be around 95ms, got %v", p95)
	}

	// p99 should be around 99ms
	if p99 < 90*time.Millisecond || p99 > 100*time.Millisecond {
		t.Errorf("Expected p99 to be around 99ms, got %v", p99)
	}
}

func TestMetrics_CalculatePercentiles_EmptyData(t *testing.T) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	p50, p95, p99 := metrics.calculatePercentiles()

	if p50 != 0 || p95 != 0 || p99 != 0 {
		t.Error("Expected all percentiles to be 0 for empty data")
	}
}

func TestMetrics_GetStats(t *testing.T) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Record some test data
	metrics.RecordRequest("GET", "/api/restaurants", 200, 10*time.Millisecond)
	metrics.RecordRequest("POST", "/api/restaurants", 201, 20*time.Millisecond)
	metrics.RecordRequest("GET", "/api/categories", 200, 15*time.Millisecond)
	metrics.RecordRequest("DELETE", "/api/restaurants/1", 500, 100*time.Millisecond)

	stats := metrics.GetStats()

	// Verify total requests
	if totalReq, ok := stats["total_requests"].(uint64); ok {
		if totalReq != 4 {
			t.Errorf("Expected total_requests to be 4, got %d", totalReq)
		}
	} else {
		t.Error("total_requests not found in stats")
	}

	// Verify total errors
	if totalErr, ok := stats["total_errors"].(uint64); ok {
		if totalErr != 1 {
			t.Errorf("Expected total_errors to be 1, got %d", totalErr)
		}
	} else {
		t.Error("total_errors not found in stats")
	}

	// Verify requests by method
	if byMethod, ok := stats["requests_by_method"].(map[string]uint64); ok {
		if byMethod["GET"] != 2 {
			t.Errorf("Expected 2 GET requests, got %d", byMethod["GET"])
		}
		if byMethod["POST"] != 1 {
			t.Errorf("Expected 1 POST request, got %d", byMethod["POST"])
		}
		if byMethod["DELETE"] != 1 {
			t.Errorf("Expected 1 DELETE request, got %d", byMethod["DELETE"])
		}
	} else {
		t.Error("requests_by_method not found in stats")
	}

	// Verify requests by status
	if byStatus, ok := stats["requests_by_status"].(map[int]uint64); ok {
		if byStatus[200] != 2 {
			t.Errorf("Expected 2 requests with status 200, got %d", byStatus[200])
		}
		if byStatus[201] != 1 {
			t.Errorf("Expected 1 request with status 201, got %d", byStatus[201])
		}
		if byStatus[500] != 1 {
			t.Errorf("Expected 1 request with status 500, got %d", byStatus[500])
		}
	} else {
		t.Error("requests_by_status not found in stats")
	}

	// Verify avg_response_time exists
	if _, ok := stats["avg_response_time"]; !ok {
		t.Error("avg_response_time not found in stats")
	}

	// Verify percentiles exist
	if _, ok := stats["p50_response_time"]; !ok {
		t.Error("p50_response_time not found in stats")
	}
	if _, ok := stats["p95_response_time"]; !ok {
		t.Error("p95_response_time not found in stats")
	}
	if _, ok := stats["p99_response_time"]; !ok {
		t.Error("p99_response_time not found in stats")
	}

	// Verify uptime exists
	if _, ok := stats["uptime"]; !ok {
		t.Error("uptime not found in stats")
	}
}

func TestMetrics_Reset(t *testing.T) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Record some data
	metrics.RecordRequest("GET", "/api/restaurants", 200, 10*time.Millisecond)
	metrics.RecordRequest("POST", "/api/restaurants", 201, 20*time.Millisecond)

	// Verify data is present
	if metrics.TotalRequests != 2 {
		t.Errorf("Expected TotalRequests to be 2 before reset, got %d", metrics.TotalRequests)
	}

	// Reset metrics
	metrics.Reset()

	// Verify all metrics are reset
	if metrics.TotalRequests != 0 {
		t.Errorf("Expected TotalRequests to be 0 after reset, got %d", metrics.TotalRequests)
	}
	if metrics.TotalErrors != 0 {
		t.Errorf("Expected TotalErrors to be 0 after reset, got %d", metrics.TotalErrors)
	}
	if len(metrics.RequestsByMethod) != 0 {
		t.Errorf("Expected RequestsByMethod to be empty after reset, got %d items", len(metrics.RequestsByMethod))
	}
	if len(metrics.RequestsByPath) != 0 {
		t.Errorf("Expected RequestsByPath to be empty after reset, got %d items", len(metrics.RequestsByPath))
	}
	if len(metrics.RequestsByStatus) != 0 {
		t.Errorf("Expected RequestsByStatus to be empty after reset, got %d items", len(metrics.RequestsByStatus))
	}
	if len(metrics.ResponseTimes) != 0 {
		t.Errorf("Expected ResponseTimes to be empty after reset, got %d items", len(metrics.ResponseTimes))
	}
}

func TestMetrics_PathLimit(t *testing.T) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Try to record 150 different paths (limit is 100)
	for i := 0; i < 150; i++ {
		path := fmt.Sprintf("/api/test/%d", i)
		metrics.RecordRequest("GET", path, 200, 10*time.Millisecond)
	}

	// Verify path count is limited to 100
	if len(metrics.RequestsByPath) > 100 {
		t.Errorf("Expected RequestsByPath to be limited to 100, got %d", len(metrics.RequestsByPath))
	}
}

func TestMetrics_ResponseTimeLimit(t *testing.T) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Try to record 1500 response times (limit is 1000)
	for i := 0; i < 1500; i++ {
		metrics.RecordRequest("GET", "/api/test", 200, time.Duration(i)*time.Millisecond)
	}

	// Verify response times are limited to 1000
	if len(metrics.ResponseTimes) > 1000 {
		t.Errorf("Expected ResponseTimes to be limited to 1000, got %d", len(metrics.ResponseTimes))
	}
}

func TestGetMetrics_Singleton(t *testing.T) {
	// Get metrics instance twice
	m1 := GetMetrics()
	m2 := GetMetrics()

	// Verify they're the same instance
	if m1 != m2 {
		t.Error("Expected GetMetrics to return the same singleton instance")
	}

	// Record the current request count
	initialCount := m1.TotalRequests

	// Record a request on one instance
	m1.RecordRequest("GET", "/api/test-singleton", 200, 10*time.Millisecond)

	// Verify it's reflected on the other instance
	if m2.TotalRequests != initialCount+1 {
		t.Errorf("Metrics not shared between singleton instances. Expected %d, got %d", initialCount+1, m2.TotalRequests)
	}
}

func BenchmarkMetrics_RecordRequest(b *testing.B) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordRequest("GET", "/api/restaurants", 200, 10*time.Millisecond)
	}
}

func BenchmarkMetrics_CalculatePercentiles(b *testing.B) {
	metrics := &Metrics{
		RequestsByMethod: make(map[string]*uint64),
		RequestsByPath:   make(map[string]*uint64),
		RequestsByStatus: make(map[int]*uint64),
		ResponseTimes:    make([]time.Duration, 0, 1000),
		lastLogTime:      time.Now(),
	}

	// Add 1000 response times
	for i := 1; i <= 1000; i++ {
		metrics.RecordRequest("GET", "/api/test", 200, time.Duration(i)*time.Millisecond)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.calculatePercentiles()
	}
}
