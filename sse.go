package tensorzero

import (
	"bufio"
	"io"
	"strings"
)

// SSEEvent represents a Server-Sent Event
type SSEEvent struct {
	Event string
	Data  string
	ID    string
	Retry string
}

// SSEScanner scans Server-Sent Events from a reader
type SSEScanner struct {
	scanner *bufio.Scanner
	event   SSEEvent
	err     error
}

// NewSSEScanner creates a new SSE scanner
func NewSSEScanner(r io.Reader) *SSEScanner {
	return &SSEScanner{
		scanner: bufio.NewScanner(r),
	}
}

// Scan scans the next SSE event
func (s *SSEScanner) Scan() bool {
	s.event = SSEEvent{}
	
	for s.scanner.Scan() {
		line := s.scanner.Text()
		
		// Empty line indicates end of event
		if line == "" {
			return true
		}
		
		// Skip comments
		if strings.HasPrefix(line, ":") {
			continue
		}
		
		// Parse field
		if colonIndex := strings.Index(line, ":"); colonIndex != -1 {
			field := line[:colonIndex]
			value := line[colonIndex+1:]
			
			// Remove leading space from value
			if len(value) > 0 && value[0] == ' ' {
				value = value[1:]
			}
			
			switch field {
			case "event":
				s.event.Event = value
			case "data":
				if s.event.Data != "" {
					s.event.Data += "\n"
				}
				s.event.Data += value
			case "id":
				s.event.ID = value
			case "retry":
				s.event.Retry = value
			}
		} else {
			// Field without colon (treat as data)
			if s.event.Data != "" {
				s.event.Data += "\n"
			}
			s.event.Data += line
		}
	}
	
	s.err = s.scanner.Err()
	return false
}

// Event returns the current event
func (s *SSEScanner) Event() SSEEvent {
	return s.event
}

// Err returns any scanning error
func (s *SSEScanner) Err() error {
	return s.err
}