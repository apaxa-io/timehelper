package timehelper

import (
	"math"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type testElement struct {
		s   string
		i   Interval
		err bool
	}

	test := []testElement{
		// 0
		testElement{
			s: "-1 year -2 mons +3 days -04:05:06",
			i: Interval{
				Months:  -14,
				Days:    3,
				Seconds: -14706,
			},
			err: false,
		},

		// 1
		testElement{
			s: "-1 year 2 mons -3 days 04:05:06.789",
			i: Interval{
				Months:  -10,
				Days:    -3,
				Seconds: 14706.789,
			},
			err: false,
		},

		// 2
		testElement{
			s:   "",
			i:   Interval{0, 0, 0},
			err: false,
		},

		// 3
		testElement{
			s:   "1 mons",
			i:   Interval{1, 0, 0},
			err: false,
		},

		// 4
		testElement{
			s: "2 year -34:56:78",
			i: Interval{
				Months:  24,
				Days:    0,
				Seconds: -125838,
			},
			err: false,
		},

		// 5
		testElement{
			s:   "00:00:00",
			i:   Interval{0, 0, 0},
			err: false,
		},

		// 6
		testElement{
			s:   "00:00",
			err: true,
		},

		// 7
		testElement{
			s:   "year mons days",
			err: true,
		},

		// 8
		testElement{
			s: "0 year 0 mons 0 days 00:00:00",
			i: Interval{
				Months:  0,
				Days:    0,
				Seconds: 0,
			},
			err: false,
		},

		// 9
		testElement{
			s:   "1.5 year",
			err: true,
		},

		// 10
		testElement{
			s:   "1,5 year",
			err: true,
		},

		// 11
		testElement{
			s:   "99999999999 year -2 mons +3 days -04:05:06",
			err: true,
		},

		// 12
		testElement{
			s:   "9 year 9999999999 mons +3 days -04:05:06",
			err: true,
		},

		// 13
		testElement{
			s:   "9 year -2 mons +99999999999 days -04:05:06",
			err: true,
		},

		// 14
		testElement{
			s:   "9 year -2 mons +9 days 040506",
			err: true,
		},

		// 15
		testElement{
			s:   "9 year -2 mons +9 days 9999999999999999999999999:05:06",
			err: true,
		},

		// 16
		testElement{
			s:   "9 year -2 mons +9 days 04:9999999999999999999999999:06",
			err: true,
		},

		// 17
		testElement{
			s:   "9 year -2 mons +9 days 04:06:99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999",
			err: true,
		},

		// 18
		testElement{
			s: "2147483647 mons 2147483647 days 00:00:00",
			i: Interval{
				Months:  2147483647,
				Days:    2147483647,
				Seconds: 0,
			},
			err: false,
		},

		// 18
		testElement{
			s: "-2147483648 mons -2147483648 days 00:00:00",
			i: Interval{
				Months:  -2147483648,
				Days:    -2147483648,
				Seconds: 0,
			},
			err: false,
		},

		//TODO waiting check overflow
		/*// 19
		testElement{
			s: "2147483647 year 2147483647 mons 2147483647 days 00:00:00",
			i: Interval{
				Months:  2147483647,
				Days:    2147483647,
				Seconds: 0,
			},
			err: false,
		},

		//-2147483648 to 2147483647

		//TODO waiting fix spaces
		// 9
		/*
			testElement{
				s:   "   ",
				err: true,
			},
		*/
	}

	for j, v := range test {
		i, err := Parse(v.s)
		if (err != nil) != v.err {
			t.Errorf("Test-%v, got error: %s", j, err)
		}
		if !v.err && (err == nil) {
			if !i.Equal(v.i) {
				t.Errorf("Test-%v. Intervals not equal.\nExpected:\n%v\ngot:\n%v", j, v.i, i)
			}
		}
	}
}

