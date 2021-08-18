package util

import (
	"encoding/binary"
	"fmt"
	"time"
)

const (
	date_format = "2006-01-02 -07:00"
)

// Date represents a date with timezone.
//
// This date can be efficiently encoded to 4 bytes using binary format similar
// to Temporenc https://temporenc.org but with custom type tag.
// The bits are 4 bits type tag ``0b1011'', 21 bits of date component
// (y=12 + m=4 + d=5) and 7 bits of time zone offset component.
type Date struct {
	tm time.Time
}

// MakeDate returns a date from t time.
func MakeDate(t time.Time) Date {
	return Date{
		tm: TruncDate(t),
	}
}

// Equal report whether d and o represent the same date. This method is
// necessary because date is implemented using time.Time.
func (d Date) Equal(o Date) bool {
	return d.Time().Equal(o.Time())
}

// Time returns the starting time of date.
func (d Date) Time() time.Time {
	return d.tm
}

// String returns date formatted using time format string
//     "2006-01-02 -07:00"
func (d Date) String() string {
	return d.tm.Format(date_format)
}

// MarshalText implements the encoding.TextMarshaler interface.
// This is basically the String() output.
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The date is expected exactly like String() format.
func (d *Date) UnmarshalText(b []byte) (err error) {
	d.tm, err = time.Parse(date_format, string(b))
	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
// See the documentation on the Date type for more details.
func (d Date) MarshalBinary() (b []byte, err error) {
	y, m, day := d.tm.Date()
	if y < 0 || y > 4094 {
		err = fmt.Errorf("year outside range 0-4094: %d", y)
		return
	}

	_, z := d.tm.Zone()

	var n uint32 = 0xB << 28
	n |= uint32(y) << 16
	n |= (uint32(m) - 1) << 12
	n |= uint32(day-1) << 7
	n |= uint32(z/900 + 64)

	b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
// See the documentation on the Date type for more details.
func (d *Date) UnmarshalBinary(b []byte) error {
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
	day := int(b[2]&0x0F<<1) | int(b[3]>>7) + 1
	if day > 31 {
		return fmt.Errorf("invalid day: %d", day)
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

	tm := time.Date(y, time.Month(m), day, 0, 0, 0, 0, loc)
	ay, am, ad := tm.Date()
	if ay != y || am != time.Month(m) || ad != day {
		return fmt.Errorf("invalid date: %04d-%02d-%02d", y, m, day)
	}

	d.tm = tm
	return nil
}
