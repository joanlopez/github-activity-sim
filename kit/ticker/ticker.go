package ticker

import (
	"context"
	"time"

	"github.com/joanlopez/github-activity-sim/platform/rand"
)

// Variability is the variability percentage
// represented as a floating number: [0,1].
// Example: 0.1 stands for 10% of variability.
type Variability float32

const defaultVariability Variability = 0.35

// Ticker is a time.Ticker wrapper that
// adds a percentage of Variability to the
// given duration.
type Ticker struct {
	C      <-chan time.Time
	d      time.Duration
	v      Variability
	cancel context.CancelFunc
}

// NewTicker initializes a new Ticker with the
// given duration and a default Variability value.
func NewTicker(ctx context.Context, d time.Duration) *Ticker {
	cCtx, cancel := context.WithCancel(ctx)
	channel := make(chan time.Time, 1)
	t := &Ticker{C: channel, d: d, v: defaultVariability, cancel: cancel}

	go func() {
		timer := time.NewTimer(variateDuration(t.d, t.v))

		for {
			select {
			case tc := <-timer.C:
				timer.Reset(variateDuration(t.d, t.v))
				select {
				case channel <- tc:
				default:
				}
			case <-cCtx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return
			}
		}
	}()

	return t
}

// Adjust adjusts ticker's variability
// for those ticks scheduled from now on.
// If given variability is not valid,
// nothing is adjusted.
func (t *Ticker) Adjust(v Variability) {
	if v >= 0 && v <= 1 {
		t.v = v
	}
}

// Stop turns off a ticker. After Stop,
// no more ticks will be sent. Stop does
// not close the channel, to prevent a
// concurrent goroutine reading from the
// channel from seeing an erroneous "tick".
func (t *Ticker) Stop() {
	t.cancel()
}

func variateDuration(d time.Duration, v Variability) time.Duration {
	p := rand.Intn(int(v * 100))
	f := float32(p) / 100
	if rand.Bool() {
		f *= -1
	}
	variance := float32(d) * f
	return time.Duration(float32(d) + variance)
}
