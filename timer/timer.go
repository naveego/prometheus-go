package timer

import "time"

// Timer represents an API for tracking how long something took to happen
type Timer interface {
	// Start sets the starting point for the timer
	Start() Timer

	// Stop stops the timer and marks the end of the duration
	Stop()

	// Elapsed returns the duration that was tracked
	Elapsed() time.Duration
}

// MemoryTimer is a memory based implementation of timer
type MemoryTimer struct {
	startTime time.Time
	elapsed   time.Duration
}

// Start the timer
func (mt *MemoryTimer) Start() Timer {
	mt.startTime = time.Now()
	return mt
}

// Stop the timer
func (mt *MemoryTimer) Stop() {
	mt.elapsed = time.Now().Sub(mt.startTime)
}

// Elapsed returns the duration
func (mt *MemoryTimer) Elapsed() time.Duration {
	return mt.elapsed
}
