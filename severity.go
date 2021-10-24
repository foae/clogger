package clogger

import "sync"

// Severity ...
type Severity int

type sev struct {
	m  map[Severity]string
	mu *sync.Mutex
}

const (
	SeverityDebug Severity = iota
	SeverityInfo
	SeverityWarn
	SeverityError
	SeverityCritical
)

var (
	sevMap = sev{
		m: map[Severity]string{
			0: "DEBUG",
			1: "INFO",
			2: "WARN",
			3: "ERROR",
			4: "CRITICAL",
		},
		mu: &sync.Mutex{},
	}
)

// String implements the stringer interface.
func (s Severity) String() string {
	sevMap.mu.Lock()
	defer sevMap.mu.Unlock()

	return sevMap.m[s]
}
