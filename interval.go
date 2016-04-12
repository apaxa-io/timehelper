package timehelper

import (
	"errors"
	"fmt"
	"github.com/apaxa-io/mathhelper"
	"math"
	"regexp"
	"strconv"
	"time"
)

// Interval represent time interval.
// It consists of 3 fields:
// 	Months - number months (integer)
// 	Days - number of days (integer)
// 	Seconds - number of seconds (real)
// All fields are signed. Sign of one field is independent from sign of other field.
// This type is equal to Postgres interval data type.
// Value from one field is never translated to value of another field, so <60*60*24 seconds> != <1 days> and so on.
// This is because of:
// 	1) compatibility with Postgres;
// 	2) day may have different amount of seconds and month may have different amount of days.
type Interval struct {
	Months  int32
	Days    int32
	Seconds float64
}

// RE for parse interval in postgres style specification.
// http://www.postgresql.org/docs/9.4/interactive/datatype-datetime.html#DATATYPE-INTERVAL-OUTPUT
var re = regexp.MustCompile(`^(?:([+-]?[0-9]+) year)? ?(?:([+-]?[0-9]+) mons)? ?(?:([+-]?[0-9]+) days)? ?(?:([+-])?([0-9]+):([0-9]+):([0-9]+(?:,|.[0-9]+)?))?$`)

// Nanosecond returns new Interval equal to 1 nanosecond
func Nanosecond() Interval {
	return Interval{Seconds: 1e-9}
}

// Microsecond returns new Interval equal to 1 microsecond
func Microsecond() Interval {
	return Interval{Seconds: 1e-6}
}

// Millisecond returns new Interval equal to 1 millisecond
func Millisecond() Interval {
	return Interval{Seconds: 1e-3}
}

// Second returns new Interval equal to 1 second
func Second() Interval {
	return Interval{Seconds: 1}
}

// Minute returns new Interval equal to 1 minute (60 seconds)
func Minute() Interval {
	return Interval{Seconds: 60}
}

// Hour returns new Interval equal to 1 hour (3600 seconds)
func Hour() Interval {
	return Interval{Seconds: 3600}
}

// Day returns new Interval equal to 1 day
func Day() Interval {
	return Interval{Days: 1}
}

// Month returns new Interval equal to 1 month
func Month() Interval {
	return Interval{Months: 1}
}

// Year returns new Interval equal to 1 year (12 months)
func Year() Interval {
	return Interval{Months: 12}
}

// Parse parses incoming string and extract interval.
// Format is postgres style specification for interval output format.
// Examples:
// 	-1 year 2 mons -3 days 04:05:06.789
// 	1 mons
// 	2 year -34:56:78
// 	00:00:00
func Parse(s string) (i Interval, err error) {
	//TODO string of 1-3 spaces are parse ok
	//TODO add check for overflow

	parts := re.FindStringSubmatch(s)
	if parts == nil || len(parts) != 8 {
		err = errors.New("Unable to parse interval from string " + s)
		return
	}

	var ti int64
	var tf float64
	var negativeTime bool

	// Store as months:

	// years
	if parts[1] != "" {
		ti, err = strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return
		}
		i.Months = int32(ti) * 12
	}

	// months
	if parts[2] != "" {
		ti, err = strconv.ParseInt(parts[2], 10, 32)
		if err != nil {
			return
		}
		i.Months += int32(ti)
	}

	// Store as days:

	// days
	if parts[3] != "" {
		ti, err = strconv.ParseInt(parts[3], 10, 32)
		if err != nil {
			return
		}
		i.Days = int32(ti)
	}

	// Store as seconds:

	negativeTime = parts[4] == "-"

	// hours
	if parts[5] != "" {
		ti, err = strconv.ParseInt(parts[5], 10, 64)
		if err != nil {
			return
		}
		i.Seconds = float64(ti) * 3600
	}

	// minutes
	if parts[6] != "" {
		ti, err = strconv.ParseInt(parts[6], 10, 64)
		if err != nil {
			return
		}
		i.Seconds += float64(ti) * 60
	}

	// seconds
	if parts[7] != "" {
		tf, err = strconv.ParseFloat(parts[7], 64)
		if err != nil {
			return
		}
		i.Seconds += tf
	}

	// fix sign if it is negative
	if negativeTime {
		i.Seconds *= -1
	}

	return
}

// FromDuration returns new Interval equivalent for given time.Duration (convert time.Duration to Interval).
func FromDuration(d time.Duration) Interval {
	return Interval{Seconds: float64(d.Nanoseconds()) / 1e9}
}

