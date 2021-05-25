package util

import (
	"encoding/binary"
	"fmt"
	"time"
)

const (
	date_format = "2006-01-02 -07:00"
)

type Date struct {
	tm time.Time
}

func MakeDate(t time.Time) Date {
	return Date{
		tm: TruncDate(t),
	}
}

func (dt Date) Equal(o Date) bool {
	return dt.Time().Equal(o.Time())
}

func (dt Date) Time() time.Time {
	return dt.tm
}

func (dt Date) String() string {
	return dt.tm.Format(date_format)
}

func (dt Date) MarshalText() ([]byte, error) {
	return []byte(dt.String()), nil
}

func (dt *Date) UnmarshalText(b []byte) (err error) {
	dt.tm, err = time.Parse(date_format, string(b))
	return
}

func (dt Date) MarshalBinary() (b []byte, err error) {
	y, m, d := dt.tm.Date()
	if y < 1 || y > 4094 {
		err = fmt.Errorf("year outside range 1-4094: %d", y)
		return
	}

	_, z := dt.tm.Zone()

	var n uint32 = 0b1011 << 28
	n |= uint32(y) << 16
	n |= (uint32(m) - 1) << 12
	n |= uint32(d-1) << 7
	n |= uint32(z/900 + 64)

	b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return
}

func (dt *Date) UnmarshalBinary(b []byte) error {
	if len(b) != 4 {
		return fmt.Errorf("invalid length: %d", len(b))
	}

	if b[0]&0xF0 != 0xB0 {
		return fmt.Errorf("invalid type bits: %04b", b[0]>>4)
	}

	y := int(b[0]&0x0F)<<8 | int(b[1])
	if y > 4094 {
		return fmt.Errorf("invalid year: %d", y)
	}
	m := int(b[2]>>4) + 1
	if m > 12 {
		return fmt.Errorf("invalid month: %d", m)
	}
	d := int(b[2]&0x0F<<1) | int(b[3]>>7) + 1
	if d > 31 {
		return fmt.Errorf("invalid day: %d", d)
	}

	z := int(b[3] & 0x7F)
	rz := (z - 64) * 15
	if z > 125 {
		return fmt.Errorf("invalid timezone: %02d:%02d", rz/60, rz%60)
	}

	var loc *time.Location
	if _, local := time.Now().Zone(); rz*60 == local {
		loc = time.Local
	} else {
		loc = time.FixedZone("", rz*60)
	}

	tm := time.Date(y, time.Month(m), d, 0, 0, 0, 0, loc)
	ay, am, ad := tm.Date()
	if ay != y || am != time.Month(m) || ad != d {
		return fmt.Errorf("invalid date: %04d-%02d-%02d", y, m, d)
	}

	dt.tm = tm
	return nil
}
