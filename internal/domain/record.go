package domain

import "time"

type LogRecord struct {
	RemoteAddr      string
	RemoteUser      string
	TimeLocal       time.Time
	Method          string
	URL             string
	ProtocolVersion string
	Status          int
	BodyBytesSent   int
	Referer         string
	UserAgent       string
}
