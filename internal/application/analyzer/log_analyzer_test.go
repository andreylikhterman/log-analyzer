package analyzer

import (
	"testing"
	"time"

	domain "analyzer/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestLogRecords() []domain.LogRecord {
	return []domain.LogRecord{
		{
			RemoteAddr:      "192.168.1.1",
			RemoteUser:      "user1",
			TimeLocal:       time.Date(2023, 10, 15, 10, 0, 0, 0, time.UTC),
			Method:          "GET",
			URL:             "/api/data",
			ProtocolVersion: "HTTP/1.1",
			Status:          200,
			BodyBytesSent:   512,
		},
		{
			RemoteAddr:      "192.168.1.2",
			RemoteUser:      "user2",
			TimeLocal:       time.Date(2023, 10, 16, 10, 0, 0, 0, time.UTC),
			Method:          "POST",
			URL:             "/api/data",
			ProtocolVersion: "HTTP/1.1",
			Status:          404,
			BodyBytesSent:   0,
		},
		{
			RemoteAddr:      "192.168.1.2",
			RemoteUser:      "user3",
			TimeLocal:       time.Date(2023, 10, 17, 10, 0, 0, 0, time.UTC),
			Method:          "GET",
			URL:             "/api/data",
			ProtocolVersion: "HTTP/1.1",
			Status:          200,
			BodyBytesSent:   1024,
		},
		{
			RemoteAddr:      "192.168.1.1",
			RemoteUser:      "user1",
			TimeLocal:       time.Date(2023, 10, 18, 10, 0, 0, 0, time.UTC),
			Method:          "DELETE",
			URL:             "/api/otherdata",
			ProtocolVersion: "HTTP/1.1",
			Status:          200,
			BodyBytesSent:   0,
		},
		{
			RemoteAddr:      "192.168.1.1",
			RemoteUser:      "user2",
			TimeLocal:       time.Date(2023, 10, 19, 10, 0, 0, 0, time.UTC),
			Method:          "GET",
			URL:             "/api/otherdata",
			ProtocolVersion: "HTTP/1.1",
			Status:          404,
			BodyBytesSent:   0,
		},
	}
}

func TestLogAnalyzer(t *testing.T) {
	analyzer := NewLogAnalyzer()
	records := createTestLogRecords()
	config := &domain.Config{}

	report, err := analyzer.Analyze(records, config)
	require.NoError(t, err)

	t.Run("TotalRequests", func(t *testing.T) {
		assert.Equal(t, 5, report.TotalRequests, "Общее количество запросов должно быть 5")
	})

	t.Run("AverageBodySize", func(t *testing.T) {
		assert.Equal(t, 307, report.AvgBodySize, "Средний размер тела запросов должен быть 507 байт")
	})

	t.Run("AvgTimeBetweenRequests", func(t *testing.T) {
		assert.Equal(t, 24*time.Hour, report.AvgTimeBetweenRequests, "Среднее время между запросами должно быть 24 часа")
	})

	t.Run("ResponseCodes", func(t *testing.T) {
		assert.Len(t, report.ResponseCodes, 2, "Должно быть 3 уникальных кода ответа")
		assert.Equal(t, report.ResponseCodes[200].Count, 3, "Количество кода 200 должно быть 3")
		assert.Equal(t, report.ResponseCodes[404].Count, 2, "Количество кода 404 должно быть 2")
	})

	t.Run("TopIPAddresses", func(t *testing.T) {
		assert.Len(t, report.TopIPAddresses, 2, "Должно быть 2 верхних IP адреса")
		assert.Equal(t, report.TopIPAddresses[0].IP, "192.168.1.1", "Топ IP адрес должен быть 192.168.1.1")
		assert.Equal(t, report.TopIPAddresses[1].IP, "192.168.1.2", "Топ IP адрес должен быть 192.168.1.2")
	})

	t.Run("Percentile95Size", func(t *testing.T) {
		assert.Equal(t, 512, report.Percentile95Size, "95-й процентиль должен быть 1024")
	})

	t.Run("SortedRequestedResources", func(t *testing.T) {
		assert.Equal(t, []string{"/api/data", "/api/otherdata"},
			report.SortedRequestedResources, "Запрашиваемые ресурсы должны быть отсортированы")
	})
}
