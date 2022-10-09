package bot

import "time"

type LatencyReader struct {
	latency        time.Duration
	LatencyChannel chan (time.Duration)
}

func NewLatencyReader() *LatencyReader {
	return &LatencyReader{
		latency:        1 * time.Second,
		LatencyChannel: make(chan (time.Duration)),
	}
}

func (l *LatencyReader) Read() {
	for latency := range l.LatencyChannel {
		l.latency = latency
	}
}
