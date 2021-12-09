package engine

// timemanager.go implements the time mangement logic which Blunder
// uses during its search phase.

import (
	"time"
)

const (
	NoValue      int64 = 0
	InfiniteTime int64 = -1
)

// A struct which holds data for a timer for Blunder's time mangement.
type TimeManager struct {
	TimeLeft  int64
	Increment int64
	MovesToGo int64
	Stop      bool

	stopTime    time.Time
	startTime   time.Time
	TimeForMove int64
}

// Start the timer, setting up the internal state.
func (tm *TimeManager) Start() {
	// Reset the flag time's up flag to false for a new search
	tm.Stop = false

	// If we're given infinite time, we're done calculating the time for the
	// current move.
	if tm.TimeLeft == InfiniteTime {
		return
	}

	// Calculate the time we can allocate for the search about to start.
	var timeForMove int64

	if tm.MovesToGo != NoValue {
		// If we have a certian amount of moves to go before the time we have left
		// is reset, use that value to divide the time currently left.
		timeForMove = tm.TimeLeft / tm.MovesToGo
	} else {
		// Otherwise get 2.5% of the current time left and use that.
		timeForMove = tm.TimeLeft / 40
	}

	// Give an bonus from the increment
	timeForMove += tm.Increment / 2

	// If the increment bonus puts us outside of the actual time we
	// have left, use the time we have left minus 500ms.
	if timeForMove >= tm.TimeLeft {
		timeForMove = tm.TimeLeft - 500
	}

	// If taking away 500ms puts us below zero, use 100ms
	// to just get a move to return.
	if timeForMove <= 0 {
		timeForMove = 100
	}

	// Calculate the time from now when we need to stop searching, based on the
	// time are allowed to spend on the current search.
	tm.startTime = time.Now()
	tm.stopTime = tm.startTime.Add(time.Duration(timeForMove) * time.Millisecond)
	tm.TimeForMove = timeForMove
	tm.Stop = false
}

// Update the alloted time for the current search.
func (tm *TimeManager) Update(newTime int64) {
	if newTime > tm.TimeLeft/8 {
		newTime = tm.TimeLeft / 8
	}

	tm.TimeForMove = newTime
	tm.stopTime = tm.startTime.Add(time.Duration(newTime) * time.Millisecond)
}

// Check if the time we alloted for picking this move has expired.
func (tm *TimeManager) Check() {
	// If we have infinite time, tm.Stop is set to false unless we've already
	// been told to stop.
	if !tm.Stop && tm.TimeLeft == InfiniteTime {
		tm.Stop = false
		return
	}

	// Otherwise figure out if our alloated time for this move is up.
	if time.Now().After(tm.stopTime) {
		tm.Stop = true
	}
}
