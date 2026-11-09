package utils

import "time"

func GetExpectedExpiryFromCommand(ttlFormat string, ttlDuration time.Duration) time.Time {
	var expectedDuration time.Time
	switch ttlFormat {
	case "px":
		expectedDuration = time.Now().Add(time.Millisecond * time.Duration(ttlDuration))
	case "ex":
		expectedDuration = time.Now().Add(time.Duration(ttlDuration) * time.Second)
	default:
		expectedDuration = time.Time{}
	}
	return expectedDuration
}
