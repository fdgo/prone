package version

import (
	"fmt"
	"os"
	"time"
)

type ServiceInfo struct {
	Version      string    `json:"version,omitempty"`
	Binary       string    `json:"binary,omitempty"`
	CompleteDate time.Time `json:"completeDate,omitempty"`
	StartTime    time.Time `json:"startTime,omitempty"`
	RunTime      string    `json:"runTime,omitempty"`
}

var (
	_DATE_       = "2017-12-12 00:00:01"
	_VERSION_    = "unkown"
	completeDate time.Time
	version      string
	binary       string
	startTime    time.Time
)

func init() {
	completeDate, _ = time.Parse("2006-01-02T15:04:05Z0700", _DATE_)
	version = _VERSION_
	binary = os.Args[0]
	startTime = time.Now()
}

func CompleteDate() time.Time {
	return completeDate
}

func Version() string {
	return version
}

func Binary() string {
	return binary
}

func RunTime() time.Duration {
	return time.Since(startTime)
}

func StartTime() time.Time {
	return startTime
}

func ServicesInfo() *ServiceInfo {
	return &ServiceInfo{
		Version:      Version(),
		Binary:       Binary(),
		CompleteDate: CompleteDate(),
		StartTime:    StartTime(),
		RunTime:      fmt.Sprintf("%v", RunTime()),
	}
}