func TestString(t *testing.T) {
	type testElement struct {
		s   string
		i   Interval
		err bool
	}

	test := []testElement{
		// 0
		testElement{
			s: "-1 year -2 mons 3 days -04:05:06",
			i: Interval{
				Months:  -14,
				Days:    3,
				Seconds: -14706,
			},
			err: false,
		},

		// 1
		testElement{
			s: "-10 mons -3 days 04:05:06.789000000000669",
			i: Interval{
				Months:  -10,
				Days:    -3,
				Seconds: 14706.789,
			},
			err: false,
		},

		// 2
		testElement{
			s:   "1 mons",
			i:   Interval{1, 0, 0},
			err: false,
		},

		// 3
		testElement{
			s: "2 year -34:57:18",
			i: Interval{
				Months:  24,
				Days:    0,
				Seconds: -125838,
			},
			err: false,
		},

		// 4
		testElement{
			s:   "00:00:00",
			i:   Interval{0, 0, 0},
			err: false,
		},

		// 5
		testElement{
			s:   "83 year 4 mons",
			i:   Interval{1000, 0, 0},
			err: false,
		},

		// 6
		testElement{
			s:   "1000 days",
			i:   Interval{0, 1000, 0},
			err: false,
		},

		// 7
		testElement{
			s:   "-1 mons",
			i:   Interval{-1, 0, 0},
			err: false,
		},

		// 8
		testElement{
			s:   "-1 mons",
			i:   Interval{-1, 0, 0},
			err: false,
		},

		//-2147483648 to 2147483647
		// 9
		testElement{
			s:   "178956970 year 7 mons 2147483647 days",
			i:   Interval{2147483647, 2147483647, 0},
			err: false,
		},

		// 10
		testElement{
			s:   "-178956970 year -8 mons -2147483648 days",
			i:   Interval{-2147483648, -2147483648, 0},
			err: false,
		},
	}

	for j, v := range test {
		s := v.i.String()
		if s != v.s {
			t.Errorf("Test-%v. Strings not equal.\nExpected:\n%s\ngot:\n%s", j, v.s, s)
		}
	}
}

func TestAdd(t *testing.T) {
	type testElement struct {
		i   Interval
		add Interval
		res Interval
	}

	test := []testElement{
		// 0
		testElement{
			Interval{-14, 3, -14706},
			Interval{1, 2, 3},
			Interval{-13, 5, -14703},
		},

		// 1
		testElement{
			Interval{},
			Interval{},
			Interval{},
		},

		// 2
		testElement{
			Interval{},
			Interval{-14, 3, -14706},
			Interval{-14, 3, -14706},
		},

		// 3
		testElement{
			Interval{-14, -15, -16},
			Interval{-14, -15, -16},
			Interval{-28, -30, -32},
		},

		// 4
		testElement{
			Interval{-14, -15, -16},
			Interval{14, 15, 16},
			Interval{},
		},

		// 5
		testElement{
			Interval{14, 15, 16},
			Interval{100, 200, 300},
			Interval{114, 215, 316},
		},
		// 6
		testElement{
			Interval{0, 0, 0},
			Interval{14, 15, 16},
			Interval{14, 15, 16},
		},

		// 7
		testElement{
			Interval{14, 15, 16},
			Interval{0, 0, 0},
			Interval{14, 15, 16},
		},
	}

	for j, v := range test {
		i := v.i.Add(v.add)
		if !i.Equal(v.res) {
			t.Errorf("Test-%v. Intervals are not equal.\nExpected:\n%v\ngot:\n%v", j, v.res, i)
		}
	}
}

func TestDuration(t *testing.T) {

	type testElement struct {
		i            Interval
		daysInMonth  uint8
		secondsInDay uint32
		d            time.Duration
	}

	test := []testElement{
		// 0
		testElement{
			Interval{0, 0, 86400},
			30,
			86400,
			86400 * time.Second,
		},

		// 1
		testElement{
			Interval{0, 10, 1},
			30,
			86400,
			864001 * time.Second,
		},

		// 2
		testElement{
			Interval{10, 10, 1},
			30,
			86400,
			26784001 * time.Second, //2562000
		},

		// 3
		testElement{
			Interval{20, 10, 1},
			0,
			0,
			time.Second,
		},

		// 4
		testElement{
			Interval{-10, -5, -1},
			30,
			84000,
			-25620001 * time.Second,
		},

		// 5
		testElement{
			Interval{},
			30,
			84000,
			0,
		},
	}

	for j, v := range test {
		d := v.i.Duration(v.daysInMonth, v.secondsInDay)
		if d != v.d {
			t.Errorf("Test-%v. Wrong duration. Expected: %v, got: %v", j, v.d, d)
		}
	}
}

