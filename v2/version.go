package util

import (
	"fmt"
	"regexp"
	"strconv"
)

// Version represents semantic versioning data.
// This struct has CBOR ``toarray'' struct tags.
type Version struct {
	_     struct{} `cbor:",toarray"`
	Major uint8
	Minor uint8
	Patch uint8
}

// MakeVersion composes a version from its components.
func MakeVersion(major, minor, patch uint8) Version {
	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

// Before reports whether the version v is before o.
func (v Version) Before(o Version) bool {
	return v.Major < o.Major ||
		(v.Major == o.Major && (v.Minor < o.Minor ||
			(v.Minor == o.Minor && v.Patch < o.Patch)))
}

// String returns a string representing the version.
func (v Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// MarshalText implements the encoding.TextMarshaler interface.
// This is basically the String() output.
func (v Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

var versionRE = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)

// UnmarshalText implements the encoding.TextUnmarshaler interface.
// The version is expected exactly like String() format.
func (v *Version) UnmarshalText(b []byte) error {
	if !versionRE.Match(b) {
		return fmt.Errorf("Invalid text for %T: %q", v, b)
	}

	m := versionRE.FindSubmatch(b)
	if major, err := strconv.ParseUint(string(m[1]), 10, 8); err != nil {
		return fmt.Errorf("Invalid major version: %s", m[1])
	} else {
		v.Major = uint8(major)
	}

	if minor, err := strconv.ParseUint(string(m[2]), 10, 8); err != nil {
		return fmt.Errorf("Invalid minor version: %s", m[2])
	} else {
		v.Minor = uint8(minor)
	}

	if patch, err := strconv.ParseUint(string(m[3]), 10, 8); err != nil {
		return fmt.Errorf("Invalid patch version: %s", m[3])
	} else {
		v.Patch = uint8(patch)
	}

	return nil
}
