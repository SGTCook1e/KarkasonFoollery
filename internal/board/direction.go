package board

import (
	"fmt"
	"strings"
)

type Direction int

const (
	Top Direction = iota
	Right
	Bottom
	Left
)

func (d Direction) Opposite() Direction {
	return (d + 2) % 4
}

func (d Direction) Rotate(orientation Direction) Direction {
	return (d + orientation) % 4
}

func (d Direction) Reset(orientation Direction) Direction {
	return (d - orientation + 4) % 4
}

func (d Direction) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	switch s {
	case "top":
		d = Top
	case "right":
		d = Right
	case "bottom":
		d = Bottom
	case "left":
		d = Left
	default:
		return fmt.Errorf("invalid direction: %q", s)
	}
	return nil
}

func (d Direction) MarshalJSON() ([]byte, error) {
	switch d {
	case Top:
		return []byte(`"top"`), nil
	case Right:
		return []byte(`"right"`), nil
	case Bottom:
		return []byte(`"bottom"`), nil
	case Left:
		return []byte(`"left"`), nil
	default:
		return nil, fmt.Errorf("invalid Direction value: %d", d)
	}
}
