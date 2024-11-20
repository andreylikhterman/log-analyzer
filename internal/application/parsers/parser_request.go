package parsers

import (
	domain "analyzer/internal/domain"
	"log"
)

type ParserRequest interface {
	Parse(map[string]string) domain.Config
}

type SimpleParserRequest struct{}

func NewParserRequest() *SimpleParserRequest {
	return &SimpleParserRequest{}
}

func (parser *SimpleParserRequest) Parse(flags map[string]string) domain.Config {
	config := domain.Config{}

	err := config.AddPath(flags["path"])
	if err != nil {
		log.Fatal(err)
	}

	err = config.AddFrom(flags["from"])
	if err != nil {
		log.Fatal(err)
	}

	err = config.AddTo(flags["to"])
	if err != nil {
		log.Fatal(err)
	}

	err = config.AddFormat(flags["format"])
	if err != nil {
		log.Fatal(err)
	}

	err = config.AddFilterField(flags["filter-field"])
	if err != nil {
		log.Fatal(err)
	}

	err = config.AddFilterValue(flags["filter-value"])
	if err != nil {
		log.Fatal(err)
	}

	return config
}
