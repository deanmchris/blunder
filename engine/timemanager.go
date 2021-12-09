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

	stopTime        time.Time
	startTime       time.Time
	hardTimeForMove int64
	SoftTimeForMove int64
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

	// If we're given a hard time limit, we're also done calculating, since we've
	// been told already how much time should be spent on the current search.
	if tm.hardTimeForMove != NoValue {
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
	tm.SoftTimeForMove = timeForMove
	tm.hardTimeForMove = NoValue
	tm.Stop = false
}

// Set a hard limit for the maximum amount of time the current
// search can use. Normally this value is set automatically to
// allow the search to use however much time it needs, but there
// are cases where we want to enforce a strict time limit.
//
// This method should not be called after TimeManger.Start has been called,
// since both methods act as the intializers of the current search, and one
// would conflict with the other. So either TimeManager.Start or
// TimeManeger.SetHardTimeForMove should be called to intialize the current
// search time logic, but not both.
func (tm *TimeManager) SetHardTimeForMove(newTime int64) {
	tm.hardTimeForMove = newTime
	tm.stopTime = time.Now().Add(time.Duration(newTime) * time.Millisecond)
}

// Set the soft time limit for current search. The soft time limit
// is a reccomendation for how long the search should continue, but
// can be changed by dynamic factors during the search.
func (tm *TimeManager) SetSoftTimeForMove(newTime int64) {
	// To avoid losing on time, we do enforce a rule that
	// any update to the soft time limit must not exceeded
	// more than 1/8th of the total time left for our side.
	if newTime > tm.TimeLeft/8 {
		newTime = tm.TimeLeft / 8
	}

	// If the hard time limit for this move has already been set, only update the
	// time limit if the hard time limit has not been set, or it is still greater
	// than or equal to the new soft time limit.
	if newTime != NoValue && (tm.hardTimeForMove == NoValue || tm.hardTimeForMove >= newTime) {
		tm.SoftTimeForMove = newTime
		tm.stopTime = tm.startTime.Add(time.Duration(newTime) * time.Millisecond)
	}
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
