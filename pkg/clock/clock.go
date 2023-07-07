package clock

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (r RealClocker) Now() time.Time {
	return time.Now()
}

type FixedClocker struct{}

func (fc FixedClocker) Now() time.Time {
	return time.Date(2023, 1, 22, 33, 44, 55, 0, time.UTC)
}

type ProgressingClocker struct {
	currnetTime time.Time
	increment   time.Duration
}

func NewProgressingClocker(currentTime time.Time, increment time.Duration) *ProgressingClocker {
	return &ProgressingClocker{
		currnetTime: currentTime,
		increment:   increment,
	}
}

func (pc *ProgressingClocker) Now() time.Time {
	now := pc.currnetTime

	pc.currnetTime = pc.currnetTime.Add(pc.increment)
	return now
}
