package utils

import "time"

// ParseToTimePtr to parse str to time pointer.
func ParseToTimePtr(layout, str string) *time.Time {
	tmp, err := time.Parse(layout, str)
	if err != nil {
		return nil
	}
	return &tmp
}
