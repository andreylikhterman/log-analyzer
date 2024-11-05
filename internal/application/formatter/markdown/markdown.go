package markdown

import (
	"analyzer/internal/domain"
	"analyzer/pkg/output"
	"fmt"
	"strings"
	"time"
)

type Formatter struct{}

func (w *Formatter) WriteGeneralInfo(builder *strings.Builder, report *domain.LogReport) {
	builder.WriteString("## Общая информация\n\n")
	builder.WriteString("| **Метрика**                     | **Значение**              |\n")
	builder.WriteString("|---------------------------------|---------------------------|\n")

	w.writeFileNames(builder, report.FileNames)
	w.writeDate(builder, "Начальная дата", report.StartDate)
	w.writeDate(builder, "Конечная дата", report.EndDate)
	fmt.Fprintf(builder, "| Количество запросов             | %-25s |\n", output.FormatNumber(report.TotalRequests))
	fmt.Fprintf(builder, "| Средний размер ответа           | %-25s |\n", fmt.Sprintf("%sb", output.FormatNumber(report.AvgBodySize)))
	fmt.Fprintf(builder, "| 95p размера ответа              | %-25s |\n", fmt.Sprintf("%sb", output.FormatNumber(report.Percentile95Size)))
	fmt.Fprintf(builder, "| Среднее время между запросами   | %-25s |\n", report.AvgTimeBetweenRequests)
	builder.WriteString("\n")
}

func (w *Formatter) WriteRequestedResources(builder *strings.Builder, report *domain.LogReport) {
	builder.WriteString("## Запрашиваемые ресурсы\n\n")
	builder.WriteString("| **Ресурс**             | **Количество**            |\n")
	builder.WriteString("|------------------------|---------------------------|\n")

	for _, resource := range report.SortedRequestedResources[:min(3, len(report.RequestedResources))] {
		fmt.Fprintf(builder, "| %-22s | %-25s |\n", fmt.Sprintf("`%s`", resource), output.FormatNumber(report.RequestedResources[resource]))
	}

	builder.WriteString("\n")
}

func (w *Formatter) WriteResponseCodes(builder *strings.Builder, report *domain.LogReport) {
	builder.WriteString("## Коды ответа\n\n")
	builder.WriteString("|**Код**| **Имя**               | **Количество**      |\n")
	builder.WriteString("|-------|-----------------------|---------------------|\n")

	for _, code := range report.SortedResponseCodes[:min(3, len(report.ResponseCodes))] {
		fmt.Fprintf(builder, "| %-5d | %-21s | %-19s |\n",
			code, report.ResponseCodes[code].Name, output.FormatNumber(report.ResponseCodes[code].Count))
	}

	builder.WriteString("\n")
}

func (w *Formatter) WriteTopIPAddresses(builder *strings.Builder, report *domain.LogReport) {
	if len(report.TopIPAddresses) == 0 {
		return
	}

	builder.WriteString("## Топ IP-адресов\n\n")
	builder.WriteString("| **IP-адрес**          | **Количество запросов**   |\n")
	builder.WriteString("|-----------------------|---------------------------|\n")

	for _, ipCount := range report.TopIPAddresses {
		fmt.Fprintf(builder, "| %-21s | %-25s |\n", ipCount.IP, output.FormatNumber(ipCount.Count))
	}

	builder.WriteString("\n")
}

func (w *Formatter) writeFileNames(builder *strings.Builder, fileNames []string) {
	for i, fileName := range fileNames {
		if i == 0 {
			fmt.Fprintf(builder, "| Файл(-ы)                        | %-25s |\n", "`"+fileName+"`")
		} else {
			fmt.Fprintf(builder, "|                                 | %-25s |\n", "`"+fileName+"`")
		}
	}
}

func (w *Formatter) writeDate(builder *strings.Builder, label string, date time.Time) {
	if date.IsZero() {
		fmt.Fprintf(builder, "| %-31s | %-25s |\n", label, "-")
	} else {
		fmt.Fprintf(builder, "| %-31s | %-25s |\n", label, date.Format("02.01.2006"))
	}
}
