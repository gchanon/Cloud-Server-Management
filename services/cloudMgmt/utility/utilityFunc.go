package utility

import (
	"strconv"
	"time"
)

func NewChronoSequence(datetime time.Time, processID int) string {
	processIDStr := strconv.Itoa(processID)
	if len(processIDStr) > 3 {
		processIDStr = processIDStr[len(processIDStr)-3:]
	} else {
		for i := len(processIDStr); i < 3; i++ {
			processIDStr = "0" + processIDStr
		}
	}

	datetimeStr := datetime.Format("060102150405")
	nanosec := strconv.Itoa(datetime.Nanosecond())
	if len(nanosec) != 9 {
		for i := len(nanosec); i < 9; i++ {
			nanosec = "0" + nanosec
		}
	}

	sequenceKey := datetimeStr + nanosec + processIDStr + "A"

	return sequenceKey

}
