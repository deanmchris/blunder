package ai

import "time"

type Timer struct {
	timeLeft    int64
	increment   int64
	timeForMove int64
	startTime   time.Time
}

func (timer *Timer) UpdateInternals(timeLeft, increment int64) {
	timer.timeLeft = timeLeft
	timer.increment = increment
}

func (timer *Timer) StartSearch() {
	timer.timeForMove = (timer.timeLeft / 40) + (timer.increment / 2)
	if timer.timeForMove >= timer.timeLeft {
		timer.timeForMove -= 500
	}

	if timer.timeForMove <= 0 {
		timer.timeForMove = 100
	}

	timer.startTime = time.Now()
}

func (timer *Timer) TimeIsUp() bool {
	return int64(time.Since(timer.startTime)/time.Millisecond) > timer.timeForMove
}

func (timer *Timer) TimeTaken() int64 {
	return int64(time.Since(timer.startTime) / time.Millisecond)
}
