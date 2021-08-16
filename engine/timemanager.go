package engine

// timemanager.go implements the time mangement logic which Blunder
// uses during its search phase.

import (
	"time"
)

const NoValue int64 = 0

// A struct which holds data for a timer for Blunder's time mangement.
type TimeManager struct {
	TimeLeft  int64
	Increment int64
	MovesToGo int64
	Stop      bool

	stopTime time.Time
}

// Start the timer, setting up the internal state.
func (tm *TimeManager) Start() {
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
	tm.stopTime = time.Now().Add(time.Duration(timeForMove) * time.Millisecond)
	tm.Stop = false
}

// Check if the time we alloted for picking this move has expired.
func (tm *TimeManager) Check() bool {
	if time.Now().After(tm.stopTime) {
		tm.Stop = true
	}
	return tm.Stop
}
