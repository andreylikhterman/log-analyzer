package parsers

import (
	"analyzer/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLogLine_ValidLine(t *testing.T) {
	parser := NewLogParser()
	line := `127.0.0.1 - - [12/Oct/2023:14:32:00 +0000] "GET /index.html HTTP/1.1" 200 1024 "http://example.com" "Mozilla/5.0"`

	expected := domain.LogRecord{
		RemoteAddr:      "127.0.0.1",
		RemoteUser:      "-",
		TimeLocal:       time.Date(2023, 10, 12, 14, 32, 0, 0, time.UTC),
		Method:          "GET",
		URL:             "/index.html",
		ProtocolVersion: "HTTP/1.1",
		Status:          200,
		BodyBytesSent:   1024,
		Referer:         "http://example.com",
		UserAgent:       "Mozilla/5.0",
	}

	logRecord, err := parser.parseLogLine(line)

	require.NoError(t, err)
	assert.Equal(t, expected, logRecord)
}

func TestParseLogLine_InvalidLine(t *testing.T) {
	parser := NewLogParser()
	line := `Invalid log line format`

	_, err := parser.parseLogLine(line)
	assert.Error(t, err)
}

func TestParseLogs_ValidLogs(t *testing.T) {
	parser := NewLogParser()
	logLines := []string{
		`127.0.0.1 - - [12/Oct/2023:14:32:00 +0000] "GET /index.html HTTP/1.1" 200 1024 "http://example.com" "Mozilla/5.0"`,
		`127.0.0.1 - - [12/Oct/2023:15:32:00 +0000] "POST /submit HTTP/1.1" 404 512 "http://example.com" "Mozilla/5.0*"`,
		`127.0.0.1 - - [12/Oct/2023:16:32:00 +0000] "PUT /image.jpg HTTP/1.1" 200 256 "http://example.com" "Mozilla*"`,
	}

	logRecords, err := parser.parseLogs(logLines)
	require.NoError(t, err)
	assert.Len(t, logRecords, 3)
}
