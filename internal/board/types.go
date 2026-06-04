package board

import (
	"fmt"
	"strings"
)

type SideType string

const (
	Field SideType = "field"
	Road  SideType = "road"
	City  SideType = "city"
	River SideType = "river"
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

func (d *Direction) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	switch s {
	case "top":
		*d = Top
	case "right":
		*d = Right
	case "bottom":
		*d = Bottom
	case "left":
		*d = Left
	default:
		return fmt.Errorf("invalid direction: %q", s)
	}
	return nil
}

type UnitType string

const (
	NoneUnit UnitType = ""
	Peasant  UnitType = "peasant"
	Priest   UnitType = "priest"
)
