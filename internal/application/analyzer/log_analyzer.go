package analyzer

import (
	domain "analyzer/internal/domain"
	"fmt"
	"path/filepath"
	"sort"
	"time"
)

var statusCodes = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",

	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",

	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	306: "Reserved",
	307: "Temporary Redirect",
	308: "Permanent Redirect",

	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	419: "Authentication Timeout",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	449: "Retry With",
	451: "Unavailable For Legal Reasons",
	499: "Client Closed Request",

	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	509: "Bandwidth Limit Exceeded",
	510: "Not Extended",
	511: "Network Authentication Required",
	520: "Unknown Error",
	521: "Web Server Is Down",
	522: "Connection Timed Out",
	523: "Origin Is Unreachable",
	524: "A Timeout Occurred",
	525: "SSL Handshake Failed",
	526: "Invalid SSL Certificate",
}

type LogAnalyzer struct {
	statusCodes map[int]string
}

func NewLogAnalyzer() *LogAnalyzer {
	return &LogAnalyzer{
		statusCodes: statusCodes,
	}
}

func (analyzer *LogAnalyzer) Analyze(records []domain.LogRecord, config *domain.Config) (domain.LogReport, error) {
	report := analyzer.initReport(records, config)

	if len(records) == 0 {
		return *report, fmt.Errorf("нет записей для анализа")
	}

	analyzer.processRecords(records, report)

	return *report, nil
}

func (analyzer *LogAnalyzer) initReport(records []domain.LogRecord, config *domain.Config) *domain.LogReport {
	return &domain.LogReport{
		RequestedResources: make(map[string]int),
		ResponseCodes:      make(map[int]domain.ResponseCode),
		FileNames:          analyzer.getFileNames(config),
		StartDate:          config.From,
		EndDate:            config.To,
		TotalRequests:      len(records),
	}
}

func (analyzer *LogAnalyzer) processRecords(records []domain.LogRecord, report *domain.LogReport) {
	var (
		totalBodySize            int
		bodySizes                = make([]int, 0)
		ipRequests               = make(map[string]int)
		totalDurationBetweenReqs time.Duration
	)

	for ind := range records {
		record := &records[ind]
		totalBodySize += record.BodyBytesSent
		bodySizes = append(bodySizes, record.BodyBytesSent)

		analyzer.updateIPRequests(ipRequests, record.RemoteAddr)
		analyzer.updateRequestedResources(report, record.URL)
		analyzer.updateResponseCodes(report, record.Status)

		totalDurationBetweenReqs = analyzer.calculateAvgRequestTime(records, ind, totalDurationBetweenReqs)
	}

	analyzer.setTopIPAddresses(report, ipRequests)
	report.AvgTimeBetweenRequests = totalDurationBetweenReqs / time.Duration(report.TotalRequests-1)
	report.AvgBodySize = totalBodySize / report.TotalRequests
	report.Percentile95Size = calculatePercentile(bodySizes, 95)
	report.SortedRequestedResources = sortRequestedResources(report.RequestedResources)
	report.SortedResponseCodes = sortResponseCodes(report.ResponseCodes)
}

func (analyzer *LogAnalyzer) updateIPRequests(ipRequests map[string]int, remoteAddr string) {
	ipRequests[remoteAddr]++
}

func (analyzer *LogAnalyzer) updateRequestedResources(report *domain.LogReport, url string) {
	report.RequestedResources[url]++
}

func (analyzer *LogAnalyzer) updateResponseCodes(report *domain.LogReport, status int) {
	if responseCode, exists := report.ResponseCodes[status]; exists {
		responseCode.Count++
		report.ResponseCodes[status] = responseCode
	} else {
		report.ResponseCodes[status] = domain.ResponseCode{
			Name:  analyzer.getStatusName(status),
			Count: 1,
		}
	}
}

func (analyzer *LogAnalyzer) calculateAvgRequestTime(records []domain.LogRecord, index int, totalDuration time.Duration) time.Duration {
	if index > 0 {
		durationBetween := records[index].TimeLocal.Sub(records[index-1].TimeLocal)
		totalDuration += durationBetween
	}

	return totalDuration
}

func (analyzer *LogAnalyzer) getFileNames(config *domain.Config) []string {
	if config.TypePath == "url" {
		return []string{config.Path}
	}

	fileNames := make([]string, 0)

	matches, _ := filepath.Glob(config.Path)

	for _, file := range matches {
		fileNames = append(fileNames, filepath.Base(file))
	}

	return fileNames
}

func (analyzer *LogAnalyzer) setTopIPAddresses(report *domain.LogReport, ipRequests map[string]int) {
	ips := make([]domain.IPCount, 0, len(ipRequests))
	for ip, count := range ipRequests {
		ips = append(ips, domain.IPCount{IP: ip, Count: count})
	}

	sort.Slice(ips, func(i, j int) bool { return ips[i].Count > ips[j].Count })

	for i := 0; i < 3 && i < len(ips); i++ {
		report.TopIPAddresses = append(report.TopIPAddresses, ips[i])
	}
}

func (analyzer *LogAnalyzer) getStatusName(code int) string {
	return analyzer.statusCodes[code]
}

func calculatePercentile(values []int, percentile int) int {
	sort.Ints(values)
	index := (percentile * len(values) / 100) - 1

	if index < 0 {
		index = 0
	}

	return values[index]
}

func sortResponseCodes(codes map[int]domain.ResponseCode) []int {
	type Codes struct {
		code    int
		summary domain.ResponseCode
	}

	sortedCodes := make([]Codes, 0, len(codes))
	for code, summary := range codes {
		sortedCodes = append(sortedCodes, Codes{code: code, summary: summary})
	}

	sort.Slice(sortedCodes, func(i, j int) bool {
		return sortedCodes[i].summary.Count > sortedCodes[j].summary.Count
	})

	sortedCodesInt := make([]int, 0, len(sortedCodes))
	for _, code := range sortedCodes {
		sortedCodesInt = append(sortedCodesInt, code.code)
	}

	return sortedCodesInt
}

func sortRequestedResources(resources map[string]int) []string {
	type Resource struct {
		resource string
		count    int
	}

	sortedResources := make([]Resource, 0, len(resources))
	for resource, count := range resources {
		sortedResources = append(sortedResources, Resource{resource: resource, count: count})
	}

	sort.Slice(sortedResources, func(i, j int) bool {
		return sortedResources[i].count > sortedResources[j].count
	})

	sortedResourcesStr := make([]string, 0, len(sortedResources))
	for _, resource := range sortedResources {
		sortedResourcesStr = append(sortedResourcesStr, resource.resource)
	}

	return sortedResourcesStr
}