func TestSub(t *testing.T) {
	type testElement struct {
		i   Interval
		sub Interval
		res Interval
	}

	test := []testElement{
		// 0
		testElement{
			Interval{1, 2, 3},
			Interval{2, 3, 4},
			Interval{-1, -1, -1},
		},

		// 1
		testElement{
			Interval{},
			Interval{2, 3, 4},
			Interval{-2, -3, -4},
		},

		// 2
		testElement{
			Interval{1, 2, 3},
			Interval{},
			Interval{1, 2, 3},
		},

		// 3
		testElement{
			Interval{0, 0, 0},
			Interval{},
			Interval{},
		},

		// 4
		testElement{
			Interval{1, 2, 3},
			Interval{1, 2, 3},
			Interval{},
		},

		// 5
		testElement{
			Interval{-1, -2, -3},
			Interval{1, 2, 3},
			Interval{-2, -4, -6},
		},

		// 6
		testElement{
			Interval{-2147483648, -2147483648, -3},
			Interval{-1, -2, -3},
			Interval{-2147483647, -2147483646, 0},
		},

		// 7
		testElement{
			Interval{2147483647, 2147483647, -3},
			Interval{1, 2, -3},
			Interval{2147483646, 2147483645, 0},
		},
	}

	for j, v := range test {
		s := v.i.Sub(v.sub)
		if !s.Equal(v.res) {
			t.Errorf("Test-%v. Wrong sub.\nExpected interval:%v\ngot:%v", j, v.res, s)
		}
	}
}

func TestMul(t *testing.T) {
	const inaccuracySeconds = 0.0005
	type testElement struct {
		i   Interval
		mul float64
		res Interval
	}

	test := []testElement{
		// 0
		testElement{
			Interval{1, 2, 3},
			2,
			Interval{2, 4, 6},
		},

		// 1
		testElement{
			Interval{1, 2, 3},
			-2,
			Interval{-2, -4, -6},
		},

		// 2
		testElement{
			Interval{1, 2, 3},
			1.05,
			Interval{1, 2, 3.15},
		},

		// 3
		testElement{
			Interval{1, 2, 3},
			0,
			Interval{},
		},

		// 4
		testElement{
			Interval{},
			-2,
			Interval{},
		},
	}

	for j, v := range test {
		i := v.i.Mul(v.mul)
		if (i.Months != v.res.Months) || (i.Days != v.res.Days) || (math.Abs(i.Seconds-v.res.Seconds) > inaccuracySeconds) {
			t.Errorf("Test-%v. Wrong interval.\nExpected:%v\ngot:%v", j, v.res, i)
		}
	}
}

func TestDiv(t *testing.T) {
	const inaccuracySeconds = 0.001
	type testElement struct {
		i   Interval
		div float64
		res Interval
	}

	test := []testElement{
		// 0
		testElement{
			Interval{4, 6, 8},
			2,
			Interval{2, 3, 4},
		},

		// 1
		testElement{
			Interval{4, 6, 8},
			1.1,
			Interval{3, 5, 7.272727},
		},
	}

	for j, v := range test {
		i := v.i.Div(v.div)
		if (i.Months != v.res.Months) || (i.Days != v.res.Days) || (math.Abs(i.Seconds-v.res.Seconds) > inaccuracySeconds) {
			t.Errorf("Test-%v. Wrong interval.\nExpected:%v\ngot:%v", j, v.res, i)
		}
	}
}

