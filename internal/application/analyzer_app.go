package application

import (
	domain "analyzer/internal/domain"
	"log"
)

type ParserLog interface {
	Parse(config *domain.Config) ([]domain.LogRecord, error)
}

type FilterLog interface {
	Filter(records []domain.LogRecord, config *domain.Config) []domain.LogRecord
}

type LogAnalyzer interface {
	Analyze(records []domain.LogRecord, config *domain.Config) (domain.LogReport, error)
}

type Formatter interface {
	Format(report *domain.LogReport, format string) (string, error)
}

type Saver interface {
	Save(output, format, name string) error
}

type AnalyzerApp struct {
	LogParser   ParserLog
	LogFilter   FilterLog
	LogAnalyzer LogAnalyzer
	Formatter   Formatter
	Saver       Saver
}

func (app *AnalyzerApp) Run(config *domain.Config) {
	logRecords, err := app.LogParser.Parse(config)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	logRecords = app.LogFilter.Filter(logRecords, config)

	logReport, err := app.LogAnalyzer.Analyze(logRecords, config)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	output, err := app.Formatter.Format(&logReport, config.Format)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	err = app.Saver.Save(output, "analyze", config.Format)
	if err != nil {
		log.Fatalf("Ошибка: %v", err)
	}
}
