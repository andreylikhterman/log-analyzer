package filter

import (
	domain "analyzer/internal/domain"
	"strconv"
)

type LogFilter struct{}

func NewLogFilter() *LogFilter {
	return &LogFilter{}
}

func (filter *LogFilter) Filter(records []domain.LogRecord, config *domain.Config) []domain.LogRecord {
	filteredRecords := make([]domain.LogRecord, 0)

	for ind := range records {
		record := &records[ind]
		if filter.checkFilterFields(record, config) && filter.checkTime(record, config) {
			filteredRecords = append(filteredRecords, *record)
		}
	}

	return filteredRecords
}

func (filter *LogFilter) checkFilterFields(record *domain.LogRecord, config *domain.Config) bool {
	switch config.FilterField {
	case "agent":
		if config.FilterValue.MatchString(record.UserAgent) {
			return true
		}
	case "address":
		if config.FilterValue.MatchString(record.RemoteAddr) {
			return true
		}
	case "user":
		if config.FilterValue.MatchString(record.RemoteUser) {
			return true
		}
	case "method":
		if config.FilterValue.MatchString(record.Method) {
			return true
		}
	case "url":
		if config.FilterValue.MatchString(record.URL) {
			return true
		}
	case "protocol":
		if config.FilterValue.MatchString(record.ProtocolVersion) {
			return true
		}
	case "status":
		if config.FilterValue.MatchString(strconv.Itoa(record.Status)) {
			return true
		}
	case "referer":
		if config.FilterValue.MatchString(record.Referer) {
			return true
		}
	case "":
		return true
	}

	return false
}

func (filter *LogFilter) checkTime(logRecord *domain.LogRecord, config *domain.Config) bool {
	if config.From.IsZero() || logRecord.TimeLocal.After(config.From) {
		if config.To.IsZero() || logRecord.TimeLocal.Before(config.To) {
			return true
		}
	}

	return false
}