func TestEqual(t *testing.T) {
	type testElement struct {
		i   Interval
		i2  Interval
		res bool
	}

	test := []testElement{
		// 0
		testElement{
			Interval{1, 2, 3},
			Interval{2, 3, 4},
			false,
		},

		// 1
		testElement{
			Interval{1, 2, 3},
			Interval{1, 2, 3},
			true,
		},

		// 2
		testElement{
			Interval{},
			Interval{},
			true,
		},

		// 3
		testElement{
			Interval{},
			Interval{1, 2, 3},
			false,
		},

		// 4
		testElement{
			Interval{1, 2, 3},
			Interval{-1, -2, -3},
			false,
		},

		// 5
		testElement{
			Interval{-1, -2, -3},
			Interval{1, 2, 3},
			false,
		},

		// 6
		testElement{
			Interval{-1, -2, -3},
			Interval{-1, -2, -3},
			true,
		},

		// 6
		testElement{
			Interval{-2147483648, -2147483648, -3},
			Interval{-2147483648, -2147483648, -3},
			true,
		},
	}

	for j, v := range test {
		b := v.i.Equal(v.i2)
		if b != v.res {
			t.Errorf("Test-%v. Intervals are not equal.\n1st interval:%v\n2nd interval:%v", j, v.i, v.i2)
		}
	}
}

func TestGreaterOrEqualAndLessOrEqual(t *testing.T) {
	type testElement struct {
		i   Interval
		i2  Interval
		res bool
	}

	test := []testElement{
		// 0
		testElement{
			Interval{1, 2, 3},
			Interval{2, 3, 4},
			false,
		},

		// 1
		testElement{
			Interval{1, 2, 3},
			Interval{1, 2, 3},
			true,
		},

		// 2
		testElement{
			Interval{2, 2, 3},
			Interval{1, 2, 3},
			true,
		},

		// 3
		//damn seconds
		testElement{
			Interval{1, 0, 86400},
			Interval{1, 1, 0},
			false,
		},

		// 4
		//damn seconds
		testElement{
			Interval{1, 0, 186400},
			Interval{1, 1, 0},
			false,
		},

		// 5
		testElement{
			Interval{1, 2, 3},
			Interval{-1, -2, -3},
			true,
		},

		// 6
		testElement{
			Interval{-1, -2, -3},
			Interval{1, 2, 3},
			false,
		},

		// 7
		testElement{
			Interval{},
			Interval{},
			true,
		},

		// 8
		testElement{
			Interval{-2147483648, 2147483647, 0},
			Interval{2147483647, -2147483648, 0},
			false,
		},
	}

	for j, v := range test {
		bG := v.i.GreaterOrEqual(v.i2)
		if bG != v.res {
			t.Errorf("TestGreaterOrEqual - %v. Intervals are not GreaterOrEqual.\n1st interval:%v\n2nd interval:%v", j, v.i, v.i2)
		}
		bL := v.i2.LessOrEqual(v.i)
		if bL != v.res {
			t.Errorf("TestLessOrEqual - %v. Intervals are not LessOrEqual.\n1st interval:%v\n2nd interval:%v", j, v.i2, v.i)
		}
	}
}

func TestLessAndGreater(t *testing.T) {
	type testElement struct {
		i   Interval
		i2  Interval
		res bool
	}

	test := []testElement{
		// 0
		testElement{
			Interval{1, 2, 3},
			Interval{2, 3, 4},
			true,
		},

		// 1
		testElement{
			Interval{1, 2, 3},
			Interval{1, 2, 3},
			false,
		},

		// 2
		testElement{
			Interval{2, 2, 3},
			Interval{1, 2, 3},
			false,
		},

		// 3
		//damn seconds
		testElement{
			Interval{1, 0, 86400},
			Interval{1, 1, 0},
			false,
		},

		// 4
		testElement{
			Interval{-2147483648, -2147483648, 0},
			Interval{2147483647, -2147483647, 0},
			true,
		},
	}

	for j, v := range test {
		bL := v.i.Less(v.i2)
		if bL != v.res {
			t.Errorf("TestLess - %v. 1st interval not less than 2nd.\n1st interval:%v\n2nd interval:%v", j, v.i, v.i2)
		}
		bG := v.i2.Greater(v.i)
		if bG != v.res {
			t.Errorf("TestGreater - %v. 1st interval not greater than 2nd.\n1st interval:%v\n2nd interval:%v", j, v.i2, v.i)
		}
	}
}

