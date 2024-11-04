package domain

type ParserRequest interface {
	Parse(map[string]string) Config
}

type ParserLog interface {
	Parse(config *Config) ([]LogRecord, error)
}

type FilterLog interface {
	Filter(records []LogRecord, config *Config) []LogRecord
}

type LogAnalyzer interface {
	Analyze(records []LogRecord, config *Config) (LogReport, error)
}

type Formatter interface {
	Format(report *LogReport, format string) (string, error)
}

type Saver interface {
	Save(output, format, name string) error
}
