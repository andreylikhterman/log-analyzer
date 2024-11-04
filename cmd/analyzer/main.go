package main

import (
	application "analyzer/internal/application"
	analyzer "analyzer/internal/application/analyzer"
	filter "analyzer/internal/application/filter"
	formatter "analyzer/internal/application/formatter"
	parsers "analyzer/internal/application/parsers"
	saver "analyzer/internal/application/saver"
	input "analyzer/internal/infrastructure/input"
)

func main() {
	requestTemplate := input.RequestTemplate{
		Name:          "analyzer",
		RequeredFlags: []string{"path"},
		OptionalFlags: []string{"from", "to", "format", "filter-field", "filter-value"},
	}

	request := input.Request(requestTemplate)

	parser := parsers.NewParserRequest()
	config := parser.Parse(request)

	app := application.AnalyzerApp{
		LogParser:   parsers.NewLogParser(),
		LogAnalyzer: analyzer.NewLogAnalyzer(),
		LogFilter:   filter.NewLogFilter(),
		Formatter:   formatter.NewFormatter(),
		Saver:       saver.NewSaver(),
	}

	app.Run(&config)
}
