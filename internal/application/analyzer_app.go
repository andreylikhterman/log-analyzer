package application

import (
	domain "analyzer/internal/domain"
	"log"
)

type AnalyzerApp struct {
	LogParser   domain.ParserLog
	LogFilter   domain.FilterLog
	LogAnalyzer domain.LogAnalyzer
	Formatter   domain.Formatter
	Saver       domain.Saver
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
