package tile

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
