package portal

import "time"

type session_date struct {
	Date    int
	TimeIn  time.Time
	TimeOut time.Time
	// LateEntryTime  float32
	// EarlyEntryTime float32
}

type filter_auth_session_month struct {
	StudentId   string
	Month       time.Time
	ListSession []session_date
	// TotalLateEntryTime  float32
	// TotalEarlyEntryTime float32
}

var (
	AuthList = filter_auth_session_month{
		StudentId: "n18dccn229",
		Month:     time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
		ListSession: []session_date{
			{
				Date:    1,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    2,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    3,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    4,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    5,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    6,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    7,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    8,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    9,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    10,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    11,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    12,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    13,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    14,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    15,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},

			{
				Date:    16,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    17,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    18,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    19,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    20,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    21,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    22,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    23,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    24,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    25,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},

			{
				Date:    26,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    27,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    28,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    29,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    30,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
			{
				Date:    31,
				TimeIn:  time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
				TimeOut: time.Date(2022, 11, 15, 9, 30, 12, 00, time.Local),
			},
		},
	}
)
