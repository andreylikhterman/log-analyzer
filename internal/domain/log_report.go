package domain

import "time"

type LogReport struct {
	FileNames                []string
	StartDate                time.Time
	EndDate                  time.Time
	TotalRequests            int
	AvgBodySize              int
	AvgTimeBetweenRequests   time.Duration
	Percentile95Size         int
	RequestedResources       map[string]int
	SortedRequestedResources []string
	ResponseCodes            map[int]ResponseCode
	SortedResponseCodes      []int
	TopIPAddresses           []IPCount
}

type ResponseCode struct {
	Name  string
	Count int
}

type IPCount struct {
	IP    string
	Count int
}
