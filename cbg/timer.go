// Copyright (c) 2021. Quirino Gervacio
// MIT License. All Rights Reserved

package cbg

import (
	"time"
)

// SleepUntil until the next timestamp
// rounded up to the next period (minutes).
// https://gist.github.com/msadakov/cdbbd979140ef7341fcfac970fc8a95b
func SleepUntil(next int, dur time.Duration) time.Time {
	i := time.Duration(next) * dur
	t := time.Now().Round(i)
	if time.Since(t) >= 0 {
		t = t.Add(i)
	}
	return t
}