func TestComparable(t *testing.T) {
	type testElement struct {
		i   Interval
		i2  Interval
		res bool
	}

	test := []testElement{
		// 0
		testElement{
			Interval{1, 2, 3},
			Interval{2, 3, 4},
			true,
		},

		// 1
		testElement{
			Interval{1, 2, 3},
			Interval{1, 2, 3},
			true,
		},

		// 2
		testElement{
			Interval{2, 2, 3},
			Interval{1, 2, 3},
			true,
		},

		// 3
		testElement{
			Interval{1, 0, 86400},
			Interval{1, 1, 0},
			false,
		},

		// 4
		testElement{
			Interval{1, 0, 186400},
			Interval{1, 1, 0},
			false,
		},

		// 5
		testElement{
			Interval{1, 0, 186400},
			Interval{2, 1, 0},
			false,
		},
	}

	for j, v := range test {
		b := v.i.Comparable(v.i2)
		if b != v.res {
			t.Errorf("Test-%v. Intervals are not Comparable.\n1st interval:%v\n2nd interval:%v", j, v.i2, v.i)
		}
	}
}

func TestAddToAndSubFrom(t *testing.T) {
	const inaccuracySeconds = 5

	type testElement struct {
		i   Interval
		t   time.Time
		res time.Time
	}

	test := []testElement{
		// 0
		testElement{
			Interval{0, 0, 1},
			time.Unix(0, 0),
			time.Unix(1, 0),
		},

		// 1
		testElement{
			Interval{0, 0, 1},
			time.Unix(0, 0),
			time.Unix(1, 0),
		},

		// 2
		testElement{
			Interval{0, 0, 1},
			time.Unix(86400, 0),
			time.Unix(86401, 0),
		},

		// 3
		testElement{
			Interval{0, 0, 1},
			time.Unix(0, 9223372035854775807),
			time.Unix(0, 9223372036854775807),
		},

		// 4
		testElement{
			Interval{0, 0, 9223372036.854775807},
			time.Unix(0, 0),
			time.Unix(9223372036, 854775807),
		},

		// 5
		testElement{
			Interval{0, 0, -9223372036.854775808},
			time.Unix(0, 0),
			time.Unix(0, -9223372036854775808),
		},
	}

	for j, v := range test {
		tA := v.i.AddTo(v.t)
		if time.Duration(math.Abs(float64(tA.Sub(v.res)))) > inaccuracySeconds*time.Second {
			//if t1 != v.res {
			t.Errorf("TestAddTo - %v. Wrong time\nExpected time:\n%v\ngot:\n%v", j, v.res.UTC(), tA.UTC())
		}
		tS := v.i.SubFrom(v.res)
		if time.Duration(math.Abs(float64(tS.Sub(v.res)))) > inaccuracySeconds*time.Second {
			t.Errorf("TestSubFrom - %v. Wrong time\nExpected time:\n%v\ngot:\n%v", j, v.t.UTC(), tS.UTC())
		}
	}

}

func TestNormal(t *testing.T) {
	type testElement struct {
		i    Interval
		year int32
		mon  int32
		day  int32
		hour int32
		min  int8
		sec  int8
		nsec int32
	}

	test := []testElement{
		// 0
		testElement{
			Interval{1001, 101, 1001.3},
			83,
			5,
			101,
			0,
			16,
			41,
			3 * 1e8,
		},

		// 1
		testElement{
			Interval{},
			0,
			0,
			0,
			0,
			0,
			0,
			0,
		},

		// 2
		testElement{
			Interval{-128, 97, 24001.789},
			-10,
			-8,
			97,
			6,
			40,
			1,
			789000000,
		},
	}

	for j, v := range test {
		y := v.i.NormalYears()
		if y != v.year {
			t.Errorf("Test-%v. Ecpected normal year: %v, got: %v", j, v.year, y)
		}
		m := v.i.NormalMonths()
		if m != v.mon {
			t.Errorf("Test-%v. Ecpected normal month: %v, got: %v", j, v.mon, m)
		}
		d := v.i.NormalDays()
		if d != v.day {
			t.Errorf("Test-%v. Ecpected normal days: %v, got: %v", j, v.day, d)
		}
		h := v.i.NormalHours()
		if h != v.hour {
			t.Errorf("Test-%v. Ecpected normal hours: %v, got: %v", j, v.hour, h)
		}
		min := v.i.NormalMinutes()
		if min != v.min {
			t.Errorf("Test-%v. Ecpected normal minutes: %v, got: %v", j, v.min, min)
		}
		s := v.i.NormalSeconds()
		if s != v.sec {
			t.Errorf("Test-%v. Ecpected normal seconds: %v, got: %v", j, v.sec, s)
		}
		ns := v.i.NormalNanoseconds()
		if ns != v.nsec {
			t.Errorf("Test-%v. Ecpected normal nanoseconds: %v, got: %v\n%v", j, v.nsec, ns, s)
		}
	}
}