// Diff calculates difference between given timestamps (time.Time) as seconds and returns result as Interval (=to-from).
// Result always have months & days parts set to zero.
func Diff(from, to time.Time) Interval {
	return Interval{Seconds: (float64(to.UnixNano()) - float64(from.UnixNano())) / 1e9}
}

// DiffExtended is similar to Diff but calculates difference in months, days & seconds instead of just seconds (=to-from).
// Result may have non-zero months & days parts.
func DiffExtended(from, to time.Time) Interval {
	fromYear, fromMonth, fromDay := from.Date()
	fromHour, fromMin, fromSec := from.Clock()
	fromNsec := from.Nanosecond()

	toYear, toMonth, toDay := to.Date()
	toHour, toMin, toSec := to.Clock()
	toNsec := to.Nanosecond()

	return Interval{
		Months:  int32((toYear-fromYear)*12 + int(toMonth-fromMonth)),
		Days:    int32(toDay - fromDay),
		Seconds: float64((toHour-fromHour)*3600+(toMin-fromMin)*60+(toSec-fromSec)) + float64(toNsec-fromNsec)/1e9,
	}
}

// Since returns elapsed time since given timestamp as Interval (=Diff(t, time.New())
// Result always have months & days parts set to zero.
func Since(t time.Time) Interval {
	return Diff(t, time.Now())
}

// SinceExtended returns elapsed time since given timestamp as Interval (=DiffExtended(t, time.New())
// Result may have non-zero months & days parts.
func SinceExtended(t time.Time) Interval {
	return DiffExtended(t, time.Now())
}

// String returns string representation of interval.
// Output format is the same as for Parse
func (i Interval) String() string {
	if i.Months == 0 && i.Days == 0 && i.Seconds == 0 {
		return "00:00:00"
	}

	y := i.NormalYears()
	mon := i.NormalMonths()

	negativeTime := i.Seconds < 0
	if negativeTime {
		i.Seconds = -i.Seconds
	}

	h := i.NormalHours()
	m := i.NormalMinutes()
	s := math.Remainder(i.Seconds, 60)

	str := ""
	if y != 0 {
		str += strconv.FormatInt(int64(y), 10) + " year "
	}
	if mon != 0 {
		str += strconv.FormatInt(int64(mon), 10) + " mons "
	}
	if i.Days != 0 {
		str += strconv.FormatInt(int64(i.Days), 10) + " days "
	}
	if i.Seconds != 0 {
		if negativeTime {
			str += "-"
		}
		secStr := strconv.FormatFloat(s, 'f', -1, 64)
		if s < 10 {
			secStr = "0" + secStr
		}
		str += fmt.Sprintf("%02d:%02d", h, m) + ":" + secStr

		return str
	}
	// As all null interval filtered at the beginning of method there is a space at the end of string
	return str[:len(str)-1]
}

// Duration convert Interval to time.Duration.
// It is required to pass number of days in mounth (usually 30 or something near)
// and number of seconds in day (usually 86400)
// because of converting mounths and days parts of original Interval to time.Duration seconds.
// Warning: this method is inaccuracy because in real life daysInMonth & secondsInDay vary and depends on relative timestamp.
func (i Interval) Duration(daysInMonth uint8, secondsInDay uint32) time.Duration {
	return time.Duration((int64(i.Months)*int64(daysInMonth)+int64(i.Days))*int64(secondsInDay)*1e9 + mathhelper.Round(i.Seconds*1e9))
}

// Add adds given Interval to original Interval.
// Original Interval will be changed.
func (i Interval) Add(add Interval) Interval {
	i.Months += add.Months
	i.Days += add.Days
	i.Seconds += add.Seconds
	return i
}

// Sub subtracts given Interval from original Interval.
// Original Interval will be changed.
func (i Interval) Sub(sub Interval) Interval {
	i.Months -= sub.Months
	i.Days -= sub.Days
	i.Seconds -= sub.Seconds
	return i
}

// Mul multiples interval by mul. Each part of Interval multiples independently.
// Multiply by non integer value give expected value only if Months and Days are zero, as them stored as integer value.
// Original Interval will be changed.
func (i Interval) Mul(mul float64) Interval {
	i.Months, i.Days, i.Seconds = int32(float64(i.Months)*mul), int32(float64(i.Days)*mul), i.Seconds*mul
	return i
}

