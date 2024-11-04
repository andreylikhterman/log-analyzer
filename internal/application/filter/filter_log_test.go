package filter

import (
	"regexp"
	"testing"
	"time"

	domain "analyzer/internal/domain"

	"github.com/stretchr/testify/assert"
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
			Referer:         "http://example.com",
			UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.66 Safari/537.36",
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
			Referer:         "http://example.com",
			UserAgent:       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
		},
		{
			RemoteAddr:      "192.168.1.3",
			RemoteUser:      "user3",
			TimeLocal:       time.Date(2023, 10, 17, 10, 0, 0, 0, time.UTC),
			Method:          "GET",
			URL:             "/api/otherdata",
			ProtocolVersion: "HTTP/1.1",
			Status:          200,
			BodyBytesSent:   1024,
			Referer:         "http://example.com",
			UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.150 Safari/537.36",
		},
		{
			RemoteAddr:      "192.168.1.4",
			RemoteUser:      "user4",
			TimeLocal:       time.Date(2023, 10, 18, 10, 0, 0, 0, time.UTC),
			Method:          "DELETE",
			URL:             "/api/data/1",
			ProtocolVersion: "HTTP/1.1",
			Status:          200,
			BodyBytesSent:   0,
			Referer:         "http://example.com",
			UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36",
		},
		{
			RemoteAddr:      "192.168.1.5",
			RemoteUser:      "user5",
			TimeLocal:       time.Date(2023, 10, 19, 10, 0, 0, 0, time.UTC),
			Method:          "GET",
			URL:             "/api/notfound",
			ProtocolVersion: "HTTP/1.1",
			Status:          404,
			BodyBytesSent:   0,
			Referer:         "http://example.com",
			UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
	}
}

func TestLogFilter(t *testing.T) {
	filter := NewLogFilter()
	records := createTestLogRecords()

	t.Run("FilterByUserAgent", func(t *testing.T) {
		config := &domain.Config{
			FilterField: "agent",
			FilterValue: regexp.MustCompile("Mozilla"),
		}

		filteredRecords := filter.Filter(records, config)
		assert.Len(t, filteredRecords, 5, "Ожидалось 5 записей, удовлетворяющих фильтру по UserAgent")
	})

	t.Run("FilterByRemoteAddr", func(t *testing.T) {
		config := &domain.Config{
			FilterField: "address",
			FilterValue: regexp.MustCompile("192.168.1.1"),
		}

		filteredRecords := filter.Filter(records, config)
		assert.Len(t, filteredRecords, 1, "Ожидалась 1 запись, удовлетворяющая фильтру по RemoteAddr")
		assert.Equal(t, "192.168.1.1", filteredRecords[0].RemoteAddr)
	})

	t.Run("FilterByStatus", func(t *testing.T) {
		config := &domain.Config{
			FilterField: "status",
			FilterValue: regexp.MustCompile("404"),
		}

		filteredRecords := filter.Filter(records, config)
		assert.Len(t, filteredRecords, 2, "Ожидалось 2 записи, удовлетворяющие фильтру по статусу")
	})

	t.Run("FilterByTimeRange", func(t *testing.T) {
		config := &domain.Config{
			From: time.Date(2023, 10, 16, 0, 0, 0, 0, time.UTC),
			To:   time.Date(2023, 10, 18, 23, 59, 59, 0, time.UTC),
		}

		filteredRecords := filter.Filter(records, config)
		assert.Len(t, filteredRecords, 3, "Ожидалось 3 записи, удовлетворяющие фильтру по времени")
	})

	t.Run("FilterByInvalidField", func(t *testing.T) {
		config := &domain.Config{
			FilterField: "referer",
			FilterValue: regexp.MustCompile(".*"),
		}

		filteredRecords := filter.Filter(records, config)
		assert.Len(t, filteredRecords, 5, "Ожидалось 5 записей, так как фильтр не должен ничего отсеивать")
	})
}
