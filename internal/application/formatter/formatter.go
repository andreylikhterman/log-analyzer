package formatter

//go:generate mockery --name=formatWriter --output=./mocks --outpkg=mocks --filename=mock_formatter.go

import (
	adoc "analyzer/internal/application/formatter/adoc"
	markdown "analyzer/internal/application/formatter/markdown"
	"analyzer/internal/domain"
	"fmt"
	"strings"
)

type formatWriter interface {
	WriteGeneralInfo(builder *strings.Builder, report *domain.LogReport)
	WriteRequestedResources(builder *strings.Builder, report *domain.LogReport)
	WriteResponseCodes(builder *strings.Builder, report *domain.LogReport)
	WriteTopIPAddresses(builder *strings.Builder, report *domain.LogReport)
}

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (formatter *Formatter) Format(report *domain.LogReport, format string) (string, error) {
	var writer formatWriter

	switch format {
	case "markdown":
		writer = &markdown.Formatter{}
	case "adoc":
		writer = &adoc.Formatter{}
	default:
		return "", fmt.Errorf("неподдерживаемый формат: %s", format)
	}

	var builder strings.Builder

	writer.WriteGeneralInfo(&builder, report)
	writer.WriteRequestedResources(&builder, report)
	writer.WriteResponseCodes(&builder, report)
	writer.WriteTopIPAddresses(&builder, report)

	return builder.String(), nil
}
