package lib

import "time"

// returns UTC time, ie: 2023-07-06T07:25:26Z
func TimeStamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