// Div divides interval by mul. Each part of Interval divides independently.
// Warning: Div may returns unexpected value because Months and Days stored as integer value.
// Original Interval will be changed.
func (i Interval) Div(div float64) Interval {
	i.Months = int32(float64(i.Months) / div)
	i.Days = int32(float64(i.Days) / div)
	i.Seconds = i.Seconds / div
	return i
}

// Comparable returns true only if it is possible to compare Intervals.
// Intervals "A" and "B" can be compared only if:
//   1) all parts of "A" are less or equal to relative parts of "B"
//   or
//   2) all parts of "B" are less or equal to relative parts of "A".
// In the other words, it is impossible to compare "30 days"-Interval with "1 month"-Interval.
func (i Interval) Comparable(i2 Interval) bool {
	return i.LessOrEqual(i2) || i.GreaterOrEqual(i2)
}

// Equal compare original Interval with given for full equality part by part.
// Warning: Seconds parts is float but compared strictly.
func (i Interval) Equal(i2 Interval) bool {
	return i.Months == i2.Months && i.Days == i2.Days && i.Seconds == i2.Seconds
}

// LessOrEqual returns true if all parts of original Interval are less or equal to relative parts of i2.
func (i Interval) LessOrEqual(i2 Interval) bool {
	return i.Months <= i2.Months && i.Days <= i2.Days && i.Seconds <= i2.Seconds
}

// Less returns true if at least one part of original Interval is less then relative part of i2 and all other parts of original Interval are less or equal to relative parts of i2.
func (i Interval) Less(i2 Interval) bool {
	return !i.Equal(i2) && i.LessOrEqual(i2)
}

// GreaterOrEqual returns true if all parts of original Interval are greater or equal to relative parts of i2.
func (i Interval) GreaterOrEqual(i2 Interval) bool {
	return i.Months >= i2.Months && i.Days >= i2.Days && i.Seconds >= i2.Seconds
}

// Greater returns true if at least one part of original Interval is greater then relative part of i2 and all other parts of original Interval are greater or equal to relative parts of i2.
func (i Interval) Greater(i2 Interval) bool {
	return !i.Equal(i2) && i.GreaterOrEqual(i2)
}

// NormalYears return number of years in month part (as i.Months / 12).
func (i Interval) NormalYears() int32 {
	// TODO what about sign?
	return i.Months / 12
}

// NormalMonths return number of months in month part after subtracting NormalYears*12 (as i.Months % 12).
// Examples: if .Months = 11 then NormalMonths = 11, but if .Months = 13 then NormalMonths = 1.
func (i Interval) NormalMonths() int32 {
	// TODO what about sign?
	return i.Months % 12
}

// NormalDays just returns Days part.
func (i Interval) NormalDays() int32 {
	return i.Days
}

// NormalHours returns number of hours in seconds part (as i.Seconds / 3600).
func (i Interval) NormalHours() int32 {
	// TODO what about sign?
	return int32(i.Seconds / 3600)
}

// NormalMinutes returns number of hours in seconds part after subtracting NormalHours*60 (as (i.Seconds - i.NormalHours()*3600) / 60).
func (i Interval) NormalMinutes() int8 {
	// TODO what about sign?
	return int8((i.Seconds - float64(i.NormalHours())*3600) / 60)
}

// NormalSeconds returns number of seconds in seconds part after subtracting NormalHours*3600 and NormalMinutes*60 (as i.Seconds % 60).
func (i Interval) NormalSeconds() int8 {
	// TODO what about sign?
	return int8(int64(i.Seconds) % 60)
}

// NormalNanoseconds returns number of nanoseconds in fraction part of seconds part.
func (i Interval) NormalNanoseconds() int32 {
	//TODO find all remainder - it isnt remainder
	// TODO what about sign?
	return int32(mathhelper.Round(math.Mod(i.Seconds, 1) * 1e9))
}

// AddTo adds original Interval to given timestamp and return result.
func (i Interval) AddTo(t time.Time) time.Time {
	//TODO report bug (not working on large seconds in interval without converting to utc)
	location := t.Location()
	t = t.UTC()

	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	nsec := t.Nanosecond()
	//return time.Date(year, month+time.Month(i.Months), day+int(i.Days), hour, min, sec+int(i.Seconds), nsec+int(i.NormalNanoseconds()), t.Location())
	t = time.Date(year, month+time.Month(i.Months), day+int(i.Days), hour, min, sec+int(i.Seconds), nsec+int(i.NormalNanoseconds()), time.UTC)
	return t.In(location)
}

// SubFrom subtract original Interval from given timestamp and return result.
func (i Interval) SubFrom(t time.Time) time.Time {
	return i.Mul(-1).AddTo(t)
}
