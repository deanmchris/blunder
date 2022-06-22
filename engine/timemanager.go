package engine

// timemanager.go implements the time mangement logic which Blunder
// uses during its search phase. It controls when the search should be
// stopped given the arguments passed through the UCI "go" command,
// as well as dynamic factors of the search.

import (
	"time"
)

const (
	NoValue      int64 = 0
	InfiniteTime int64 = -1
)

// A struct which holds data for a timer for Blunder's time mangement.
type TimeManager struct {
	// Fields for UCI go command arguments
	TimeLeft     int64
	Increment    int64
	MoveTime     int64
	MovesToGo    int16
	MaxNodeCount uint64
	MaxDepth     uint8

	// Fields to calculate when the search should be stopped.
	Stop        bool
	TimeForMove int64
	stopTime    time.Time
}

// Setup the interals of the timer given the "go" command arguments.
func (tm *TimeManager) Setup(timeLeft, increment, moveTime int64,
	movesToGo int16, maxDepth uint8, maxNodeCount uint64) {

	tm.TimeLeft = timeLeft
	tm.Increment = increment
	tm.MovesToGo = movesToGo
	tm.MoveTime = moveTime
	tm.MaxDepth = maxDepth
	tm.MaxNodeCount = maxNodeCount
}

// Start the timer, setting up the internal state.
func (tm *TimeManager) Start() {
	// Reset the flag time's up flag to false for a new search
	tm.Stop = false

	// Prioritize the "movetime" argument if a value is given and use that.
	if tm.MoveTime != NoValue {
		tm.stopTime = time.Now().Add(time.Duration(tm.MoveTime) * time.Millisecond)
		tm.TimeLeft = NoValue
		return
	}

	// If we're given infinite time, we're done calculating the time for the
	// current move.
	if tm.TimeLeft == InfiniteTime {
		return
	}

	// Otherwise, automatically calculate the time we can allocate for the search about to start.

	timeForMove := int64(0)
	if int64(tm.MovesToGo) != NoValue {
		// If we have a certian amount of moves to go before the time we have left
		// is reset, use that value to divide the time currently left.
		timeForMove = tm.TimeLeft / int64(tm.MovesToGo)
	} else {
		// Otherwise, calculate the amount of remaining time that will
		// be used by spending more time as the game progresses,
		// assuming that the longer the game continues, the quicker
		// it will end from the current position.
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
	tm.TimeForMove = timeForMove
}

// Update the alloted search time.
func (tm *TimeManager) Update(newTimeForMove int64) {
	// If we've been given an explict amount of search time,
	// respect it.
	if tm.MoveTime != NoValue {
		return
	}

	// To avoid losing on time, enforce a rule that
	// any update to the time left must not exceeded
	// more than 1/8th of the total time left for our side.
	if newTimeForMove > tm.TimeLeft/8 {
		newTimeForMove = tm.TimeLeft / 8
	}

	// Set the new time for the current search.
	tm.stopTime = time.Now().Add(time.Duration(newTimeForMove) * time.Millisecond)
	tm.TimeForMove = newTimeForMove
}

// Check if the time we alloted for picking this move has expired.
func (tm *TimeManager) Check() {
	// If we've already been told to stop before now,
	// no more work needs to be done.
	if tm.Stop {
		return
	}

	// If we have infinite time, we don't need to check if our time is up.
	if tm.TimeLeft == InfiniteTime {
		return
	}

	// Otherwise check if our alloted time is over.
	if time.Now().After(tm.stopTime) {
		tm.Stop = true
	}
}
