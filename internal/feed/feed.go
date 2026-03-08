package feed

import (
	"math"
	"math/rand"
	"time"
)

// DataPoint represents a single time series value.
type DataPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// SineProducer generates a sine wave with added noise.
type SineProducer struct {
	Frequency  float64
	Amplitude  float64
	NoiseLevel float64
	SampleRate float64
	ticker     *time.Ticker
	quit       chan struct{}
}

// NewSineProducer creates a new producer with default 50Hz sample rate.
func NewSineProducer(frequency, amplitude, noiseLevel float64) *SineProducer {
	return &SineProducer{
		Frequency:  frequency,
		Amplitude:  amplitude,
		NoiseLevel: noiseLevel,
		SampleRate: 50.0, // 50 Hz
		quit:       make(chan struct{}),
	}
}

// Start begins generating data points and sending them to the output channel.
func (p *SineProducer) Start() <-chan DataPoint {
	out := make(chan DataPoint)
	interval := time.Second / time.Duration(p.SampleRate)
	p.ticker = time.NewTicker(interval)

	go func() {
		defer close(out)
		start := time.Now()
		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)

		for {
			select {
			case <-p.ticker.C:
				elapsed := time.Since(start).Seconds()
				// val = A * sin(2 * pi * f * t) + noise
				value := p.Amplitude*math.Sin(2*math.Pi*p.Frequency*elapsed) + (r.Float64()*2-1)*p.NoiseLevel
				out <- DataPoint{
					Timestamp: time.Now().UnixMilli(),
					Value:     value,
				}
			case <-p.quit:
				p.ticker.Stop()
				return
			}
		}
	}()

	return out
}

// Stop halts the data generation.
func (p *SineProducer) Stop() {
	close(p.quit)
}
