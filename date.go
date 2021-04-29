package util

import (
	"fmt"
	"regexp"
	"time"
)

const (
	DATE_FORMAT = "2006-01-02"
)

type Date struct {
	_    struct{} `cbor:",toarray"`
	Time time.Time
	TZ   TZ
}

func MakeDate(t time.Time) Date {
	return Date{
		Time: TruncDate(t),
		TZ:   TimeToTZ(t),
	}
}

func (d Date) GetTime() time.Time {
	return d.Time.In(d.TZ.Location(d.Time))
}

func (d Date) String() string {
	return d.GetTime().Format(DATE_FORMAT) + " " + d.TZ.String()
}

func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

var dateRE = regexp.MustCompile(`^` +
	`(\d{4}-(?:0[1-9]|1[0-2])-(?:0[1-9]|[12]\d|3[01]))` + // YYYY-MM-DD
	` ` +
	`([+-]\d\d:\d\d)` + // +- hh:mm
	`$`)

func (d *Date) UnmarshalText(b []byte) error {
	if len(b) != len(DATE_FORMAT)+1+6 || !dateRE.Match(b) {
		return fmt.Errorf("Invalid text for %T: %q", d, b)
	}
	m := dateRE.FindSubmatch(b)
	if err := d.TZ.UnmarshalText(m[2]); err != nil {
		return fmt.Errorf("Invalid tz text for %T: %q", d, b)
	}

	var err error
	if d.Time, err = time.ParseInLocation(
		DATE_FORMAT, string(m[1]), d.TZ.Location(d.Time),
	); err != nil {
		return fmt.Errorf("Invalid date text for %T: %q", d, b)
	}
	return nil
}
