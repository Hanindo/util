package util

import (
	"fmt"
	"time"
)

type TZ int

func TimeToTZ(t time.Time) TZ {
	_, offset := t.Zone()
	return TZ(offset / 60)
}

func (tz TZ) Location(t time.Time) *time.Location {
	if _, local := t.Local().Zone(); local == int(tz)*60 {
		return time.Local
	} else {
		return time.FixedZone("", int(tz)*60)
	}
}

func (tz TZ) String() string {
	return time.Time{}.In(time.FixedZone("", int(tz)*60)).Format("-07:00")
}

func (tz TZ) MarshalText() ([]byte, error) {
	return []byte(tz.String()), nil
}

func (tz *TZ) UnmarshalText(b []byte) error {
	t, err := time.Parse(time.RFC3339,
		time.RFC3339[:len(time.RFC3339)-6]+string(b))
	if err != nil {
		return fmt.Errorf("Invalid TZ for %T: %q", tz, b)
	}
	*tz = TimeToTZ(t)
	return nil
}
