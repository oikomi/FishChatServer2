package xtime

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// Time be used to MySql timestamp converting.
type Time int64

func (jt *Time) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case time.Time:
		*jt = Time(sc.Unix())
	case string:
		var i int64
		i, err = strconv.ParseInt(sc, 10, 64)
		*jt = Time(i)
	}
	return
}

func (jt Time) Value() (driver.Value, error) {
	return time.Unix(int64(jt), 0), nil
}

func (jt Time) Time() time.Time {
	return time.Unix(int64(jt), 0)
}

// Duration be used toml unmarshal string time, like 1s, 500ms.
type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	tmp, err := time.ParseDuration(string(text))
	if err == nil {
		*d = Duration(tmp)
	}
	return err
}
