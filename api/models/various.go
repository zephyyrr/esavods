package models

import (
	"fmt"
)

type Id string

type MilliSecondDuration int64

const (
	MilliSecond MilliSecondDuration = 1
	Second                          = 1000 * MilliSecond
	Minute                          = 60 * Second
	Hour                            = 60 * Minute
)

func (d MilliSecondDuration) Hours() int64 {
	return int64(d / Hour)
}

func (d MilliSecondDuration) Minutes() int64 {
	return int64(d / Minute)
}

func (d MilliSecondDuration) Seconds() int64 {
	return int64(d / Second)
}

func (d MilliSecondDuration) String() string {
	hours, d := d/Hour, d%Hour
	minutes, d := d/Minute, d%Minute
	seconds, d := d/Second, d%Second
	hundredths := d / 10
	return fmt.Sprintf("%02d:%02d:%02d.%02d", hours, minutes, seconds, hundredths)
}