func TestAll(t *testing.T) {
	i := Nanosecond()
	if i.Seconds != 1e-9 {
		t.Errorf("Error")
	}

	i = Microsecond()
	if i.Seconds != 1e-6 {
		t.Errorf("Error")
	}

	i = Millisecond()
	if i.Seconds != 1e-3 {
		t.Errorf("Error")
	}

	i = Second()
	if i.Seconds != 1 {
		t.Errorf("Error")
	}

	i = Minute()
	if i.Seconds != 60 {
		t.Errorf("Error")
	}

	i = Hour()
	if i.Seconds != 3600 {
		t.Errorf("Error")
	}

	i = Day()
	if i.Days != 1 {
		t.Errorf("Error")
	}

	i = Month()
	if i.Months != 1 {
		t.Errorf("Error")
	}

	i = Year()
	if i.Months != 12 {
		t.Errorf("Error")
	}
}

func TestFromDuration(t *testing.T) {
	type testElement struct {
		i Interval
		d time.Duration
	}
	test := []testElement{
		// 0
		testElement{
			Interval{0, 0, 86400},
			86400 * time.Second,
		},

		// 1
		testElement{
			Interval{0, 0, 8},
			8 * time.Second,
		},

		// 2
		testElement{
			Interval{0, 0, -9223372036.854775808},
			-9223372036854775808,
		},

		// 3
		testElement{
			Interval{0, 0, 9223372036.854775807},
			9223372036854775807,
		},

		// 4
		testElement{
			Interval{},
			0,
		},

		// 5
		testElement{
			Interval{0, 0, -0.000000001},
			-1,
		},
	}
	for j, v := range test {
		i := FromDuration(v.d)
		if i != v.i {
			t.Errorf("Test-%v. Wrong interval. Expected: %v, got: %v", j, v.i, i)
		}
	}
}

func TestDiff(t *testing.T) {
	type testElement struct {
		i    Interval
		from time.Time
		to   time.Time
	}
	test := []testElement{
		// 0
		testElement{
			Interval{0, 0, 1},
			time.Unix(0, 0),
			time.Unix(1, 0),
		},

		// 1
		testElement{
			Interval{0, 0, -1},
			time.Unix(1, 0),
			time.Unix(0, 0),
		},

		// 2
		testElement{
			Interval{0, 0, -0.000000001},
			time.Unix(0, 1),
			time.Unix(0, 0),
		},

		// 3
		testElement{
			Interval{0, 0, 0.000000001},
			time.Unix(0, 0),
			time.Unix(0, 1),
		},

		// 4
		testElement{
			Interval{},
			time.Unix(0, 0),
			time.Unix(0, 0),
		},

		// 5
		testElement{
			Interval{0, 0, 9223372035.854775807},
			time.Unix(1, 0),
			time.Unix(0, 9223372036854775807),
		},

		// 6
		testElement{
			Interval{0, 0, -9223372036.854775807},
			time.Unix(0, 9223372036854775807),
			time.Unix(0, 0),
		},

		// 7
		testElement{
			Interval{0, 0, 9223372036.854775807},
			time.Unix(0, 0),
			time.Unix(9223372036, 854775807),
		},

		// 8
		testElement{
			Interval{0, 0, -9223372036.854775808},
			time.Unix(0, 0),
			time.Unix(0, -9223372036854775808),
		},
	}

	for j, v := range test {
		i := Diff(v.from, v.to)
		if i != v.i {
			t.Errorf("Test-%v. Wrong interval\nExpected:\n%v\ngot:\n%v", j, v.i, i)
		}
	}
}

