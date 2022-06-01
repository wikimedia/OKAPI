package utils

import "time"

// DaysMap get map that's filled with days in the past (max amount of days back)
func DaysMap(max int) map[string]bool {
	days := map[string]bool{}

	for i := 0; i < max; i++ {
		days[time.Now().UTC().Add(time.Duration(i)*-24*time.Hour).Format(DateFormat)] = true
	}

	return days
}
