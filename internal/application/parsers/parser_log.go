package parsers

import (
	domain "analyzer/internal/domain"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type LogParser struct {
	LogPattern *regexp.Regexp
}

func NewLogParser() *LogParser {
	return &LogParser{
		LogPattern: regexp.MustCompile(`(?P<remote_addr>\S+) - (?P<remote_user>\S+) ` +
			`\[(?P<time_local>\S+\s\S+)\] "(?P<request>[^"]*)" ` +
			`(?P<status>\d+) (?P<body_bytes_sent>\d+) ` +
			`"(?P<referer>[^"]*)" "(?P<user_agent>[^"]*)"`),
	}
}

func (parser *LogParser) Parse(config *domain.Config) ([]domain.LogRecord, error) {
	logs := make([]string, 0)

	var err error

	switch config.TypePath {
	case "url":
		logs, err = parser.getLogsFromURL(config.Path)
	case "local":
		logs, err = parser.getLogsFromLocal(config.Path)
	}

	if err != nil {
		return nil, err
	}

	logRecords, err := parser.parseLogs(logs)
	if err != nil {
		return nil, err
	}

	return logRecords, err
}

func (parser *LogParser) parseLogs(logs []string) ([]domain.LogRecord, error) {
	logRecords := make([]domain.LogRecord, 0)

	for _, log := range logs {
		logRecord, err := parser.parseLogLine(log)
		if err != nil {
			return nil, err
		}

		logRecords = append(logRecords, logRecord)
	}

	return logRecords, nil
}

func (parser *LogParser) getLogsFromURL(url string) ([]string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить запрос по URL %s: %v", url, err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось получить файл: %s", response.Status)
	}

	scanner := bufio.NewScanner(response.Body)
	logs := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		logs = append(logs, line)
	}

	return logs, nil
}

func (parser *LogParser) getLogsFromLocal(path string) ([]string, error) {
	matches, _ := filepath.Glob(path)
	logs := make([]string, 0)

	for _, file := range matches {
		file, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("не удалось прочитать файл %s: %v", file.Name(), err)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log := scanner.Text()
			logs = append(logs, log)
		}

		file.Close()
	}

	return logs, nil
}

func (parser *LogParser) parseLogLine(line string) (domain.LogRecord, error) {
	matches := parser.LogPattern.FindStringSubmatch(line)
	if matches == nil {
		return domain.LogRecord{}, fmt.Errorf("не удалось разобрать строку: %s", line)
	}

	remoteAddr := matches[1]
	remoteUser := matches[2]
	timeLocalStr := matches[3]
	request := strings.Split(matches[4], " ")
	statusStr := matches[5]
	bodyBytesSentStr := matches[6]
	referer := matches[7]
	userAgent := matches[8]

	timeLocal, err := time.Parse("02/Jan/2006:15:04:05 +0000", timeLocalStr)
	if err != nil {
		return domain.LogRecord{}, fmt.Errorf("не удалось разобрать время: %v", err)
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil {
		return domain.LogRecord{}, fmt.Errorf("не удалось преобразовать status: %v", err)
	}

	bodyBytesSent, err := strconv.Atoi(bodyBytesSentStr)
	if err != nil {
		return domain.LogRecord{}, fmt.Errorf("не удалось преобразовать body_bytes_sent: %v", err)
	}

	return domain.LogRecord{
		RemoteAddr:      remoteAddr,
		RemoteUser:      remoteUser,
		TimeLocal:       timeLocal,
		Method:          request[0],
		URL:             request[1],
		ProtocolVersion: request[2],
		Status:          status,
		BodyBytesSent:   bodyBytesSent,
		Referer:         referer,
		UserAgent:       userAgent,
	}, nil
}