func TestDiffExtended(t *testing.T) {
	type testElement struct {
		i     Interval
		sFrom string
		sTo   string
	}
	test := []testElement{
		// 0
		testElement{
			Interval{3507, 10, 85636.854775807},
			"1970-01-01T00:00:00Z",
			"2262-04-11T23:47:16.854775807Z",
		},

		// 1
		testElement{
			Interval{},
			"1970-01-01T00:00:00Z",
			"1970-01-01T00:00:00Z",
		},

		// 2
		testElement{
			Interval{0, 0, 1},
			"1970-01-01T00:00:58Z",
			"1970-01-01T00:00:59Z",
		},

		// 3
		testElement{
			Interval{12, 0, 1260},
			"1970-01-01T00:11:00Z",
			"1971-01-01T00:32:00Z",
		},

		// 4
		testElement{
			Interval{12, 0, 3600},
			"1970-01-01T22:00:00Z",
			"1971-01-01T23:00:00Z",
		},

		// 5
		testElement{
			Interval{0, 11, 0},
			"1970-01-14T00:00:00Z",
			"1970-01-25T00:00:00Z",
		},

		// 6
		testElement{
			Interval{8, 11, 0},
			"1970-03-14T00:00:00Z",
			"1970-11-25T00:00:00Z",
		},

		// 7
		testElement{
			Interval{12, 0, 0},
			"1970-01-01T00:00:00Z",
			"1971-01-01T00:00:00Z",
		},

		// 8
		testElement{
			Interval{852, 0, 0},
			"1900-01-01T00:00:00Z",
			"1971-01-01T00:00:00Z",
		},

		// 9
		testElement{
			Interval{-3507, -10, -85636.854775807},
			"2262-04-11T23:47:16.854775807Z",
			"1970-01-01T00:00:00Z",
		},

		// 10
		testElement{
			Interval{-1192, 11, 0},
			"2000-03-01T00:00:00Z",
			"1900-11-12T00:00:00Z",
		},
	}

	for j, v := range test {
		// RFC3339Nano = "2006-01-02T15:04:05.999 999 999Z07:00"
		timeFrom, err := time.Parse(time.RFC3339Nano, v.sFrom)
		if err != nil {
			t.Errorf("Test-%v. Parsing string:%v\ngot err: %v", j, v.sFrom, err)
		}
		timeTo, err1 := time.Parse(time.RFC3339Nano, v.sTo)
		if err1 != nil {
			t.Errorf("Test-%v. Got err: %v, while parsing:%v", j, v.sTo, err1)
		}
		i := DiffExtended(timeFrom, timeTo)
		if i != v.i {
			t.Errorf("Test-%v. Wrong interval\nExpected:\n%v\ngot:\n%v", j, v.i, i)
		}
	}

}

func TestSince(t *testing.T) {
	const inaccuracySeconds = 5
	test := []time.Time{time.Unix(1, 0), time.Unix(1e9, 1e18), time.Unix(0, 0)}
	//TODO check whats wrong with big values
	// max time: time.Unix(1<<63-62135596801, 999999999)
	//time.Unix(- 9223372036854775808, -9223372036854775808)
	for j, v := range test {
		sec := float64(time.Since(v)) / 1e9
		i := Since(v)
		if (i.Months != 0) || (i.Days != 0) || (math.Abs(i.Seconds-sec) > inaccuracySeconds) {
			t.Errorf("Test-%v. Wrong time since: %v\nExpected (time.Since):\n%f\ngot (Since):\n%f", j, v, sec, i.Seconds)
		}
	}

}

func TestSinceExtended(t *testing.T) {
	const inaccuracySeconds = 5
	test := []time.Time{time.Unix(1, 0), time.Unix(1e9, 1e18), time.Unix(0, 0)}
	for j, v := range test {
		i := SinceExtended(v)
		v1 := v.AddDate(0, int(i.Months), int(i.Days))
		v1 = v1.Add(time.Duration(i.Seconds) * time.Second)
		if time.Since(v1) > inaccuracySeconds*time.Second || time.Since(v1) < -inaccuracySeconds*time.Second {
			t.Errorf("Test-%v\nWrong time since: %v\nGit interval:%v\ntime now(v1):%v\nexpected time since(ts):%v", j, v, i, v1, time.Since(v1))
		}

	}
}
