package utils

import (
	"strconv"
)
// 
const (
	MINUTE = 60
	HOUR = MINUTE * 60
	DAY = HOUR * 24
	MONTH = DAY * 30
	YEAR = MONTH * 12
)

// PadString will add totalSize spaces evenly to the right and left side of str.
// Returns str after applying the pad.
func PadString(str string, totalSize int) string {
	turn := 0
	for {
		if len(str) >= totalSize {
			break
		}
		if turn == 0 {
			str = " " + str
			turn ^= 1
		} else {
			str = str + " "
			turn ^= 1
		}
	}
	return str
}

// FormatUptime will figure out current proc uptime
func FormatUptime(startTime, currentTime int64) string {
	val := currentTime - startTime
	if val < MINUTE {
		return strconv.Itoa(int(val)) + "s"
	} else if val >= MINUTE && val < HOUR {
		return strconv.Itoa(int(val / MINUTE)) + "m"
	} else if val >= HOUR && val < DAY {
		return strconv.Itoa(int(val / HOUR)) + "H"
	} else if val >= DAY && val < MONTH {
		return strconv.Itoa(int(val / DAY)) + "d"
	} else if val >= MONTH && val < YEAR {
		return strconv.Itoa(int(val / MONTH)) + "M"
	} else {
		return strconv.Itoa(int(val / YEAR)) + "y"
	}
}
