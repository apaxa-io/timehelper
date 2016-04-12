package timehelper

import "time"

// SimpleLayout is additional predefined layout for use in Time.Format and Time.Parse.
// The reference time used in the layouts is the specific time: "Mon Jan 2 15:04:05 MST 2006".
const SimpleLayout = "2006-01-02 15:04:05"

// UnixEpoch returns local Time corresponding to the beginning of UNIX epoch.
func UnixEpoch() time.Time {
	return time.Unix(0, 0)
}
